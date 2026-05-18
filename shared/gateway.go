package shared

import "context"

type Gateway interface {
	RequestToPay(ctx context.Context, req RequestToPayRequest) (*RequestToPayResponse, error)
	CheckPaymentStatus(ctx context.Context, referenceID string) (*PaymentStatusResponse, error)
}
