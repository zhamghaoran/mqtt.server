package handler

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	ExpireTime = time.Minute * 5
	DeleteTime = time.Millisecond * 100
)

type connection struct {
	conn      net.Conn
	expire    time.Duration
	topicName []string
}

var connMap sync.Map
var subscribeMap sync.Map

func SetConn(conn net.Conn, messageId uint16) {
	c := &connection{}
	c.conn = conn
	c.expire = ExpireTime
	connMap.Store(messageId, c)
}
func GetConn(messageId uint16) (net.Conn, error) {
	conn, ok := connMap.LoadAndDelete(messageId)
	if ok {
		connect, _ := conn.(connection)
		connect.expire = ExpireTime
		connMap.Store(messageId, connect)
		return connect.conn, nil
	} else {
		return nil, fmt.Errorf("未定义的连接")
	}

}
func DeleteConn(messageId uint16) {
	connMap.Delete(messageId)
}

func init() {
	go scanMapDeleteExpireKey()
}
func scanMapDeleteExpireKey() {
	connMap.Range(func(key, value any) bool {
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
func SetConnName(messageId uint16, topicName string) error {
	value, ok := connMap.LoadAndDelete(messageId)
	if ok != true {
		return fmt.Errorf("连接异常")
	}
	conn, err := value.(connection)
	if err != true {
		return fmt.Errorf("转化异常")
	}
	conn.topicName = []string{topicName}
	connMap.Store(messageId, conn)
	return nil
}
func getConnectionByTopicName(topicName string) []connection {
	var connList []connection
	connMap.Range(func(key, value any) bool {
		conn, err := value.(connection)
		if err == true {
			connList = append(connList, conn)
		}
		return true
	})
	return connList
}

func publish(topicName string, payload []byte) error {
	load, o := subscribeMap.Load(topicName)
	if !o {
		return fmt.Errorf("获取连接失败")
	}
	connList, err := load.([]net.Conn)
	if !err {
		return fmt.Errorf("连接断言错误")
	}
	for _, v := range connList {
		_, err := v.Write(payload)
		if err != nil {
			return err
		}
	}
	return nil
}
func subscribe(messageId uint16, topicName []string) error {
	value, ok := connMap.Load(messageId)
	if !ok {
		return fmt.Errorf("获取连接错误")
	}
	conn, err := value.(connection)
	if !err {
		return fmt.Errorf("连接断言失败")
	}
	andDelete, loaded := subscribeMap.LoadAndDelete(conn.conn)
	if !loaded {
		subscribeMap.Store(conn.conn, topicName)
	} else {
		topicList, err := andDelete.([]string)
		if !err {
			return fmt.Errorf("参数列表断言错误")
		}
		topicList = append(topicList, topicName...)
		subscribeMap.Store(conn.conn, topicList)
	}
	return nil
}
