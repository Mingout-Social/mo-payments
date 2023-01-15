package payments

import (
	"errors"
	"github.com/Mingout-Social/mo-payments/providers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ProviderRazorpay = "RAZORPAY"
const ProviderCashfree = "CASHFREE"

type PaymentStatus string

const (
	Success PaymentStatus = "SUCCESS"
	Failed  PaymentStatus = "FAILED"
	Pending PaymentStatus = "PENDING"
)

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

type PaymentDetail struct {
	ID        primitive.ObjectID `bson:"id" json:"id"`
	OrderID   string             `bson:"order_id" json:"order_id"`
	PaymentID string             `bson:"payment_id" json:"payment_id"`
	Amount    int64              `bson:"amount" json:"amount"`
	Status    PaymentStatus      `bson:"status" json:"status"`
	Provider  string             `bson:"provider" json:"provider"`
}

func GenerateOrder(orderAmount int64, userId primitive.ObjectID, mobileNumber string, email string, provider string) (PaymentDetail, error) {
	var order OrderResponse
	var err error
	var paymentDetail PaymentDetail

	if provider == "" {
		err = errors.New("No Payment Provider Configured!")
	}

	if provider == ProviderRazorpay {
		order, err = providers.CreateRazorpayOrder(orderAmount)
		paymentDetail = PaymentDetail{
			ID:       primitive.NewObjectID(),
			OrderID:  order.ID,
			Amount:   order.Amount,
			Status:   Pending,
			Provider: provider,
		}
	} else if provider == ProviderCashfree {
		orderId := primitive.NewObjectID()
		order, err = providers.CreateCashFreeOrder(orderAmount, orderId, userId, mobileNumber, email)
		paymentDetail = PaymentDetail{
			ID:       orderId,
			OrderID:  order.ID,
			Amount:   order.Amount,
			Status:   Pending,
			Provider: provider,
		}
	}

	return paymentDetail, err
}

func VerifyPayment(orderId string, provider string) (VerifyPaymentResponse, error) {
	var verifyPaymentResponse VerifyPaymentResponse

	if provider != ProviderCashfree {
		return verifyPaymentResponse, nil
	}

	verifyPaymentResponse, err := providers.VerifyPaymentOrder(orderId)

	return verifyPaymentResponse, err
}
