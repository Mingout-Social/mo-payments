package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Mingout-Social/mo-payments/responses"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Success = "SUCCESS"
	Failed  = "FAILED"
	Pending = "PENDING"
)

func CreateCashFreeOrder(amount int64, orderId primitive.ObjectID, userId primitive.ObjectID, mobileNo string, emailId string) (responses.OrderResponse, error) {
	var order responses.OrderResponse
	amount = amount / 100

	body, err := json.Marshal(map[string]interface{}{
		"order_id":       orderId.Hex(),
		"order_amount":   amount,
		"order_currency": "INR",
		"customer_details": map[string]interface{}{
			"customer_id":    userId.Hex(),
			"customer_email": emailId,
			"customer_phone": mobileNo,
		},
	})

	if err != nil {
		logrus.Errorf("Error While Marshaling createCashFreeOrder params: %v", err)
		return order, err
	}

	url := os.Getenv("CASHFREE_BASE_URI") + "/orders"

	payload := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPost, url, payload)

	if err != nil {
		logrus.Errorf("Error while create new request to createCashFreeOrder: %v", err)
		return order, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-client-id", os.Getenv("CASHFREE_APP_ID"))
	req.Header.Set("x-client-secret", os.Getenv("CASHFREE_SECRET_KEY"))
	req.Header.Set("x-api-version", os.Getenv("CASHFREE_API_VERSION"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error while making cashfree order: %v", err, resp, req)
		return order, err
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error while making cashfree order: %v", err, resp, req)
		return order, err
	}
	var cashFreeResponse map[string]interface{}

	err = json.Unmarshal(response, &cashFreeResponse)
	if err != nil {
		logrus.Errorf("Error while Unmarshaling cashfree order: %v", err, resp, req)
		return order, err
	}

	if (cashFreeResponse["code"]) != nil {
		logrus.Errorf("Invalid Order Amount Error %v", err, resp, req)
		return order, err
	}

	order.Amount = int64((cashFreeResponse["order_amount"]).(float64))
	order.Currency = (cashFreeResponse["order_currency"]).(string)
	order.Entity = (cashFreeResponse["entity"]).(string)
	order.Status = (cashFreeResponse["order_status"]).(string)
	order.ID = (cashFreeResponse["order_token"]).(string)

	return order, nil
}

func VerifyPaymentOrder(orderId string) (responses.VerifyPaymentResponse, error) {

	var verifyPaymentResponse responses.VerifyPaymentResponse
	var body io.Reader

	url := os.Getenv("CASHFREE_BASE_URI") + "/orders/" + orderId + "/payments"

	req, err := http.NewRequest(http.MethodGet, url, body)

	if err != nil {
		logrus.Errorf("Error while verifying payment on cashfree: %v", err)
		return verifyPaymentResponse, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-client-id", os.Getenv("CASHFREE_APP_ID"))
	req.Header.Set("x-client-secret", os.Getenv("CASHFREE_SECRET_KEY"))
	req.Header.Set("x-api-version", os.Getenv("CASHFREE_API_VERSION"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error while making cashfree order: %v", err, resp, req)
		return verifyPaymentResponse, err
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error while making cashfree order: %v", err, resp, req)
		return verifyPaymentResponse, err
	}
	var cashFreeResponse []map[string]interface{}

	err = json.Unmarshal(response, &cashFreeResponse)
	if err != nil {
		logrus.Errorf("Error while Unmarshaling cashfree order: %v", err, resp, req)
		return verifyPaymentResponse, err
	}

	verifyPaymentResponse.Status = string(Failed)
	verifyPaymentResponse.OrderId = orderId

	for i := range cashFreeResponse {
		verifyPaymentResponse.PaymentId = fmt.Sprintf("%f", (cashFreeResponse[i]["cf_payment_id"]).(float64))
		if (cashFreeResponse[i]["payment_status"]).(string) == string(Success) {
			verifyPaymentResponse.Status = string(Success)
			break
		}
	}

	return verifyPaymentResponse, nil
}
