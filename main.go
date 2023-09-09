package main

import (
	config2 "leetcode/config"
	"leetcode/service"
)

func main() {
	config := config2.Config{Port: "1883"}
	service.CreateService(config)
}
