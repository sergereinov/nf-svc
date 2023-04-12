package main

import "fmt"

// ref: https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers

var l4proto = map[uint32]string{
	1:  "ICMP",
	2:  "IGMP",
	6:  "TCP",
	17: "UDP",
}

func L4Proto(proto uint32) string {
	if name, ok := l4proto[proto]; ok {
		return name
	}
	return fmt.Sprint(proto)
}
