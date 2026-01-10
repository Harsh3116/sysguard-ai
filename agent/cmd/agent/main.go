package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("SysGuard AI Agent started")

	for {
		fmt.Println("Agent heartbeat: system running")
		time.Sleep(10 * time.Second)
	}
}
