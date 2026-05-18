package momo

type RequestToPayPayload struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalID   string `json:"externalId"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`

	Payer struct {
		PartyIDType string `json:"partyIdType"`
		PartyID     string `json:"partyId"`
	} `json:"payer"`
}
type Payer struct {
	PartyIDType string `json:"partyIdType"`
	PartyID     string `json:"partyId"`
}
