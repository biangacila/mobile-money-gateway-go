package main

import (
	"fmt"

	"github.com/biangacila/mobile-money-gateway-go/momo"
	"github.com/biangacila/mobile-money-gateway-go/shared"
	"github.com/joho/godotenv"
)

func main() {

	if shared.Getenv("MOMO_COLLECTION_PRIMARY_KEY", "") == "" {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("Error loading .env file: " + err.Error())
		}
	}

	// Get environment variables
	subscriptionKey := shared.Getenv("MOMO_COLLECTION_PRIMARY_KEY", "")
	apiUser := shared.Getenv("MOMO_API_USER", "YOUR_API_USER_UUID")
	apiKey := shared.Getenv("MOMO_API_KEY", "YOUR_API_KEY")
	baseURL := shared.Getenv("MOMO_BASE_URL", "https://sandbox.momodeveloper.mtn.com")
	targetEnvironment := shared.Getenv("MOMO_TARGET_ENVIRONMENT", "sandbox")

	client := momo.NewClient(subscriptionKey, apiUser, apiKey, baseURL, targetEnvironment)

	// Generate a collection token
	token, err := client.GenerateCollectionToken()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Collection token generated: %s\n", token)

	payload := momo.RequestToPayPayload{
		Amount:     "10",
		Currency:   "EUR",
		ExternalID: "INV-1003",
		Payer: momo.Payer{
			PartyIDType: "MSISDN",
			PartyID:     "46733123453",
		},
		PayerMessage: "Payment for order biacibenga Solution",
		PayeeNote:    "Thank you",
	}
	referenceID, err := client.RequestToPay(token, payload)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Payment request sent: %s\n", referenceID)

	// Check payment status
	status, err := client.GetRequestToPayStatus(token, referenceID)
	if err != nil {
		panic(err)

	}
	fmt.Printf("Payment status: %s\n", status)

}
