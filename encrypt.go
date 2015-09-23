/*
   encrypt.go encrypts a message using secret key NaCl
   and a user provided key (default key: qwerty).
   Then it writes the result to encrypted.dat
*/
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"io/ioutil"
	"log"
)

// These are defined in golang.org/x/crypto/nacl/secretbox
const keySize = 32
const nonceSize = 24

// If a key is not provided, “qwerty” will be used
var userKey = flag.String("k", "qwerty", "encryption key")

// NaCl's key must be 32 bytes. If the provided key is less than that,
// we will pad it with the appropriate number of bytes from pad.
// pad should be the same for encrypter and decrypter
var pad = []byte("«super jumpy fox jumps all over»")

// Message to encrypt (plaintext)
var message = []byte("Hello world")

func main() {
	flag.Parse()

	// key is a temporary holder for the real key (naclKey)
	key := []byte(*userKey)
	// NaCl's key has a constant size of 32 bytes.
	// The user provided key probably is less than that. We pad it with
	// a long enough string and truncate anything we don't need later on.
	key = append(key, pad...)

	// NaCl's key should be of type [32]byte.
	// Here we create it and truncate key bytes beyond 32
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])

	// Nonce is a [24]byte variable that should be unique amongst messages.
	// Nonce is generated by the encrypter.
	// 24 bytes is large enough to avoid collisions while reading from rand.
	// If it suits you better, you may use a non random nonce (e.g a counter).
	// It only has to change per message.
	nonce := new([nonceSize]byte)
	// Read bytes from random and put them in nonce until it is full.
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		log.Fatalln("Could not read from random:", err)
	}

	// out will hold the nonce and the encrypted message (ciphertext)
	out := make([]byte, nonceSize)
	// Copy the nonce to the start of out
	copy(out, nonce[:])
	// Encrypt the message and append it to out, assign the result to out
	out = secretbox.Seal(out, message, nonce, naclKey)

	// Write the result to a file
	err = ioutil.WriteFile("encrypted.dat", out, 0644)
	if err != nil {
		log.Fatalln("Error while writing encrypted.dat: ", err)
	}

	fmt.Printf("Message encrypted succesfully. Total size is %d bytes,"+
		" of which %d bytes is the message, "+
		"%d bytes is the nonce and %d bytes is the overhead.\n",
		len(out), len(message), nonceSize, secretbox.Overhead)
	fmt.Printf("The encryption key is: '%s'\n", naclKey[:])
	// Nonce may contain non-printable characters. We print it as []byte
	fmt.Printf("The nonce is: '%v'\n", nonce[:])
}
