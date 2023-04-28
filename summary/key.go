package summary

import (
	"fmt"
	"net"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

const (
	_MULTICAST = "multicast"
)

type key struct {
	*flowmessage.FlowMessage
}

// key (and text) for this flow's sampler
func (k key) sampler() string {
	return fmt.Sprintf(
		"%v(%v)",
		k.Type,
		net.IP(k.SamplerAddress),
	)
}

func (k key) destination() string {
	return fmt.Sprintf("%v", net.IP(k.DstAddr))
}

// main key for this flow in forward direction
func (k key) asToDst() string {
	if isMulticast(k.DstAddr) {
		return k.format(net.IP(k.SrcAddr).String(), _MULTICAST, k.DstPort)
	}
	return k.format(net.IP(k.SrcAddr).String(), net.IP(k.DstAddr).String(), k.DstPort)
}

// main key for this flow in backward direction
func (k key) asFromSrc() string {
	return k.format(net.IP(k.SrcAddr).String(), net.IP(k.DstAddr).String(), k.SrcPort)
}

// probable main key for the flow as if this flow were a response to the main flow
func (k key) parent() string {
	return k.format(net.IP(k.DstAddr).String(), net.IP(k.SrcAddr).String(), k.SrcPort)
}

// probable main key for the multicast flow as if this flow were a response to the main multicast flow
func (k key) parentMulticast() string {
	return k.format(net.IP(k.DstAddr).String(), _MULTICAST, k.SrcPort)
}

// format the key as "L3,L4,Src,Dst,Port"
func (k key) format(src, dst string, port uint32) string {
	return fmt.Sprintf(
		"%d,%d,%s,%s,%d",
		k.Etype,
		k.Proto,
		src,
		dst,
		port,
	)
}
