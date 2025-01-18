package main

import (
	"blockchain-mining/config"
	"flag"
	"fmt"
)

var (
	configFlag = flag.String("env", "./env.toml", "env.toml file not found")
	difficulty = flag.Int("diff", 12, "difficulty err")
)

func main() {
	flag.Parse()

	c := config.NewConfig(*configFlag)
	fmt.Println(c.Info.Version)
	fmt.Println("test")
}
