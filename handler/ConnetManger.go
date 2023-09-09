package handler

import (
	"bytes"
	"fmt"
	"log"
	"mqtt-server/constant"
	packets "mqtt-server/packet"
	"net"
	"sync"
	"time"
)

const (
	ExpireTime = time.Minute * 5
	DeleteTime = time.Millisecond * 100
)

type connection struct {
	conn   net.Conn
	expire time.Duration
}

/*
key  remoteAddress
val  connection
*/
var connMap sync.Map

/*
key topic
value connection
*/
var subscribeMap sync.Map

func SetConn(conn net.Conn, remoteAdd string) {
	c := &connection{}
	c.conn = conn
	c.expire = ExpireTime
	connMap.Store(remoteAdd, c)
}

func DeleteConn(remoteAdd string) {
	value, _ := connMap.LoadAndDelete(remoteAdd)
	conn, _ := value.(connection)
	_ = conn.conn.Close()
}

func init() {
	go scanMapDeleteExpireKey()
}
func scanMapDeleteExpireKey() {
	connMap.Range(func(key, value interface{}) bool {
		connection, ok := value.(connection)
		if ok {
			if connection.expire < DeleteTime {
				connMap.Delete(key)
			} else {
				connection.expire -= DeleteTime
				connMap.Store(key, connection)
			}
		} else {
			log.Printf("存在异常连接%v", value)
		}
		return true
	})
	time.Sleep(DeleteTime)
}

func publish(topicName string, payload []byte) error {
	var i bytes.Buffer
	load, o := subscribeMap.Load(topicName)
	if !o {
		return fmt.Errorf("获取连接失败")
	}
	connList, err := load.([]net.Conn)
	if !err {
		return fmt.Errorf("连接断言错误")
	}
	for _, v := range connList {
		packet := packets.NewControlPacket(byte(constant.PUBLISH)).(*packets.PublishPacket)
		packet.Payload = payload
		packet.TopicName = topicName
		packet.Qos = 1
		_ = packet.Write(&i)
		_, sendError := v.Write(i.Bytes())
		if sendError != nil {
			return sendError
		}
	}
	return nil
}

func subscribe(remoteAdd string, topicName []string) error {
	value, ok := connMap.Load(remoteAdd)
	if !ok {
		return fmt.Errorf("获取连接错误")
	}
	conn, err := value.(*connection)
	if !err {
		return fmt.Errorf("连接断言失败")
	}
	for _, val := range topicName {
		andDelete, loaded := subscribeMap.LoadAndDelete(val)
		if !loaded {
			subscribeMap.Store(val, []net.Conn{conn.conn})
		} else {
			connList, err := andDelete.([]net.Conn)
			if !err {
				return fmt.Errorf("类型断言错误")
			}
			connList = append(connList, conn.conn)
			subscribeMap.Store(val, connList)
		}
	}
	return nil
}
func StateVerification(remoteAdd string) error {
	value, loaded := connMap.LoadAndDelete(remoteAdd)
	if !loaded {
		return fmt.Errorf("未定义的连接")
	}
	conn, err := value.(*connection)
	if !err {
		return fmt.Errorf("错误的连接")
	}
	conn.expire = ExpireTime
	connMap.Store(remoteAdd, conn)
	return nil
}
