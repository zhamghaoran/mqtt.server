package main

import (
	"fmt"
	packets "leetcode/packet"
	"log"
	"net"
)

const (
	Reversed       int = 0  //保留
	CONNECT        int = 1  // 客户端请求连接服务端
	CONNACK        int = 2  //连接确认报文确定
	PUBLISH        int = 3  //发布消息
	PUBACK         int = 4  //Qos1消息发布确认
	PUBREC         int = 5  //发布收到保证第一步
	PUBREL         int = 6  //发布释放 保证交付第二部
	PUBCOMP        int = 7  //Qos2消息发布完成  保证交付第二部
	SUBSCRIBE      int = 8  //客户端订阅发布
	SUBACK         int = 9  //订阅请求报文确定
	UNSUBSCRIBE    int = 10 //客户端取消订阅
	UNSUBSCRIBEACK int = 11 //客户端取消订阅请求
	PINGREG        int = 12 //心跳请求
	PINGRESQ       int = 13 //心跳响应
	DISCONNECT     int = 14 //客户端断开连接
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
	defer conn.Close()
	for {
		log.Println("New client connected:", conn.RemoteAddr())
		packet, err := packets.ReadPacket(conn)
		if err != nil {
			break
		}
		println(packet.String())
		packet.Details()
	}
	log.Println("Client disconnected:", conn.RemoteAddr())
}
func sendACK(conn net.Conn, messageType int) {
	var err error
	switch messageType {
	case CONNECT:
		_, err = conn.Write([]byte{0x20, byte(CONNACK), 0x00, 0x00})
	case PINGREG:
		_, err = conn.Write([]byte{byte(PINGRESQ), 0x00})
	}
	if err != nil {
		log.Printf("sendACK err : %s", err.Error())
	}
}

type MQTTHeader struct {
	MessageType byte
	Dup         bool
	QosLevel    byte
	Retain      bool
	Remaining   int
}

func (header *MQTTHeader) ParseHeader(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid header length")
	}

	header.MessageType = data[0] >> 4
	header.Dup = ((data[0] >> 3) & 0x01) == 1
	header.QosLevel = (data[0] >> 1) & 0x03
	header.Retain = (data[0] & 0x01) == 1

	var multiplier uint32 = 1
	var value uint32 = 0
	var pos int = 1
	var b byte

	for {
		if pos >= len(data) {
			return fmt.Errorf("invalid header length")
		}

		b = data[pos]
		value += uint32(b&0x7F) * multiplier
		multiplier *= 128

		if multiplier > 128*128*128 {
			return fmt.Errorf("invalid header length")
		}
		pos++
		if b&0x80 == 0 {
			break
		}
	}
	header.Remaining = int(value)

	return nil
}
