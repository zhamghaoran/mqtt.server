package packets

import (
	"io"
	"leetcode/constant"
)

type DisconnectPacket struct {
	FixedHeader
}

func (d *DisconnectPacket) Type() int {
	return constant.DISCONNECT
}

func (d *DisconnectPacket) String() string {
	return d.FixedHeader.String()
}

func (d *DisconnectPacket) Write(w io.Writer) error {
	packet := d.FixedHeader.pack()
	_, err := packet.WriteTo(w)

	return err
}

func (d *DisconnectPacket) Unpack(b io.Reader) error {
	return nil
}

func (d *DisconnectPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}
