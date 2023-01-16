package responses

type OrderResponse struct {
	ID         string
	Entity     string
	Amount     int64
	AmountPaid int
	AmountDue  int64
	Currency   string
	Status     string
}

type VerifyPaymentResponse struct {
	OrderId   string
	PaymentId string
	Status    string
}
