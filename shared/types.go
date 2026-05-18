package shared

type Payer struct {
	PartyIDType string `json:"partyIdType"`
	PartyID     string `json:"partyId"`
}

type RequestToPayRequest struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalID   string `json:"externalId"`
	Payer        Payer  `json:"payer"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`
}

type RequestToPayResponse struct {
	ReferenceID string `json:"referenceId"`
	Status      string `json:"status"`
}

type PaymentStatusResponse struct {
	ReferenceID            string `json:"referenceId"`
	Amount                 string `json:"amount"`
	Currency               string `json:"currency"`
	FinancialTransactionID string `json:"financialTransactionId"`
	ExternalID             string `json:"externalId"`
	Payer                  Payer  `json:"payer"`
	PayerMessage           string `json:"payerMessage"`
	PayeeNote              string `json:"payeeNote"`
	Status                 string `json:"status"`
	Reason                 string `json:"reason"`
}
