package main

import (
	"fmt"
	"time"

	"github.com/iamBelugaa/cache/cache"
)

func main() {
	c := cache.NewCache(
		time.Second*2,
		time.Second*3,
		&cache.Actions[string, string]{
			OnCleanup: func(key, val string) {
				fmt.Printf("%s -> %s", key, val)
				println()
			},
		},
	)
	defer c.Close()

	c.Set("1", "a")
	c.Set("2", "b")
	c.Set("3", "c")
	c.Set("4", "d")

	v, ok := c.Get("1")
	println(v, ok)

	c.Delete("1")
	c.Delete("3")
}
