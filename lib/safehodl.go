package safehodl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/alexedwards/argon2id"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/jedib0t/go-pretty/v6/table"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

func getSafeHodlDbPath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/.safehodl", userHomeDir)
}

func openSafeHodlDb() *bolt.DB {
	db, err := bolt.Open(getSafeHodlDbPath(), 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// StartInteractiveSafeHodlConfiguration SafeHODL interactive configuration.
func StartInteractiveSafeHodlConfiguration() {
	fmt.Printf("========================\n SafeHODL Configuration\n========================\n")
	fmt.Println("Enter your Bitcoin holdings amount: ")
	btcAmount, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	btcAmountFloat, err := strconv.ParseFloat(string(btcAmount), 64)
	if err != nil {
		fmt.Printf("Invalid BTC amount. %s", err)
		os.Exit(1)
	}

	trySetupPassphrase()
	persistBtcAmountHoldings(btcAmountFloat)

	fmt.Println("All set! now you can run \"safehodl\".")
}

// trySetupPassphrase checks if a passphrase hash is already in the db.
func trySetupPassphrase() {
	db := openSafeHodlDb()
	defer db.Close()

	var existingPassphraseHashBlob []byte
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("auth"))
		if err != nil {
			return err
		}
		existingPassphraseHashBlob = b.Get([]byte("passphrase"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(existingPassphraseHashBlob) == 0 {
		fmt.Println("Enter a secret passphrase for secure access: ")
		plainNewPassphraseBlob, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}

		hashedNewPassphrase, err := argon2id.CreateHash(string(plainNewPassphraseBlob), argon2id.DefaultParams)
		if err != nil {
			log.Fatal(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("auth"))
			return b.Put([]byte("passphrase"), []byte(hashedNewPassphrase))
		})
		if err != nil {
			log.Fatal(err)
		}
	}

}

// AssertPassphrase asks the user for a passphrase and if a wrong passphrase is provided then optimistically removes the data.
func AssertPassphrase() {
	fmt.Println("Enter passphrase: ")
	plainPassphrase, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	db := openSafeHodlDb()
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("auth"))
		passphraseHash := b.Get([]byte("passphrase"))

		matchPassphrase, err := argon2id.ComparePasswordAndHash(string(plainPassphrase), string(passphraseHash))
		if err != nil {
			return err
		}

		if !matchPassphrase {
			os.Remove(getSafeHodlDbPath()) // nuke data
			fmt.Print("Incorrect passphrase!!")
			os.Exit(0)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

}

func getStoredPassphrase() (encKey []byte, salt []byte, params *argon2id.Params) {
	db := openSafeHodlDb()
	defer db.Close()

	var passphraseHashBlob []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("auth"))
		passphraseHashBlob = b.Get([]byte("passphrase"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	hashParams, passphraseSalt, encPassphrase, err := argon2id.DecodeHash(string(passphraseHashBlob))
	if err != nil {
		log.Fatal(err)
	}

	return encPassphrase, passphraseSalt, hashParams
}

// deriveIDKeyFromStoredPassphraseHash derives a AES-256 private key from stored passphrase.
func deriveIDKeyFromStoredPassphraseHash() []byte {
	encKey, salt, _ := getStoredPassphrase()

	return argon2.IDKey(encKey, salt, 1, 64*1024, 2, 32)
}

func encHoldingAmount(amount float64) []byte {
	key := deriveIDKeyFromStoredPassphraseHash()

	c, err := aes.NewCipher(key)
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
func decBtcHoldingAmountFromDb() float64 {
	var hodlingsAmountBlob []byte

	db := openSafeHodlDb()

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("holdings"))
		hodlingsAmountBlob = b.Get([]byte("BTC"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	key := deriveIDKeyFromStoredPassphraseHash()

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

	nonceSize := gcm.NonceSize()
	if len(hodlingsAmountBlob) < nonceSize {
		log.Fatal(err)
	}

	nonce, ciphertext := hodlingsAmountBlob[:nonceSize], hodlingsAmountBlob[nonceSize:]
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

func persistBtcAmountHoldings(btcAmount float64) {
	encBtcAmount := encHoldingAmount(btcAmount)

	db := openSafeHodlDb()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("holdings"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("BTC"), encBtcAmount)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

func hasBtcAmountSet() bool {
	var hasBtcHoldingsAmountSet bool

	db := openSafeHodlDb()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("holdings"))
		if err != nil {
			return err
		}
		btcHoldingsAmountBlob := b.Get([]byte("BTC"))
		hasBtcHoldingsAmountSet = len(btcHoldingsAmountBlob) > 0
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return hasBtcHoldingsAmountSet
}

// GetBtcAmount decrypts holdings amount from encrypted .safehodl file and returns it.
func GetBtcAmount() (hasBtcHoldingsAmountSet bool, storedBtcAmount float64) {
	if !hasBtcAmountSet() {
		return false, 0
	}

	return true, decBtcHoldingAmountFromDb()
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
	_, btcAmount := GetBtcAmount()

	btcUSD, btcEUR := getBitcoinPrices()

	const oneBitcoinSatoshi float64 = 0.00000001
	satoshiUSD := btcUSD * oneBitcoinSatoshi
	satoshiEUR := btcEUR * oneBitcoinSatoshi

	btcAmountUSD := btcAmount * btcUSD
	btcAmountEUR := btcAmount * btcEUR

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"BTC", "USD", "EUR"})
	t.AppendRows([]table.Row{
		{"1 Bitcoin", fmt.Sprintf("$%.2f", btcUSD), fmt.Sprintf("%.2f€", btcEUR)},
		{"1 Satoshi", fmt.Sprintf("$%.4f", satoshiUSD), fmt.Sprintf("%.4f€", satoshiEUR)},
	})
	t.AppendSeparator()
	t.AppendRow([]interface{}{fmt.Sprintf("%.8f", btcAmount), fmt.Sprintf("$%.2f", btcAmountUSD), fmt.Sprintf("%.2f€", btcAmountEUR)})

	t.Render()
}
