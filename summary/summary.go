package summary

import (
	"fmt"
	"sort"
	"strings"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

/*
Flows thoughts
--
	Src IP     ->    Dst IP
	|                |
	Src Port         Dst Port
	(random)         (concrete svc)
--
	Src IP     ->    Multicast IP
	|                |
	Src Port         Dst Port
	(random)         (concrete svc)
--
	Unicast traffic schema:
		From IP1:r to IP2:c
		From IP2:c to IP1:r
	The key to chain it in one sequence is {IP1, IP2, c} and it applies as
		forth: Src IP, Dst IP, Dst Port
			   =       =       =
		back:  Dst IP, Src IP, Src Port
		       =       =       =
			   IP1,    IP2,    concrete port
--
	Multicast traffic schema:
		From IP1:r to Multicast:c
		From IP2:c to IP1:r
	The keys to chain it in several rows of one sequence will look like
		reference row
			forth: Src IP, Dst Multicast IP, Dst Port
			       =       =                 =
			back:  Dst IP, Src Multicast IP, Src Port
		           =       =                 =
			       IP1,    Multicast IP,     concrete port
		response to multicast rows
			back-mulicast: Dst IP, Src IP, Src Port
      		               =       =       =
			               IP1,    IP2,    concrete port
*/

/*
	Summary tree:

	Samplers
	-	Flows groups
		-	Flows rows of the same group
			- 	Packet counters for identical flows
*/

type summary struct {
	samplers map[string]*flows
}

type flows struct {
	flows map[string]*rows
}

type rows struct {
	dst map[string]*packets
}

func NewSummary() *summary {
	return &summary{}
}

func (s *summary) Add(fmsg *flowmessage.FlowMessage) {
	if s == nil {
		return
	}

	if s.samplers == nil {
		s.samplers = make(map[string]*flows)
	}

	msg := message{FlowMessage: fmsg}

	// find sampler to aggregate this new message
	samplerKey := key(msg).sampler()
	sampler, ok := s.samplers[samplerKey]
	if !ok {
		sampler = &flows{
			flows: make(map[string]*rows),
		}
		s.samplers[samplerKey] = sampler
	}

	// check if this flow is answer to some other flow
	var isAnswer bool
	if _, ok := sampler.flows[key(msg).parent()]; ok {
		isAnswer = true
	}
	if _, ok := sampler.flows[key(msg).parentMulticast()]; ok {
		isAnswer = true
	}

	// select flows key
	var flowKey string
	if isAnswer {
		flowKey = key(msg).asFromSrc()
	} else {
		flowKey = key(msg).asToDst()
	}

	// find flows group for new flow
	flowsGroup, ok := sampler.flows[flowKey]
	if !ok {
		flowsGroup = &rows{
			dst: make(map[string]*packets),
		}
		sampler.flows[flowKey] = flowsGroup
	}

	// find appropriate row with the same destination
	dstKey := key(msg).destination()
	row, ok := flowsGroup.dst[dstKey]
	if !ok {
		row = &packets{
			description: msg.description(isAnswer),
		}
		flowsGroup.dst[dstKey] = row
	}

	// append bytes and packets
	row.bytes += msg.Bytes
	row.packets += msg.Packets
}

func (s *summary) Reset() {
	if s == nil {
		return
	}

	// firstly clear counters so it leaves back-and-forth keys for furthers flows
	// next Reset() will delete all empty untouched flows

	for samplerKey, sampler := range s.samplers {
		for flowKey, flowsGroup := range sampler.flows {
			for dstKey, row := range flowsGroup.dst {

				// Deletion during iteration is safe,
				// see https://stackoverflow.com/a/23230406 and later clarifications

				if row.empty() {
					delete(flowsGroup.dst, dstKey)
				} else {
					row.clearCounters()
				}
			}

			// delete flows group when empty
			if len(flowsGroup.dst) == 0 {
				delete(sampler.flows, flowKey)
			}
		}

		// delete sampler when no more data for/from it
		if len(sampler.flows) == 0 {
			delete(s.samplers, samplerKey)
		}
	}
}

func (s *summary) Dump(maxCount int, lineBreak string) string {
	var sb strings.Builder

	for samplerKey, sampler := range s.samplers {

		var lines []*packets

		//get slice of lines
		for _, flowsGroup := range sampler.flows {
			for _, row := range flowsGroup.dst {
				if row.empty() {
					continue
				}
				lines = append(lines, row)
			}
		}

		if len(lines) == 0 {
			continue
		}

		//print header
		if len(lines) <= maxCount {
			sb.WriteString(fmt.Sprintf("%s%s", samplerKey, lineBreak))
		} else {
			sb.WriteString(fmt.Sprintf("%s top %v of %v%s",
				samplerKey,
				maxCount,
				len(lines),
				lineBreak))
		}

		//sort in reverse order, from big to small
		sort.Slice(lines, func(a, b int) bool {
			return lines[a].Bytes() > lines[b].Bytes()
		})

		//cut to top count
		if len(lines) > maxCount {
			lines = lines[:maxCount]
		}

		//print lines
		for _, v := range lines {
			sb.WriteString(fmt.Sprintf("  %v, %+v%s", v.Description(), v.Value(), lineBreak))
		}
	}

	sb.WriteString(lineBreak)

	return sb.String()
}
