package collectors

import (
	"fmt"
	"net"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type GroupMsg struct {
	Bytes   uint64
	Packets uint64
}

type Msg struct {
	*flowmessage.FlowMessage

	trackingClients []string
}

func (m Msg) Partition() string {
	return fmt.Sprintf(
		"%v(%v)",
		m.Type,
		net.IP(m.SamplerAddress),
	)
}

func isMulticast(addr []byte) bool {
	ip := net.IP(addr)
	return ip.Equal(net.IPv4bcast) ||
		ip.IsMulticast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() ||
		(len(ip) == net.IPv4len && ip[3] == 0xff)
}

func (m Msg) describeOppositePort() string {
	if isMulticast(m.SrcAddr) {
		//is it answer from multicast?
		//should we consider the SrcPort?
		return fmt.Sprintf("SrcPort=%v", m.SrcPort)
	}

	src := net.IP(m.SrcAddr).String()
	for _, v := range m.trackingClients {
		if v == src {
			//all packets from trackingClients treated as a response to a random port
			//so we should consider the SrcPort
			return fmt.Sprintf("SrcPort=%v", m.SrcPort)
		}
	}
	return fmt.Sprintf("DstPort=%v", m.DstPort)
}

func (m Msg) GroupKey() string {
	return fmt.Sprintf(
		"L3=%v, L4=%v, Src=%v, Dst=%v, %v",
		Ethertype(m.Etype),
		L4Proto(m.Proto),
		net.IP(m.SrcAddr),
		net.IP(m.DstAddr),
		m.describeOppositePort(),
	)
}

func (m Msg) Aggregate(acc GroupMsg) GroupMsg {
	acc.Bytes += m.Bytes
	acc.Packets += m.Packets
	return acc
}
