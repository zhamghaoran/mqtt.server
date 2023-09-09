package main

import (
	"leetcode/config"
	"leetcode/service"
)

func main() {
	Mqttconfig := config.Config{Port: "1883"}
	service.CreateService(Mqttconfig)
}
