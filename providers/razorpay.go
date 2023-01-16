package providers

import (
	"os"

	"github.com/Mingout-Social/mo-payments/responses"
	"github.com/mitchellh/mapstructure"
	"github.com/razorpay/razorpay-go"
	"github.com/sirupsen/logrus"
)

const (
	OrderStatusCreated   = "created"
	OrderStatusAttempted = "attempted"
	OrderStatusPaid      = "paid"

	PaymentStatusCaptured = "captured"
)

func getClient() *razorpay.Client {
	key := os.Getenv("RAZORPAY_KEY")
	secret := os.Getenv("RAZORPAY_SECRET")

	return razorpay.NewClient(key, secret)
}

func CreateRazorpayOrder(amount int64) (responses.OrderResponse, error) {
	var order responses.OrderResponse

	payload := map[string]interface{}{
		"amount":          amount,
		"currency":        "INR",
		"partial_payment": false,
		"notes":           map[string]interface{}{},
	}

	client := getClient()
	response, err := client.Order.Create(payload, nil)

	if err != nil {
		logrus.Errorf("CreateOrder: %v", err)
		return order, err
	}

	err = mapstructure.Decode(response, &order)
	if err != nil {
		logrus.Errorf("CreateOrder: %v", err)
	}

	return order, err
}
