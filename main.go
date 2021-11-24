// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

// command line interface to interact with blockchain and program
type CommandLine struct{}

// printUsage will display what options are available to the user
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("getbalance - address ADDRESS - get balance for ADDRESS")
	fmt.Println("createblockchain - address ADDRESS creates a blockchain and rewards the mining fee")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount of coins from one address to another")
	fmt.Println("createwallet - Creates a new wallet")
	fmt.Println("listaddresses - Lists the addresses in the wallet file")
}

// lists addresses
func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

// create wallet
func (cli *CommandLine) createWallet() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

// create blockchain
func (cli *CommandLine) createBlockChain(address string) {
	newChain := blockchain.InitBlockChain(address)
	newChain.Database.Close()
	fmt.Println("Finished creating chain")
}

// get balance
func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain > FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// send money
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)

	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

// validateArgs ensures the CLI was given valid input
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		// Goexit will exit the application by shutting down the goroutine
		// if you were to use os.exit you might corrupt the data
		runtime.Goexit()
	}
}

// printChain will display the entire contents of the blockchain
func (cli *CommandLine) printChain() {
	iterator := cli.blockchain.Iterator()

	for {
		block := iterator.Next()
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// this works because the genesis block has not PrevHash
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// run will start the command line
func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagset("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	// Parsed() will return true if the object has been called
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmoutn <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if listAddressesCmd > Parsed() {
		cli.listAddresses()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
}

// run the code
func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.run()
}
