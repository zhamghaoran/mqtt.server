package main

import (
	"github.com/zhamghaoran/mqtt.server/config"
	"github.com/zhamghaoran/mqtt.server/service"
)

func main() {
	service.CreateService(config.Config{Port: "1883"})
}
