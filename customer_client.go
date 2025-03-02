package metering

import (
	"encoding/json"
	"errors"
	"fmt"
)

type CustomerClient struct {
	BaseClient
}

type Customer struct {
	CustomerId    string            `json:"customerId"`
	CustomerName  string            `json:"customerName"`
	CustomerEmail string            `json:"customerEmail"`
	Traits        map[string]string `json:"traits,omitempty"`
	Enabled       bool              `json:"enabled"`
	UpdateTime    int64             `json:"updateTime,omitempty"`
	CreateTime    int64             `json:"createTime,omitempty"`
}

func NewCustomerClient(apiKey string, opts ...ClientOption) *CustomerClient {
	bc := NewBaseClient(apiKey, opts...)
	c := &CustomerClient{BaseClient: *bc}
	c.logf("Instantiating amberflo.io Signals Client")
	return c
}

func (m *CustomerClient) AddorUpdateCustomer(customer *Customer, createInStripe bool) (*Customer, error) {
	if customer.CustomerId == "" || customer.CustomerName == "" {
		return nil, errors.New("customer info 'CustomerId' and 'CustomerName' are required fields")
	}

	return m.sendCustomerToApi(customer, createInStripe)
}

func (c *CustomerClient) GetCustomer(customerId string) (*Customer, error) {
	signature := fmt.Sprintf("GetCustomer(%s)", customerId)
	var customer *Customer
	urlGet := fmt.Sprintf("%s/customers/?customerId=%s", Endpoint, customerId)
	data, err := c.AmberfloHttpClient.sendHttpRequest("Customers", urlGet, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err)
	}
	if data != nil && string(data) != "{}" {
		err = json.Unmarshal(data, &customer)
		if err != nil {
			return nil, fmt.Errorf("%s Error reading JSON body: %s", signature, err)
		}
	}

	return customer, nil
}

func (c *CustomerClient) sendCustomerToApi(payload *Customer, createInStripe bool) (*Customer, error) {
	signature := fmt.Sprintf("sendCustomerToApi(%v)", payload)

	c.logf("Checking if customer deatils exist %s", payload.CustomerId)
	customer, _ := c.GetCustomer(payload.CustomerId)

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("%s error marshalling payload: %s", signature, err)
	}

	url := fmt.Sprintf("%s/customers", Endpoint)
	httpMethod := ""
	if customer != nil && customer.CustomerId == payload.CustomerId {
		httpMethod = "PUT"
	} else {
		httpMethod = "POST"
		url = fmt.Sprintf("%s/customers?autoCreateCustomerInStripe=%t", Endpoint, createInStripe)
	}
	b, err = c.AmberfloHttpClient.sendHttpRequest("customers", url, httpMethod, b)
	if err != nil {
		return nil, fmt.Errorf("%s error making %s http call: %s", signature, httpMethod, err)
	}

	if b != nil {
		err = json.Unmarshal(b, &customer)
		if err != nil {
			return nil, fmt.Errorf("%s Error reading JSON body: %s", signature, err)
		}
	}

	return customer, nil
}
