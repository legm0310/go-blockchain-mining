package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
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
	}

	a.service = service.NewService(a.config, a.repository, a.config.Info.Difficulty)

	a.log.Info("Module Started", "time", time.Now().Unix(), "difficulty", a.config.Info.Difficulty)

	sc := bufio.NewScanner(os.Stdin)

	useCase()

	for {
		sc.Scan()

		input := strings.Split(sc.Text(), " ")
		if err = a.inputValueAssessment(input); err != nil {
			a.log.Error("Failed to parse input", "err", err, "input", input)
		}
	}
}

func (a *App) inputValueAssessment(input []string) error {
	msg := errors.New("check Use Case")
	if len(input) == 0 {
		return msg
	} else {
		switch input[0] {
		case CreateWallet:
			fmt.Println("Create Wallet -----------------")
			if wallet := a.service.MakeWallet(); wallet == nil {
				panic("Failed To Create Wallet")
			} else {
				fmt.Println("Success To Create Wallet")
			}

		case TransferCoin:
			fmt.Println("TransferCoin in Switch")
		case MintCoin:
			fmt.Println("MintCoin in Switch")
		default:
			return errors.New("failed to find CLI order")
		}
		fmt.Println()
	}
	return nil
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
