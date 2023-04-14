package collectors

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

const (
	_FIELDS_SEP = " "
	_VALUE_SEP  = ":"
)

func FormatFlowMessage(fmsg *flowmessage.FlowMessage) string {
	if fmsg == nil {
		return ""
	}

	srcmac := make([]byte, 8)
	dstmac := make([]byte, 8)
	binary.BigEndian.PutUint64(srcmac, fmsg.SrcMac)
	binary.BigEndian.PutUint64(dstmac, fmsg.DstMac)
	srcmac = srcmac[2:8]
	dstmac = dstmac[2:8]

	var buf []string
	buf = appendNonEmpty[string](buf, "Type", fmsg.Type.String())
	buf = appendNonEmpty[uint64](buf, "TimeRecv", fmsg.TimeReceived)
	buf = appendNonEmpty[uint32](buf, "SequenceNum", fmsg.SequenceNum)
	buf = appendNonEmpty[uint64](buf, "SamplingRate", fmsg.SamplingRate)
	buf = appendNonEmpty[string](buf, "Sampler", net.IP(fmsg.SamplerAddress).String())
	buf = appendNonEmpty[uint64](buf, "TimeFlowStart", fmsg.TimeFlowStart)
	buf = appendNonEmpty[uint64](buf, "TimeFlowEnd", fmsg.TimeFlowEnd)
	buf = appendNonEmpty[uint64](buf, "Bytes", fmsg.Bytes)
	buf = appendNonEmpty[uint64](buf, "Packets", fmsg.Packets)
	buf = appendNonEmpty[string](buf, "SrcAddr", net.IP(fmsg.SrcAddr).String())
	buf = appendNonEmpty[string](buf, "DstAddr", net.IP(fmsg.DstAddr).String())
	buf = appendNonEmpty[uint32](buf, "Etype", fmsg.Etype)
	buf = appendNonEmpty[uint32](buf, "Proto", fmsg.Proto)
	buf = appendNonEmpty[uint32](buf, "SrcPort", fmsg.SrcPort)
	buf = appendNonEmpty[uint32](buf, "DstPort", fmsg.DstPort)
	buf = appendNonEmpty[uint32](buf, "InIf", fmsg.InIf)
	buf = appendNonEmpty[uint32](buf, "OutIf", fmsg.OutIf)
	buf = appendNonEmpty[string](buf, "SrcMac", net.HardwareAddr(srcmac).String())
	buf = appendNonEmpty[string](buf, "DstMac", net.HardwareAddr(dstmac).String())
	buf = appendNonEmpty[uint32](buf, "SrcVlan", fmsg.SrcVlan)
	buf = appendNonEmpty[uint32](buf, "DstVlan", fmsg.DstVlan)
	buf = appendNonEmpty[uint32](buf, "VlanId", fmsg.VlanId)
	buf = appendNonEmpty[uint32](buf, "IngressVrfID", fmsg.IngressVrfID)
	buf = appendNonEmpty[uint32](buf, "EgressVrfID", fmsg.EgressVrfID)
	buf = appendNonEmpty[uint32](buf, "IPTos", fmsg.IPTos)
	buf = appendNonEmpty[uint32](buf, "ForwardingStatus", fmsg.ForwardingStatus)
	buf = appendNonEmpty[uint32](buf, "IPTTL", fmsg.IPTTL)
	buf = appendNonEmpty[uint32](buf, "TCPFlags", fmsg.TCPFlags)
	buf = appendNonEmpty[uint32](buf, "IcmpType", fmsg.IcmpType)
	buf = appendNonEmpty[uint32](buf, "IcmpCode", fmsg.IcmpCode)
	buf = appendNonEmpty[uint32](buf, "IPv6FlowLabel", fmsg.IPv6FlowLabel)
	buf = appendNonEmpty[uint32](buf, "FragmentId", fmsg.FragmentId)
	buf = appendNonEmpty[uint32](buf, "FragmentOffset", fmsg.FragmentOffset)
	buf = appendNonEmpty[uint32](buf, "BiFlowDirection", fmsg.BiFlowDirection)
	buf = appendNonEmpty[uint32](buf, "SrcAS", fmsg.SrcAS)
	buf = appendNonEmpty[uint32](buf, "DstAS", fmsg.DstAS)
	buf = appendNonEmpty[string](buf, "NextHop", net.IP(fmsg.NextHop).String(), "<nil>", "0.0.0.0")
	buf = appendNonEmpty[uint32](buf, "NextHopAS", fmsg.NextHopAS)
	buf = appendNonEmpty[uint32](buf, "SrcNet", fmsg.SrcNet)
	buf = appendNonEmpty[uint32](buf, "DstNet", fmsg.DstNet)

	buf = appendNonEmpty[bool](buf, "HasEncap", fmsg.HasEncap)
	buf = appendNonEmpty[string](buf, "SrcAddrEncap", net.IP(fmsg.SrcAddrEncap).String(), "<nil>", "0.0.0.0")
	buf = appendNonEmpty[string](buf, "DstAddrEncap", net.IP(fmsg.DstAddrEncap).String(), "<nil>", "0.0.0.0")
	buf = appendNonEmpty[uint32](buf, "ProtoEncap", fmsg.ProtoEncap)
	buf = appendNonEmpty[uint32](buf, "EtypeEncap", fmsg.EtypeEncap)
	buf = appendNonEmpty[uint32](buf, "IPTosEncap", fmsg.IPTosEncap)
	buf = appendNonEmpty[uint32](buf, "IPTTLEncap", fmsg.IPTTLEncap)
	buf = appendNonEmpty[uint32](buf, "IPv6FlowLabelEncap", fmsg.IPv6FlowLabelEncap)
	buf = appendNonEmpty[uint32](buf, "FragmentIdEncap", fmsg.FragmentIdEncap)
	buf = appendNonEmpty[uint32](buf, "FragmentOffsetEncap", fmsg.FragmentOffsetEncap)

	buf = appendNonEmpty[bool](buf, "HasMPLS", fmsg.HasMPLS)
	buf = appendNonEmpty[uint32](buf, "MPLSCount", fmsg.MPLSCount)
	buf = appendNonEmpty[uint32](buf, "MPLS1TTL", fmsg.MPLS1TTL)
	buf = appendNonEmpty[uint32](buf, "MPLS1Label", fmsg.MPLS1Label)
	buf = appendNonEmpty[uint32](buf, "MPLS2TTL", fmsg.MPLS2TTL)
	buf = appendNonEmpty[uint32](buf, "MPLS2Label", fmsg.MPLS2Label)
	buf = appendNonEmpty[uint32](buf, "MPLS3TTL", fmsg.MPLS3TTL)
	buf = appendNonEmpty[uint32](buf, "MPLS3Label", fmsg.MPLS3Label)
	buf = appendNonEmpty[uint32](buf, "MPLSLastTTL", fmsg.MPLSLastTTL)
	buf = appendNonEmpty[uint32](buf, "MPLSLastLabel", fmsg.MPLSLastLabel)

	buf = appendNonEmpty[bool](buf, "HasPPP", fmsg.HasPPP)
	buf = appendNonEmpty[uint32](buf, "PPPAddressControl", fmsg.PPPAddressControl)

	return strings.Join(buf, _FIELDS_SEP)
}

func appendNonEmpty[T comparable](buf []string, name string, value T, omit ...T) []string {
	for _, o := range omit {
		if value == o {
			return buf
		}
	}
	var zeroValue T
	if value == zeroValue {
		return buf
	}
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString(_VALUE_SEP)
	sb.WriteString(fmt.Sprintf("%v", value))
	return append(buf, sb.String())
}
