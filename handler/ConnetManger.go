package handler

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	ExpireTime = time.Minute * 5
	DeleteTime = time.Millisecond * 100
)

type connection struct {
	lock   sync.Mutex
	conn   net.Conn
	expire time.Duration
}

var connMap map[uint16]*connection

func SetConn(conn net.Conn, messageId uint16) {
	c := &connection{}
	c.conn = conn
	c.expire = ExpireTime
	connMap[messageId] = c
}
func GetConn(messageId uint16) (net.Conn, error) {
	if connMap[messageId] == nil {
		return nil, fmt.Errorf("未定义的连接: %d", messageId)
	}
	connMap[messageId].lock.Lock()
	conn := connMap[messageId].conn
	connMap[messageId].expire = ExpireTime
	connMap[messageId].lock.Unlock()
	return conn, nil
}
func DeleteConn(messageId uint16) {
	connMap[messageId] = nil
}
func init() {
	scanMapDeleteExpireKey()
}
func scanMapDeleteExpireKey() {
	for k, v := range connMap {
		v.lock.Lock()
		if v.expire < DeleteTime {
			delete(connMap, k)
			v.lock.Unlock()
		} else {
			v.expire -= DeleteTime
			v.lock.Unlock()
		}
	}
	time.Sleep(DeleteTime)
}
