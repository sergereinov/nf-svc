package summary

import (
	"fmt"
	"net"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type message struct {
	*flowmessage.FlowMessage
}

func (m message) description(isAnswer bool) string {
	var port string
	if isAnswer {
		port = fmt.Sprintf("SrcPort=%v", m.SrcPort)
	} else {
		port = fmt.Sprintf("DstPort=%v", m.DstPort)
	}

	return fmt.Sprintf(
		"L3=%v, L4=%v, Src=%v, Dst=%v, %v",
		Ethertype(m.Etype),
		L4Proto(m.Proto),
		net.IP(m.SrcAddr),
		net.IP(m.DstAddr),
		port,
	)
}
