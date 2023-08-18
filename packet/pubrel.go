package packets

import (
	"fmt"
	"io"
	"leetcode/constant"
)

type PubrelPacket struct {
	FixedHeader
	MessageID uint16
}

func (pr *PubrelPacket) Type() int {
	return constant.PUBREL
}

func (pr *PubrelPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", pr.FixedHeader, pr.MessageID)
}

func (pr *PubrelPacket) Write(w io.Writer) error {
	var err error
	pr.FixedHeader.RemainingLength = 2
	packet := pr.FixedHeader.pack()
	packet.Write(encodeUint16(pr.MessageID))
	_, err = packet.WriteTo(w)

	return err
}

func (pr *PubrelPacket) Unpack(b io.Reader) error {
	var err error
	pr.MessageID, err = decodeUint16(b)

	return err
}

func (pr *PubrelPacket) Details() Details {
	return Details{Qos: pr.Qos, MessageID: pr.MessageID}
}
