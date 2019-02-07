package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/etherparty/bitthief/btckey"
)

// TODO: Pass list of words as CLI arg and not hardcoded

func main() {
	wordsPtr := flag.String("file", "corncob_lowercase.txt", `Filename you want to scan through. Must be delineated by lines. EG:
		word
		word2
		word3
		word4`)

	wordPtr := flag.String("word", "", `Specific word you want to search for`)
	timestampPtr := flag.Bool("timestamps", true, "Defaults to true. Set to true if you want it to loop through all timestamps as seed values since Jan 1 2009. False if you don't.")

	flag.Parse()

	// if a single word was passed in, let it chill.
	if *wordPtr != "" {
		handle([]byte(*wordPtr), true)

		os.Exit(0)
	}

	file, err := os.Open(*wordsPtr)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	// loop through timestamps
	if *timestampPtr {
		janFirst2018 := 1167609600
		now := 1519173144

		go func() {
			for i := janFirst2018; i < now; i++ {
				seed := strconv.Itoa(i)

				handle([]byte(seed), false)
			}
		}()
	}
	// loop through words

	// loop through all words
	go func() {
		for _, word := range words {
			handle([]byte(word), false)
		}
	}()

	forever := make(chan bool)
	<-forever

}

func handle(seed []byte, singleton bool) {
	hashedBytes := sha256.Sum256(seed)

	privateKey := btckey.PrivateKey{}
	if err := privateKey.FromBytes(hashedBytes[:]); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	btcAddress := privateKey.PublicKey.ToAddressUncompressed()

	// convert private key to address
	url := fmt.Sprintf("https://blockchain.info/q/addressbalance/%s", btcAddress)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	contentValue, err := strconv.Atoi(string(contents))
	returnString := fmt.Sprintf("Address %s has balance %s, Password: %s, private key: %s \n", btcAddress, string(contents), string(seed), privateKey.ToWIF())
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	if contentValue > 0 {
		fmt.Println("========================== \n THE IMPOSSIBLE IS POSSIBLE! \n =========================")
		filename := fmt.Sprintf("%s.txt", privateKey.ToWIF())
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		file.WriteString(returnString)
		file.Close()
	} else {
		file, err := os.Open("attempts.txt")
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		file.WriteString(returnString)
		file.Close()
	}

	// if its a single one, lets also check tx count.
	if singleton {
		url := fmt.Sprintf("https://blockchain.info/q/getreceivedbyaddress/%s", btcAddress)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("Total Satoshis Received: %s\n", contents)
	}

	fmt.Println(returnString)
	fmt.Printf("Transaction URL: https://blockexplorer.com/address/%s \n", btcAddress)
}
