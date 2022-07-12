package cron

import (
	"log"
	"strconv"

	"egox/database/model"
	"egox/pkg/get_nft_owner_of"
	"egox/pkg/get_nft_token_uri"
	"egox/pkg/get_nft_total_supply"

	"github.com/robfig/cron/v3"
)

func Cronjob() {
	log.Println("Cron Starting...")

	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	c.AddFunc("* * * * *", func() {
		contracts := model.GetAllContracts()
		for _, contract := range contracts {
			getItemsOwnerOfByContract(contract)
		}

	})
	c.Start()
	defer c.Stop()
	select {}
}

func getItemsOwnerOfByContract(contract string) {
	totalSupply, err := get_nft_total_supply.GetTotalSupply(contract)
	if err != nil {
		log.Printf("Get NFT total supply failed, err: %v\n", err)
	}

	log.Printf("Contract: %s total supply: %s", contract, totalSupply)

	totalAmount, _ := strconv.Atoi(totalSupply.String())

	for i := 1; i <= totalAmount; i++ {
		owner, err := get_nft_owner_of.GetOwnerOf(contract, i)
		if err != nil {
			log.Printf("Get NFT owner of failed, err: %v\n", err)
		}

		tokenURI, err := get_nft_token_uri.GetNftTokenUri(contract, i)

		if err != nil {
			log.Printf("Get NFT tokenURI failed, err: %v\n", err)
		}

		log.Printf("Token Id: %v owner: %s tokenURI: %s", i, owner, *tokenURI)

	}
}
