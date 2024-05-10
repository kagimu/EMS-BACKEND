package controllers

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateNewInvoiceNumber() string {
	now := time.Now()
	year := now.Format("06")  // YY format
	month := now.Format("01") // MM format

	// Generate random four-digit number
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(10000) // Generate random number between 0 and 9999

	// Format invoice number
	return fmt.Sprintf("%s%s%04d", year, month, randomNum)
}
