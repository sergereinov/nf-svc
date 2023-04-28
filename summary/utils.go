package summary

import "net"

func isMulticast(addr []byte) bool {
	ip := net.IP(addr)
	return ip.Equal(net.IPv4bcast) ||
		ip.IsMulticast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() ||
		(len(ip) == net.IPv4len && ip[3] == 0xff)
}
