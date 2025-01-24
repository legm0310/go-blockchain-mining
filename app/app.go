package app

import (
	"blockchain-mining/global"
	"bufio"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
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

func NewApp(config *config.Config, difficulty int64) {
	a := &App{
		config: config,
		log:    log15.New("module", "app"),
	}

	var err error

	if a.repository, err = repository.NewRepository(a.config); err != nil {
		panic(err)
	}

	a.service = service.NewService(a.config, a.repository, difficulty)

	a.log.Info("Module Started", "time", time.Now().Unix(), "difficulty", difficulty)

	sc := bufio.NewScanner(os.Stdin)

	useCase()

	for {

		from := global.FROM()

		if from != "" {
			a.log.Info("Current Connected Wallet", "from", from)
			fmt.Println()
		}

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
		from := global.FROM()

		switch input[0] {
		case TransferCoin:
			a.service.CreateBlock([]*Transaction{}, []byte{}, 0)

		case MintCoin:
			fmt.Println("MintCoin in Switch")

		case CreateWallet:
			if wallet := a.service.MakeWallet(); wallet == nil {
				panic("Failed To Create Wallet")
			} else {
				fmt.Println()
				a.log.Info("Success To Create Wallet", "PrivateKey", wallet.PrivateKey, "PublicKey", wallet.PublicKey)
				fmt.Println()
			}

		case ConnectWallet:
			if from != "" {
				a.log.Debug("Already Connected", "from", from)
				fmt.Println()
			} else {
				if wallet, err := a.service.GetWallet(input[1]); err != nil {
					if err == mongo.ErrNoDocuments {
						a.log.Debug("Failed To Find Wallet PK is Nil", "PrivateKey", input[1])
					} else {
						a.log.Crit("Failed To Find Wallet", "PrivateKey", input[1], "err", err)
					}
					fmt.Println()
				} else {
					global.SetFrom(wallet.PublicKey)
					fmt.Println()
					a.log.Info("Success To Connect Wallet", "from", wallet.PublicKey)
					fmt.Println()
				}
			}

		case ChangeWallet:
			if from == "" {
				a.log.Debug("Connect Wallet First")
				fmt.Println()
			} else {
				if wallet, err := a.service.GetWallet(input[1]); err != nil {
					if err == mongo.ErrNoDocuments {
						a.log.Debug("Failed To Find Wallet PK is Nil", "PrivateKey", input[1])
					} else {
						a.log.Crit("Failed To Find Wallet", "PrivateKey", input[1], "err", err)
					}
				} else {
					global.SetFrom(wallet.PublicKey)
					fmt.Println()
					a.log.Info("Success To Change Wallet", "from", wallet.PublicKey)
					fmt.Println()
				}
			}
		case "":
			fmt.Println()
		default:
			return errors.New("failed to find CLI order")
		}
	}
	return nil
}

func useCase() {
	fmt.Println()

	fmt.Println("This is Akaps Module For BlockChain Core With Mongo")
	fmt.Println()
	fmt.Println("Use Case")

	fmt.Println("1. ", CreateWallet)
	fmt.Println("2. ", ConnectWallet, " <PK>")
	fmt.Println("3. ", ChangeWallet, " <PK>")
	fmt.Println("4. ", TransferCoin, " <To> <Amount>")
	fmt.Println("5. ", MintCoin, " <To> <Amount>")
	fmt.Println()
}
