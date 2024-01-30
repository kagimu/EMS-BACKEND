package controllers

import (
	"math/rand"
	"strconv"
	"time"
)

// TODO: Decide on payment reference structure
func GeneratePaymentReference() (string, error) {
	rand.NewSource(time.Now().UnixNano())

	randomNumber := rand.Intn(20000000000 + 1)
	randomNumberString := strconv.Itoa(randomNumber)

	return randomNumberString, nil
}
