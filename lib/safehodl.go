package safehodl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
)

var Secret32BytesKeyAES string

func getSafeHodlDotFilePath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/.safehodl", userHomeDir)
}

// StartInteractiveSafeHodlConfiguration SafeHODL interactive configuration.
func StartInteractiveSafeHodlConfiguration() {
	fmt.Printf("========================\n SafeHODL Configuration\n========================\n")
	fmt.Println("Enter your Bitcoin holdings amount: ")
	amount, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	amountFloat, err := strconv.ParseFloat(string(amount), 64)
	if err != nil {
		fmt.Printf("Invalid amount. %s", err)
		os.Exit(1)
	}
	persistHodlAmount(amountFloat)

	fmt.Println("All set! now you can run \"safehodl\".")
}

// AssertPinCodeForUsage asks for access pin code and if a wrong pin code is entered then removes the data.
func AssertPinCodeForUsage(pinCode string) {
	fmt.Println("Enter PIN code: ")
	passwd, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	if string(passwd) != pinCode {
		os.Remove(getSafeHodlDotFilePath()) // nuke data
		fmt.Print("Incorrect PIN code!!")
		os.Exit(0)
	}
}

func encHodlAmount(amount float64) []byte {
	c, err := aes.NewCipher([]byte(Secret32BytesKeyAES))
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	return gcm.Seal(nonce, nonce, []byte(fmt.Sprintf("%f", amount)), nil)
}
func persistHodlAmount(amount float64) {
	encAmount := encHodlAmount(amount)

	err := ioutil.WriteFile(getSafeHodlDotFilePath(), encAmount, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func decHodlAmountFromFile() float64 {
	hodlAmountEncBlob, err := ioutil.ReadFile(getSafeHodlDotFilePath())
	if err != nil {
		log.Fatal(err)
	}

	c, err := aes.NewCipher([]byte(Secret32BytesKeyAES))
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

	nonceSize := gcm.NonceSize()
	if len(hodlAmountEncBlob) < nonceSize {
		log.Fatal(err)
	}

	nonce, ciphertext := hodlAmountEncBlob[:nonceSize], hodlAmountEncBlob[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}

	hodlAmount, err := strconv.ParseFloat(string(plaintext[:]), 64)
	if err != nil {
		log.Fatal(err)
	}

	return hodlAmount
}

func hasHodlAmount() bool {
	_, err := os.Stat(getSafeHodlDotFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal(err)
	}

	return true
}

// GetHodlAmount decrypts holdings amount from encrypted .safehodl file and returns it.
func GetHodlAmount() (hasStoredAmount bool, storedAmount float64) {
	if !hasHodlAmount() {
		return false, 0
	}

	return true, decHodlAmountFromFile()
}

func getBitcoinPrices() (usd float64, eur float64) {
	api := krakenapi.New("", "")
	ticker, err := api.Ticker(krakenapi.XXBTZUSD, krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	btcUSD, err := strconv.ParseFloat(ticker.XXBTZUSD.Close[0], 64)
	if err != nil {
		log.Fatal(err)
	}

	btcEUR, err := strconv.ParseFloat(ticker.XXBTZEUR.Close[0], 64)
	if err != nil {
		log.Fatal(err)
	}

	return btcUSD, btcEUR
}

// DisplayHodlInfo prints on screen holdings information.
func DisplayHodlInfo() {
	_, hodlAmount := GetHodlAmount()

	btcUSD, btcEUR := getBitcoinPrices()

	const oneBitcoinSatoshi float64 = 0.00000001
	satoshiUSD := btcUSD * oneBitcoinSatoshi
	satoshiEUR := btcEUR * oneBitcoinSatoshi

	hodlAmountUSD := hodlAmount * btcUSD
	hodlAmountEUR := hodlAmount * btcEUR

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"BTC", "USD", "EUR"})
	t.AppendRows([]table.Row{
		{"1 Bitcoin", fmt.Sprintf("$%.2f", btcUSD), fmt.Sprintf("%.2f€", btcEUR)},
		{"1 Satoshi", fmt.Sprintf("$%.4f", satoshiUSD), fmt.Sprintf("%.4f€", satoshiEUR)},
	})
	t.AppendSeparator()
	t.AppendRow([]interface{}{fmt.Sprintf("%.8f", hodlAmount), fmt.Sprintf("$%.2f", hodlAmountUSD), fmt.Sprintf("%.2f€", hodlAmountEUR)})

	t.Render()
}
