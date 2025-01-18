package app

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"blockchain-mining/config"
	"blockchain-mining/repository"
	"blockchain-mining/service"

	. "blockchain-mining/types"

	"github.com/inconshreveable/log15"
)

type App struct {
	config *config.Config

	service    *service.Service
	repository *repository.Repository

	log log15.Logger
}

func NewApp(config *config.Config) {
	a := &App{
		config: config,
		log:    log15.New("module", "app"),
	}

	var err error

	if a.repository, err = repository.NewRepository(a.config); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(a.config, a.repository)

		a.log.Info("Module Started", "time", time.Now().Unix())

		sc := bufio.NewScanner(os.Stdin)

		useCase()

		for {
			sc.Scan()
			fmt.Println(sc.Text())
		}
	}
}

func useCase() {
	fmt.Println()

	fmt.Println("This is Akaps Module For BlockChain Core With Mongo")
	fmt.Println()
	fmt.Println("Use Case")

	fmt.Println("1. ", CreateWallet)
	fmt.Println("2. ", TransferCoin, " <To> <Amount>")
	fmt.Println("3. ", MintCoin, " <To> <Amount>")
	fmt.Println()
}
