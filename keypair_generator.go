package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"log"
	"os"
)

func main() {
	// argv
	var privateKeyFileName = flag.String("private", "private.key", "the name of private key filename")
	var publicKeyFileName = flag.String("public", "public.key", "the name of public key filename")
	flag.Parse()

	// key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := &privateKey.PublicKey

	// marshal
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)

	// save
	err = os.WriteFile(*privateKeyFileName, privateKeyBytes, 0600)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(*publicKeyFileName, publicKeyBytes, 0600)
	if err != nil {
		log.Fatal(err)
	}
}
