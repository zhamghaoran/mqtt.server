package packets

import (
	"fmt"
	"github.com/zhamghaoran/mqtt.server/constant"
	"io"
)

type PubrecPacket struct {
	FixedHeader
	MessageID uint16
}

func (pr *PubrecPacket) Type() int {
	return constant.PUBREC
}

func (pr *PubrecPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", pr.FixedHeader, pr.MessageID)
}

func (pr *PubrecPacket) Write(w io.Writer) error {
	var err error
	pr.FixedHeader.RemainingLength = 2
	packet := pr.FixedHeader.pack()
	packet.Write(encodeUint16(pr.MessageID))
	_, err = packet.WriteTo(w)

	return err
}

func (pr *PubrecPacket) Unpack(b io.Reader) error {
	var err error
	pr.MessageID, err = decodeUint16(b)

	return err
}

func (pr *PubrecPacket) Details() Details {
	return Details{Qos: pr.Qos, MessageID: pr.MessageID, Address: pr.RemoteAddress}
}
