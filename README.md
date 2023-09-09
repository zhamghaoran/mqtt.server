# mqtt-server

## 介绍 

使用原生golang 构建的mqtt服务端，高效，简单

## 使用介绍

- 添加依赖

`go get github.com/zhamghaoran/mqtt.server`

- 创建一个简单的mqtt 服务端

```go
func main() {
	Mqttconfig := config.Config{Port: "1883"}
	service.CreateService(Mqttconfig)
}
// 在1883端口上启动一个mqtt服务
```

## 使用自定义的handler

- 首先需要自己实现HandlerI 接口

  ```go
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
  ```

  然后调用SetHandler方法传入自定义的实现

  ```go
  func main() {
  	config := config2.Config{Port: "1883"}
  	handler.SetHandler(UserHandler)
  	service.CreateService(config)
  }
  ```

  

