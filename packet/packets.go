package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ControlPacket interface {
	Write(io.Writer) error
	Unpack(io.Reader) error
	String() string
	Details() Details
	Type() int
}

var PacketNames = map[uint8]string{
	1:  "CONNECT",
	2:  "CONNACK",
	3:  "PUBLISH",
	4:  "PUBACK",
	5:  "PUBREC",
	6:  "PUBREL",
	7:  "PUBCOMP",
	8:  "SUBSCRIBE",
	9:  "SUBACK",
	10: "UNSUBSCRIBE",
	11: "UNSUBACK",
	12: "PINGREQ",
	13: "PINGRESP",
	14: "DISCONNECT",
}

const (
	Connect     = 1
	Connack     = 2
	Publish     = 3
	Puback      = 4
	Pubrec      = 5
	Pubrel      = 6
	Pubcomp     = 7
	Subscribe   = 8
	Suback      = 9
	Unsubscribe = 10
	Unsuback    = 11
	Pingreq     = 12
	Pingresp    = 13
	Disconnect  = 14
)

const (
	Accepted                        = 0x00
	ErrRefusedBadProtocolVersion    = 0x01
	ErrRefusedIDRejected            = 0x02
	ErrRefusedServerUnavailable     = 0x03
	ErrRefusedBadUsernameOrPassword = 0x04
	ErrRefusedNotAuthorised         = 0x05
	ErrNetworkError                 = 0xFE
	ErrProtocolViolation            = 0xFF
)

var ConnackReturnCodes = map[uint8]string{
	0:   "Connection Accepted",
	1:   "Connection Refused: Bad Protocol Version",
	2:   "Connection Refused: Client Identifier Rejected",
	3:   "Connection Refused: Server Unavailable",
	4:   "Connection Refused: Username or Password in unknown format",
	5:   "Connection Refused: Not Authorised",
	254: "Connection Error",
	255: "Connection Refused: Protocol Violation",
}

var (
	ErrorRefusedBadProtocolVersion    = errors.New("unacceptable protocol version")
	ErrorRefusedIDRejected            = errors.New("identifier rejected")
	ErrorRefusedServerUnavailable     = errors.New("server Unavailable")
	ErrorRefusedBadUsernameOrPassword = errors.New("bad user name or password")
	ErrorRefusedNotAuthorised         = errors.New("not Authorized")
	ErrorNetworkError                 = errors.New("network Error")
	ErrorProtocolViolation            = errors.New("protocol Violation")
)

var ConnErrors = map[byte]error{
	Accepted:                        nil,
	ErrRefusedBadProtocolVersion:    ErrorRefusedBadProtocolVersion,
	ErrRefusedIDRejected:            ErrorRefusedIDRejected,
	ErrRefusedServerUnavailable:     ErrorRefusedServerUnavailable,
	ErrRefusedBadUsernameOrPassword: ErrorRefusedBadUsernameOrPassword,
	ErrRefusedNotAuthorised:         ErrorRefusedNotAuthorised,
	ErrNetworkError:                 ErrorNetworkError,
	ErrProtocolViolation:            ErrorProtocolViolation,
}

func ReadPacket(r io.Reader, remoteAddress string) (ControlPacket, error) {
	var fh FixedHeader
	b := make([]byte, 1)

	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	err = fh.unpack(b[0], r)
	if err != nil {
		return nil, err
	}
	fh.RemoteAddress = remoteAddress
	cp, err := NewControlPacketWithHeader(fh)
	if err != nil {
		return nil, err
	}

	packetBytes := make([]byte, fh.RemainingLength)
	n, err := io.ReadFull(r, packetBytes)
	if err != nil {
		return nil, err
	}
	if n != fh.RemainingLength {
		return nil, errors.New("failed to read expected data")
	}

	err = cp.Unpack(bytes.NewBuffer(packetBytes))
	return cp, err
}

func NewControlPacket(packetType byte) ControlPacket {
	switch packetType {
	case Connect:
		return &ConnectPacket{FixedHeader: FixedHeader{MessageType: Connect}}
	case Connack:
		return &ConnackPacket{FixedHeader: FixedHeader{MessageType: Connack}}
	case Disconnect:
		return &DisconnectPacket{FixedHeader: FixedHeader{MessageType: Disconnect}}
	case Publish:
		return &PublishPacket{FixedHeader: FixedHeader{MessageType: Publish}}
	case Puback:
		return &PubackPacket{FixedHeader: FixedHeader{MessageType: Puback}}
	case Pubrec:
		return &PubrecPacket{FixedHeader: FixedHeader{MessageType: Pubrec}}
	case Pubrel:
		return &PubrelPacket{FixedHeader: FixedHeader{MessageType: Pubrel, Qos: 1}}
	case Pubcomp:
		return &PubcompPacket{FixedHeader: FixedHeader{MessageType: Pubcomp}}
	case Subscribe:
		return &SubscribePacket{FixedHeader: FixedHeader{MessageType: Subscribe, Qos: 1}}
	case Suback:
		return &SubackPacket{FixedHeader: FixedHeader{MessageType: Suback}}
	case Unsubscribe:
		return &UnsubscribePacket{FixedHeader: FixedHeader{MessageType: Unsubscribe, Qos: 1}}
	case Unsuback:
		return &UnsubackPacket{FixedHeader: FixedHeader{MessageType: Unsuback}}
	case Pingreq:
		return &PingreqPacket{FixedHeader: FixedHeader{MessageType: Pingreq}}
	case Pingresp:
		return &PingrespPacket{FixedHeader: FixedHeader{MessageType: Pingresp}}
	}
	return nil
}

func NewControlPacketWithHeader(fh FixedHeader) (ControlPacket, error) {
	switch fh.MessageType {
	case Connect:
		return &ConnectPacket{FixedHeader: fh}, nil
	case Connack:
		return &ConnackPacket{FixedHeader: fh}, nil
	case Disconnect:
		return &DisconnectPacket{FixedHeader: fh}, nil
	case Publish:
		return &PublishPacket{FixedHeader: fh}, nil
	case Puback:
		return &PubackPacket{FixedHeader: fh}, nil
	case Pubrec:
		return &PubrecPacket{FixedHeader: fh}, nil
	case Pubrel:
		return &PubrelPacket{FixedHeader: fh}, nil
	case Pubcomp:
		return &PubcompPacket{FixedHeader: fh}, nil
	case Subscribe:
		return &SubscribePacket{FixedHeader: fh}, nil
	case Suback:
		return &SubackPacket{FixedHeader: fh}, nil
	case Unsubscribe:
		return &UnsubscribePacket{FixedHeader: fh}, nil
	case Unsuback:
		return &UnsubackPacket{FixedHeader: fh}, nil
	case Pingreq:
		return &PingreqPacket{FixedHeader: fh}, nil
	case Pingresp:
		return &PingrespPacket{FixedHeader: fh}, nil
	}
	return nil, fmt.Errorf("unsupported packet type 0x%x", fh.MessageType)
}

type Details struct {
	Qos       byte
	MessageID uint16
	Address   string
}

type FixedHeader struct {
	MessageType     byte
	Dup             bool
	Qos             byte
	Retain          bool
	RemainingLength int
	RemoteAddress   string
}

func (fh FixedHeader) String() string {
	return fmt.Sprintf("%s: dup: %t qos: %d retain: %t rLength: %d", PacketNames[fh.MessageType], fh.Dup, fh.Qos, fh.Retain, fh.RemainingLength)
}

func boolToByte(b bool) byte {
	switch b {
	case true:
		return 1
	default:
		return 0
	}
}

func (fh *FixedHeader) pack() bytes.Buffer {
	var header bytes.Buffer
	header.WriteByte(fh.MessageType<<4 | boolToByte(fh.Dup)<<3 | fh.Qos<<1 | boolToByte(fh.Retain))
	header.Write(encodeLength(fh.RemainingLength))
	return header
}

func (fh *FixedHeader) unpack(typeAndFlags byte, r io.Reader) error {
	fh.MessageType = typeAndFlags >> 4
	fh.Dup = (typeAndFlags>>3)&0x01 > 0
	fh.Qos = (typeAndFlags >> 1) & 0x03
	fh.Retain = typeAndFlags&0x01 > 0

	var err error
	fh.RemainingLength, err = decodeLength(r)
	return err
}

func decodeByte(b io.Reader) (byte, error) {
	num := make([]byte, 1)
	_, err := b.Read(num)
	if err != nil {
		return 0, err
	}

	return num[0], nil
}

func decodeUint16(b io.Reader) (uint16, error) {
	num := make([]byte, 2)
	_, err := b.Read(num)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(num), nil
}

func encodeUint16(num uint16) []byte {
	bytesResult := make([]byte, 2)
	binary.BigEndian.PutUint16(bytesResult, num)
	return bytesResult
}

func encodeString(field string) []byte {
	return encodeBytes([]byte(field))
}

func decodeString(b io.Reader) (string, error) {
	buf, err := decodeBytes(b)
	return string(buf), err
}

func decodeBytes(b io.Reader) ([]byte, error) {
	fieldLength, err := decodeUint16(b)
	if err != nil {
		return nil, err
	}

	field := make([]byte, fieldLength)
	_, err = b.Read(field)
	if err != nil {
		return nil, err
	}

	return field, nil
}

func encodeBytes(field []byte) []byte {
	fieldLength := make([]byte, 2)
	binary.BigEndian.PutUint16(fieldLength, uint16(len(field)))
	return append(fieldLength, field...)
}

func encodeLength(length int) []byte {
	var encLength []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encLength = append(encLength, digit)
		if length == 0 {
			break
		}
	}
	return encLength
}

func decodeLength(r io.Reader) (int, error) {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	for multiplier < 27 {
		_, err := io.ReadFull(r, b)
		if err != nil {
			return 0, err
		}

		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength), nil
}
