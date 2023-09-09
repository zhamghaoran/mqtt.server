package main

import (
	"mqtt/config"
	"mqtt/service"
)

func main() {
	Mqttconfig := config.Config{Port: "1883"}
	service.CreateService(Mqttconfig)
}
