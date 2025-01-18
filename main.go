package main

import (
	"flag"

	"blockchain-mining/app"
	"blockchain-mining/config"
)

var (
	configFlag = flag.String("env", "./env.toml", "env.toml file not found")
	difficulty = flag.Int("diff", 12, "difficulty err")
)

func main() {
	flag.Parse()
	c := config.NewConfig(*configFlag)

	app.NewApp(c)
}
