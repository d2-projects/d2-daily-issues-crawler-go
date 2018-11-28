package test

import (
	"fmt"
	"github.com/robfig/cron"
)

func main () {
	i := 0
	c := cron.New()
	spec := "0 0/1 * * * ?"
	c.AddFunc(spec, func() {
		i++
		fmt.Println("count :", i)
	})
	c.Start()
	select{}
}
