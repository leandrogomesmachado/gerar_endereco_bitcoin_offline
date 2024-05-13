package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	// Open the mnemonic.txt file
	file, err := os.Open("mnemonic.txt")
	if err != nil {
		fmt.Println("Error opening mnemonic.txt:", err)
		return
	}
	defer file.Close()

	// Read the words from the file
	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading mnemonic.txt:", err)
		return
	}

	// Select 12 random words from the slice
	var selectedWords []string
	for i := 0; i < 12; i++ {
		index := rand.Intn(len(words))
		selectedWords = append(selectedWords, words[index])
	}

	// Generate a BIP39 mnemonic from the selected words
	mnemonic := strings.Join(selectedWords, " ")

	// Convert the mnemonic to seed
	seed := bip39.NewSeed(mnemonic, "")

	// Generate a BIP32 HD wallet
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		fmt.Println("Error generating master key:", err)
		return
	}

	// Get the raw private key bytes
	privateKeyBytes := masterKey.Key

	// Encode the raw private key bytes to hexadecimal string
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// Derive a BIP32 Extended Public Key
	publicKey := masterKey.PublicKey()

	// Derive the Bitcoin address from the public key
	pubKeyHash := btcutil.Hash160(publicKey.Key)
	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Error deriving address:", err)
		return
	}

	// Write the mnemonic, raw private key, and address to a file
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	writer.WriteString("Mnemonic: " + mnemonic + "\n")
	writer.WriteString("Raw Private Key: " + privateKeyHex + "\n")
	writer.WriteString("Address: " + address.EncodeAddress() + "\n")

	fmt.Println("Generated Bitcoin Address:", address.EncodeAddress())
	fmt.Println("Raw Private Key:", privateKeyHex)
}
