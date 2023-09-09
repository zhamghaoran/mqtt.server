package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mqtt-server/config"
	"mqtt-server/constant"
	"mqtt-server/handler"
	packets "mqtt-server/packet"
	"net"
)

func CreateService(Config config.Config) {
	// 创建一个tcp服务
	tcpAddr, _ := net.ResolveTCPAddr("tcp", Config.Port)
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
	defer conn.Close()
	for {
		log.Println("New client connected:", conn.RemoteAddr())
		packet, err := packets.ReadPacket(conn, conn.RemoteAddr().String())
		if err != nil && err != io.EOF {
			break
		}
		log.Println(packet.String())

		if packet.Type() == packets.Disconnect {
			// 取消连接
			handler.DeleteConn(packet.Details().Address)
			sendACK(conn, packet.Type(), 0)
			break
		}

		if packet.Type() == packets.Connect {
			// 这是一个连接请求，保存连接
			handler.SetConn(conn, packet.Details().Address)
			sendACK(conn, packet.Type(), 0)
		} else {
			// 登录状态校验
			err := handler.StateVerification(conn.RemoteAddr().String())
			if err != nil {
				log.Println(err.Error())
				return
			}
			// 处理数据
			typeCode, err := handleDeclaredStruct(packet)
			if err != nil {
				log.Printf("未定义的连接: %v", err.Error())
				return
			}
			err = handler.SendACK(conn, typeCode, packet.Details().MessageID)
			if err != nil {
				log.Println(err.Error())
			}
			//sendACK(conn, typeCode, packet.Details().MessageID)
		}
	}
	log.Println("Client disconnected:", conn.RemoteAddr())
}

func handleDeclaredStruct(packet packets.ControlPacket) (int, error) {
	// 获取到方法列表
	handlers := handler.GetHandler()
	// 将方法列表传入
	typeCode, err := ExecuteHandler(packet, handlers)
	if err != nil {
		return typeCode, err
	}
	return typeCode, nil
}
func ExecuteHandler(packet packets.ControlPacket, handler handler.HandlerI) (int, error) {
	typeCode := packet.Type()
	var err error
	switch typeCode {
	case constant.CONNECT:
		return packet.(*packets.ConnectPacket).Type(), handler.ConnectHandle(packet.(*packets.ConnectPacket))
	case constant.CONNACK:
		packet = packet.(*packets.ConnackPacket)
		return packet.Type(), nil
	case constant.PUBLISH:
		packet = packet.(*packets.PublishPacket)
		return packet.Type(), handler.PublishHandle(packet.(*packets.PublishPacket))
	case constant.PUBACK:
		packet = packet.(*packets.PubackPacket)
		return packet.Type(), nil
	case constant.PUBREC:
		packet = packet.(*packets.PubrecPacket)
		return packet.Type(), nil
	case constant.PUBREL:
		packet = packet.(*packets.PubrelPacket)
		return packet.Type(), nil
	case constant.PUBCOMP:
		packet = packet.(*packets.PubcompPacket)
		return packet.Type(), nil
	case constant.SUBSCRIBE:
		packet = packet.(*packets.SubscribePacket)
		return packet.Type(), handler.SubscribeHandle(packet.(*packets.SubscribePacket))
	case constant.SUBACK:
		packet = packet.(*packets.SubackPacket)
		return packet.Type(), nil
	case constant.UNSUBSCRIBE:
		packet = packet.(*packets.UnsubscribePacket)
		return packet.Type(), nil
	case constant.UNSUBSCRIBEACK:
		packet = packet.(*packets.UnsubackPacket)
		return packet.Type(), nil
	case constant.PINGREG:
		packet = packet.(*packets.PingreqPacket)
		return packet.Type(), nil
	case constant.PINGRESQ:
		packet = packet.(*packets.PingrespPacket)
		return packet.Type(), nil
	case constant.DISCONNECT:
		packet = packet.(*packets.DisconnectPacket)
		return packet.Type(), nil
	default:
		err = fmt.Errorf("unsupported packet type : %d", typeCode)
		return 0, err
	}

}
func sendACK(conn net.Conn, messageType int, id uint16) {
	var err error
	var i bytes.Buffer
	switch messageType {
	case constant.CONNECT:
		i.Reset()
		connackPacket := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
		_ = connackPacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	case constant.PINGREG:
		i.Reset()
		pingrespPacket := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)
		_ = pingrespPacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	case constant.SUBSCRIBE:
		i.Reset()
		subackpacket := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
		subackpacket.MessageID = id
		_ = subackpacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	case constant.PUBLISH:
		i.Reset()
		subackpacket := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
		subackpacket.MessageID = id
		_ = subackpacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	case constant.UNSUBSCRIBE:
		i.Reset()
		subackpacket := packets.NewControlPacket(packets.Unsuback).(*packets.UnsubackPacket)
		subackpacket.MessageID = id
		_ = subackpacket.Write(&i)
		fmt.Println(i.Bytes())
		_, _ = conn.Write(i.Bytes())
	}
	if err != nil {
		log.Printf("sendACK err : %s", err.Error())
	}
}
