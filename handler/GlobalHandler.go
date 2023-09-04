package handler

import (
	packets "leetcode/packet"
)

var HandlerMap map[string]interface{}

type HandlerI interface {
	ConnectHandle(packet *packets.ConnectPacket) error
	ConnectAckHandle(packet *packets.ConnackPacket) error
	PublishHandle(packet *packets.PublishPacket) error
	PubackHandle(packet *packets.PubackPacket) error
	PubrelHandle(packet *packets.PubrelPacket) error
	PubcompHandle(packet *packets.PubcompPacket) error
	SubscribeHandle(packet *packets.SubscribePacket) error
	SubackHandle(packet *packets.SubscribePacket) error
	UnsubscribeHandle(packet *packets.UnsubscribePacket) error
	UnsubackHandle(packet *packets.UnsubackPacket) error
	PingreqHandle(packet *packets.PingreqPacket) error
	PingrespHandle(packet *packets.PingrespPacket) error
	DisconnectHandle(packet *packets.DisconnectPacket) error
}
type DefaultHandler struct{}

func (DefaultHandler) ConnectHandle(packet *packets.ConnectPacket) error {
	//todo  密码校验
	return nil
}
func (DefaultHandler) ConnectAckHandle(packet *packets.ConnackPacket) error {
	return nil
}
func (DefaultHandler) PublishHandle(packet *packets.PublishPacket) error {
	return publish(packet.TopicName, packet.Payload)
}
func (DefaultHandler) PubackHandle(packet *packets.PubackPacket) error {
	return nil
}
func (DefaultHandler) PubrelHandle(packet *packets.PubrelPacket) error {
	return nil
}
func (DefaultHandler) PubcompHandle(packet *packets.PubcompPacket) error {
	return nil
}
func (DefaultHandler) SubscribeHandle(packet *packets.SubscribePacket) error {

	return nil
}
func (DefaultHandler) SubackHandle(packet *packets.SubscribePacket) error {
	return nil
}
func (DefaultHandler) UnsubscribeHandle(packet *packets.UnsubscribePacket) error {
	return nil
}
func (DefaultHandler) UnsubackHandle(packet *packets.UnsubackPacket) error {
	return nil
}
func (DefaultHandler) PingreqHandle(packet *packets.PingreqPacket) error {
	return nil
}
func (DefaultHandler) PingrespHandle(packet *packets.PingrespPacket) error {
	return nil
}
func (DefaultHandler) DisconnectHandle(packet *packets.DisconnectPacket) error {
	return nil
}
