package main

import (
	"bytes"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"leetcode/constant"
	"leetcode/handler"
	. "leetcode/packet"
	"log"
	"net"
)

func main() {
	// 创建一个tcp服务
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":1883")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("MQTT server listening on localhost:1883")

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	for {
		log.Println("New client connected:", conn.RemoteAddr())
		packet, err := ReadPacket(conn)
		if err != nil {
			break
		}
		log.Println(packet.String())
		// 处理数据
		typeCode, err := handleDeclaredStruct(packet)
		if err != nil {
			return
		}
		if packet.Type() == Connect {
			// 这是一个连接请求，保存连接
			handler.SetConn(conn, packet.Details().MessageID)
		}
		if packet.Type() == Disconnect {
			// 取消连接
			handler.DeleteConn(packet.Details().MessageID)
			sendACK(conn, typeCode)
			break
		}
		sendACK(conn, typeCode)
	}
	log.Println("Client disconnected:", conn.RemoteAddr())
}

func handleDeclaredStruct(packet ControlPacket) (int, error) {
	// 获取到方法列表
	handlers := handler.GetHandler()
	// 将方法列表传入
	typeCode, err := ExecuteHandler(packet, handlers)
	if err != nil {
		return typeCode, err
	}
	return typeCode, nil
}
func ExecuteHandler(packet ControlPacket, handler handler.HandlerI) (int, error) {
	typeCode := packet.Type()
	var err error
	switch typeCode {
	case constant.CONNECT:
		return packet.(*ConnectPacket).Type(), handler.ConnectHandle(packet.(*ConnectPacket))
	case constant.CONNACK:
		packet = packet.(*ConnackPacket)
		return packet.Type(), nil
	case constant.PUBLISH:
		packet = packet.(*PublishPacket)
		return packet.Type(), nil
	case constant.PUBACK:
		packet = packet.(*PubackPacket)
		return packet.Type(), nil
	case constant.PUBREC:
		packet = packet.(*PubrecPacket)
		return packet.Type(), nil
	case constant.PUBREL:
		packet = packet.(*PubrelPacket)
		return packet.Type(), nil
	case constant.PUBCOMP:
		packet = packet.(*PubcompPacket)
		return packet.Type(), nil
	case constant.SUBSCRIBE:
		packet = packet.(*SubscribePacket)
		return packet.Type(), handler.SubscribeHandle(packet.(*SubscribePacket))
	case constant.SUBACK:
		packet = packet.(*SubackPacket)
		return packet.Type(), nil
	case constant.UNSUBSCRIBE:
		packet = packet.(*UnsubscribePacket)
		return packet.Type(), nil
	case constant.UNSUBSCRIBEACK:
		packet = packet.(*UnsubackPacket)
		return packet.Type(), nil
	case constant.PINGREG:
		packet = packet.(*PingreqPacket)
		return packet.Type(), nil
	case constant.PINGRESQ:
		packet = packet.(*PingrespPacket)
		return packet.Type(), nil
	case constant.DISCONNECT:
		packet = packet.(*DisconnectPacket)
		return packet.Type(), nil
	default:
		err = fmt.Errorf("unsupported packet type : %d", typeCode)
		return 0, err
	}

}
func sendACK(conn net.Conn, messageType int) {
	var err error
	var i bytes.Buffer
	switch messageType {
	case constant.CONNECT:
		i.Reset()
		connackPacket := packets.NewControlPacket(Connack).(*packets.ConnackPacket)
		_ = connackPacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	case constant.PINGREG:
		i.Reset()
		pingrespPacket := packets.NewControlPacket(Pingresp).(*packets.PingrespPacket)
		_ = pingrespPacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	case constant.SUBSCRIBE:
		i.Reset()
		subackpacket := packets.NewControlPacket(Suback).(*packets.SubackPacket)
		_ = subackpacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	}
	if err != nil {
		log.Printf("sendACK err : %s", err.Error())
	}
}
