package main

import (
	"fmt"
	"github.com/adlio/trello"
	"github.com/martinlindhe/inputbox"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	AppKey string `yaml:"appKey"`
	IdList string `yaml:"idList"`
}

func main() {
	service := "trello-fast-card"
	user := "trello-api-key"

	cfg := config{}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not resolve home dir %v", err)
	}
	configPath := filepath.Join(home, ".config", "trello-fast-card.yaml")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	trelloApiKey, err := keyring.Get(service, user)
	if err != nil {
		got, ok := inputbox.InputBox("trello api key", "input trello api key", "")
		if ok {
			err := keyring.Set(service, user, got)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("No value entered")
		}

		trelloApiKey, err = keyring.Get(service, user)
		if err != nil {
			log.Fatal(err)
		}
	}

	client := trello.NewClient(cfg.AppKey, trelloApiKey)

	got, ok := inputbox.InputBox("Card Text", "Input card text", "")
	if ok {
		card := trello.Card{
			Name:   got,
			IDList: cfg.IdList,
		}
		err := client.CreateCard(&card, trello.Defaults())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No value entered")
	}
}
