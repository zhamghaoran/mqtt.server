package packets

import (
	"github.com/zhamghaoran/mqtt.server/constant"
	"io"
)

type PingrespPacket struct {
	FixedHeader
}

func (pr *PingrespPacket) Type() int {
	return constant.PINGRESQ
}

func (pr *PingrespPacket) String() string {
	return pr.FixedHeader.String()
}

func (pr *PingrespPacket) Write(w io.Writer) error {
	packet := pr.FixedHeader.pack()
	_, err := packet.WriteTo(w)

	return err
}

func (pr *PingrespPacket) Unpack(b io.Reader) error {
	return nil
}

func (pr *PingrespPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0, Address: pr.RemoteAddress}
}
