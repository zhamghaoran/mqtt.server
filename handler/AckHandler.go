package handler

import (
	"bytes"
	"fmt"
	"leetcode/constant"
	packets "leetcode/packet"
	"net"
)

func SendACK(conn net.Conn, messageType int, id uint16) error {
	var i bytes.Buffer
	factory, err := AckFactory(messageType, id)
	if err != nil {
		return err
	}
	err = factory.Write(&i)
	if err != nil {
		return err
	}
	_, err = conn.Write(i.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func AckFactory(messageType int, messageID uint16) (packets.ControlPacket, error) {
	switch messageType {
	case constant.CONNECT:
		conn, _ := packets.NewControlPacket(byte(constant.CONNACK)).(*packets.ConnackPacket)
		return conn, nil
	case constant.PUBLISH:
		conn, _ := packets.NewControlPacket(byte(constant.PUBACK)).(*packets.PubackPacket)
		conn.MessageID = messageID
		return conn, nil
	case constant.PINGREG:
		conn, _ := packets.NewControlPacket(byte(constant.PINGRESQ)).(*packets.PingrespPacket)
		return conn, nil
	case constant.SUBSCRIBE:
		conn, _ := packets.NewControlPacket(byte(constant.SUBACK)).(*packets.SubackPacket)
		conn.MessageID = messageID
		return conn, nil
	case constant.UNSUBSCRIBE:
		conn, _ := packets.NewControlPacket(byte(constant.UNSUBSCRIBEACK)).(*packets.UnsubackPacket)
		conn.MessageID = messageID
		return conn, nil
	}
	return nil, fmt.Errorf("生成ACK消息失败")

}
