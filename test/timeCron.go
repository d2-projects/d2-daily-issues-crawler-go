package test

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second *3)
	for _ = range ticker.C {
		fmt.Println("time at %v", time.Now())
	}
}