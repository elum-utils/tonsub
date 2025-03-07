package main

import (
	"encoding/json"
	"fmt"

	"github.com/elum-utils/tonsub"
)

func main() {

	subs, err := tonsub.New(
		"UQAno5PEMnsMt26bPgnXeFMOBVzSNHor2ctgyALrg3oD5aF6",
		"https://ton.org/global.config.json",
	)

	if err != nil {
		panic(err.Error())
	}

	subs.OnJetton(func(t *tonsub.RootJetton) {
		jsonData, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		fmt.Println("Transaction JETTON")
		fmt.Println(string(jsonData))
	})

	subs.OnTON(func(t *tonsub.RootTON) {
		jsonData, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		fmt.Println("Transaction TON")
		fmt.Println(string(jsonData))
	})

	subs.OnNFT(func(t *tonsub.RootNFT) {
		jsonData, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		fmt.Println("Transaction NFT")
		fmt.Println(string(jsonData))
	})

	select {}

}
