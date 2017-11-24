package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

func Appender() func(string) (int, string) {
	nonce := -1
	return func(x string) (int, string) {
		nonce += 1
		return nonce, x + strconv.Itoa(nonce)
	}
}

func main() {
	AppendNonce := Appender()
	h := sha256.New()
	target := "00000"

	for {
		nonce, AppendedNonce := AppendNonce("hello")
		h.Write([]byte(AppendedNonce))

		calculated_hash := h.Sum(nil)

		fmt.Printf("%d %x\n", nonce, calculated_hash)
		calculated_hash_s := fmt.Sprintf("%x", calculated_hash)

		if calculated_hash_s[:len(target)] == target {
			fmt.Println(calculated_hash[:len(target)/2])
			fmt.Println(target)
			break
		}
	}

}
