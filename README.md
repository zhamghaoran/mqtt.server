# mqtt-server

## 介绍 

使用原生golang 构建的mqtt服务端，高效，简单

## 使用介绍

创建一个简单的mqtt 服务端

```go
func main() {
	Mqttconfig := config.Config{Port: "1883"}
	service.CreateService(Mqttconfig)
}
// 在1883端口上启动一个mqtt服务
```



