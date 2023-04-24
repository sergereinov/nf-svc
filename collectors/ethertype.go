package collectors

import "fmt"

// ref: https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml

var ethertype = map[uint32]string{
	2048:   "IPv4",
	0x86DD: "IPv6",
}

func Ethertype(etype uint32) string {
	if name, ok := ethertype[etype]; ok {
		return name
	}

	return fmt.Sprint(etype)
}
