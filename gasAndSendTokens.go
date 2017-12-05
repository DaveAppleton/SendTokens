package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"./admin"

	"github.com/DaveAppleton/etherUtils"
	"github.com/DaveAppleton/ether_go/ethKeys"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"context"

	"./parityclient"
)

var client *parityclient.Client
var fromCSV [][]string

var contracts [][]string
var addresses [][]string
var skipped [][]string
var theToken *admin.ERC20
var tokenAddress common.Address
var mainAddress common.Address
var defaultGas *big.Int

var skipRows int
var addressCol int
var amountCol int
var gasPrice *big.Int
var decimalPlaces int64

func getClient() (client *parityclient.Client, err error) {
	endPoint := "http://localhost:8545"
	if len(endPoint) == 0 {
		endPoint = "/Users/daveappleton/Library/Ethereum/geth.ipc"
	}
	//deadline := time.Now().Add(20 * time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), deadline)
	client, err = parityclient.Dial(endPoint)
	return
}

var input string

func readCSVFile() (data [][]string, err error) {
	sf, err := os.Open(input)
	if err != nil {
		fmt.Println(err, "[", input, "]")
		log.Fatal(err, input)

	}
	data, err = csv.NewReader(sf).ReadAll()
	return
}

func isItAContract(add common.Address) (isContract bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	code, err := client.CodeAt(ctx, add, nil)
	if err != nil {
		return
	}
	isContract = len(code) != 0
	return
}

func sortOutGas(data [][]string) (sumGas *big.Int, sumVal *big.Int, err error) {
	var msg ethereum.CallMsg
	//var dataString []byte
	sumGas = big.NewInt(0)
	sumVal = big.NewInt(0)
	//ctx := context.Background()

	msg.From = mainAddress
	msg.To = &tokenAddress
	msg.Value = big.NewInt(0)

	for _, rec := range data {

		add := common.HexToAddress(rec[addressCol])
		val, ok := etherUtils.StrToDecimals(rec[amountCol], decimalPlaces)
		if !ok {
			log.Fatal("Error should not happen here")
			continue
		}
		// = still working on estimation
		// dataString =
		// msg.Data = dataString
		// est, err := client.EstimateGas(ctx, msg)
		// if err != nil {
		// 	fmt.Println("Estimate : ", err)
		est := defaultGas
		// }
		//

		fmt.Println(add.Hex(), val, "gas : ", est)
		sumGas = new(big.Int).Add(sumGas, est)
		sumVal = new(big.Int).Add(sumVal, val)
	}
	return
}

func payAll(payer *ethKeys.AccountKey, list [][]string) {

	for _, rec := range list {
		add := common.HexToAddress(rec[addressCol])
		val, _ := etherUtils.StrToDecimals(rec[amountCol], decimalPlaces)

		tx := bind.NewKeyedTransactor(payer.GetKey())
		tx.GasPrice = gasPrice
		txn, err := theToken.Transfer(tx, add, val)
		if err != nil {
			fmt.Println("Error sending ", rec[amountCol], " to ", rec[addressCol], err)
		} else {
			fmt.Println(txn.Hash().Hex(), rec[amountCol], rec[addressCol])
		}
	}
}

func main() {
	var err error
	var pay string

	flag.StringVar(&input, "input", "", " path to input file (CSV)")
	flag.StringVar(&pay, "pay", "no", " pay to users?")
	flag.Parse()
	if len(input) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err = viper.ReadInConfig()    // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	tokenAddressString := viper.GetString("TOKEN_ADDRESS")
	tokenAddress = common.HexToAddress(tokenAddressString)

	mainAddressString := viper.GetString("MAIN_ADDRESS")
	mainAddress = common.HexToAddress(mainAddressString)

	gasPriceString := viper.GetString("GAS_PRICE")
	gasPrice, ok := new(big.Int).SetString(gasPriceString, 10)
	if !ok {
		log.Fatalf("Bad gas price in config.json : %s is not a number\n", gasPriceString)
	}

	gasLimitString := viper.GetString("DEFAULT_GAS")
	defaultGas, ok = new(big.Int).SetString(gasLimitString, 10)
	if !ok {
		log.Fatalf("Bad gas limit in config.json : %s is not a number\n", gasLimitString)
	}

	skipRows = viper.GetInt("SKIP_ROWS")
	addressCol = viper.GetInt("ADDRESS_COL")
	amountCol = viper.GetInt("AMOUNT_COL")

	client, err = getClient()
	if err != nil {
		log.Fatal(err)
	}
	theToken, err = admin.NewERC20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	theTokenName, err := theToken.Name(nil)
	if err != nil {
		log.Fatal("Error loading token name : ", err)
	}
	fmt.Println("Token : ", theTokenName)

	useDecimals := viper.GetBool("AMOUNT_IN_DECIMALS")
	if useDecimals {
		decimalPlaces8, err := theToken.Decimals(nil)
		if err != nil {
			log.Fatal("Error loading token decimals : ", err)
		}
		decimalPlaces = int64(decimalPlaces8)
	} else {
		decimalPlaces = 0
	}

	fromCSV, err := readCSVFile()
	if err != nil {
		log.Fatal(err)
	}

	payer := ethKeys.NewKey("keys/payment")

	err = payer.RestoreOrCreate()
	if err != nil {
		log.Fatal(err)
	}

	tokenBal, err := theToken.BalanceOf(nil, payer.PublicKey())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Payment Address ", payer.PublicKeyAsHexString())
	fmt.Println()
	fmt.Println()
	fmt.Println()
	lineNo := 0
	for _, rec := range fromCSV {
		if lineNo < skipRows {
			lineNo++
			continue
		}
		rec[amountCol] = strings.TrimSpace(rec[amountCol])
		rec[addressCol] = strings.TrimSpace(rec[addressCol])
		add := common.HexToAddress(rec[addressCol])
		_, ok := etherUtils.StrToDecimals(rec[amountCol], decimalPlaces)
		if !ok {
			skipped = append(skipped, []string{rec[addressCol], "Number error : " + rec[amountCol]})
			continue
		}
		isContract, err := isItAContract(add)
		if err != nil {
			fmt.Println(err)
			skipped = append(skipped, []string{rec[addressCol], err.Error()})
			continue
		}
		if isContract {
			contracts = append(contracts, rec)
			continue
		}
		addresses = append(addresses, rec)
	}
	fmt.Println("Addresses----------------------------------->", len(addresses))
	addGas, addVal, addErr := sortOutGas(addresses)
	fmt.Println("Contracts----------------------------------->", len(contracts))
	conGas, conVal, conErr := sortOutGas(contracts)
	fmt.Println("SKIPPED ------------------------------------>", len(skipped))
	for _, rec := range skipped {
		fmt.Println(rec[0], rec[1])
	}

	AmountToPay := new(big.Int).Add(addVal, conVal)
	TotalGasLimit := new(big.Int).Add(addGas, conGas)
	TotalGasUse := new(big.Int).Mul(TotalGasLimit, gasPrice)
	fmt.Println()

	fmt.Printf("Tokens Required %s, total gas %d, total gas price %s\n",
		etherUtils.CoinToStr(AmountToPay, int(decimalPlaces)),
		TotalGasLimit,
		etherUtils.EtherToStr(TotalGasUse))

	if strings.Compare(pay, "yes") != 0 {
		fmt.Println()
		fmt.Println("Payment not enabled")
		os.Exit(0)
	}
	if len(skipped) > 0 {
		fmt.Println("Payment not allowed when we have skipped entries")
		os.Exit(1)
	}
	if addErr != nil || conErr != nil {
		fmt.Println("Payment not allowed when we have errors")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	balance, err := client.BalanceAt(ctx, payer.PublicKey(), nil)
	if err != nil {
		log.Fatal(err)
	}
	if balance.Cmp(TotalGasUse) < 0 {
		fmt.Println("insufficient ether - need ", etherUtils.EtherToStr(TotalGasUse), " found - ", etherUtils.EtherToStr(balance))
		return
	}
	if tokenBal.Cmp(AmountToPay) < 0 {
		fmt.Println("insufficient Tokens - need ", etherUtils.CoinToStr(AmountToPay, 8), " found - ", etherUtils.CoinToStr(tokenBal, 8))
		return

	}
	payAll(payer, addresses)
	payAll(payer, contracts)

}
