package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func computeMd5Base64(filepath string) (b64digest string, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	b64digest = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return b64digest, err
}

func computeSha256(filepath string) (hexdigest string, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	hexdigest = hex.EncodeToString(h.Sum(nil))
	return hexdigest, err
}
