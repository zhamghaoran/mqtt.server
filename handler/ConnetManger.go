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
	conn   net.Conn
	expire time.Duration
}

var connMap sync.Map

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
