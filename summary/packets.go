package summary

import "fmt"

type packets struct {
	description string
	bytes       uint64
	packets     uint64
}

func (p *packets) Description() string {
	return p.description
}

func (p *packets) Value() string {
	return fmt.Sprintf("{Bytes:%d Packets:%d}", p.Bytes(), p.Packets())
}

func (p *packets) Bytes() uint64 {
	return p.bytes
}

func (p *packets) Packets() uint64 {
	return p.packets
}

func (p *packets) clearCounters() {
	p.bytes = 0
	p.packets = 0
}

func (p *packets) empty() bool {
	return p.bytes == 0
}
