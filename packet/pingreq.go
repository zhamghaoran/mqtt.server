package packets

import (
	"io"
	"leetcode/constant"
)

type PingreqPacket struct {
	FixedHeader
}

func (pr *PingreqPacket) Type() int {
	return constant.PINGREG
}

func (pr *PingreqPacket) String() string {
	return pr.FixedHeader.String()
}

func (pr *PingreqPacket) Write(w io.Writer) error {
	packet := pr.FixedHeader.pack()
	_, err := packet.WriteTo(w)

	return err
}

func (pr *PingreqPacket) Unpack(b io.Reader) error {
	return nil
}

func (pr *PingreqPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0, Address: pr.RemoteAddress}
}
