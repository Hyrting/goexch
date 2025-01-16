package goexch

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

var RateLimitExceeded = errors.New("rate limit exceeded, please wait")

// Client represents the API client with rate limiting.
type Client struct {
	baseURL     string
	apiKey      string
	client      *http.Client
	rateLimiter *RateLimiter // Added rate limiter
}

// New initializes and returns a new Client with rate limiting.
func New(key string) *Client {
	return &Client{
		baseURL:     "https://exch.cx/api",
		apiKey:      key,
		client:      &http.Client{},
		rateLimiter: nil,
	}
}

func (c *Client) Client(client *http.Client) {
	c.client = client
}

func (c *Client) RateLimiter(max int, interval time.Duration) {
	c.rateLimiter = NewRateLimiter(max, interval)
}

func (c *Client) request(path, method string, params map[string]string) (int, []byte, error) {
	// Enforce rate limiting
	if c.rateLimiter != nil {
		if !c.rateLimiter.Allow() {
			return 0, nil, RateLimitExceeded
		}
	}

	fullURL := fmt.Sprintf("%s/%s", c.baseURL, path)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return 0, []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	q := req.URL.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return res.StatusCode, body, nil
}

// Volume fetches 24-hour volume data.
func (c *Client) Volume() (*GetVolumeResponse, error) {
	statusCode, body, err := c.request("volume", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", statusCode)
	}

	var result *GetVolumeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Status retrieves network statuses.
func (c *Client) Status() (map[string]interface{}, error) {
	statusCode, body, err := c.request("status", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", statusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Order creates a new exchange order.
func (c *Client) Order(from, to CryptoCurrency, address string, opts *OrderOptions) (*CreateOrderResposnse, error) {
	if from == "" || to == "" || address == "" {
		return nil, fmt.Errorf("from, to, and address are required")
	}

	params := map[string]string{
		"from_currency": string(from),
		"to_currency":   string(to),
		"to_address":    address,
	}

	if opts != nil {
		if opts.RefundAddress != "" {
			params["refund_address"] = opts.RefundAddress
		}
		if opts.RateMode != "" {
			params["rate_mode"] = opts.RateMode
		}
		if opts.ReferrerID != "" {
			params["ref"] = opts.ReferrerID
		}
		if opts.FeeOption != "" {
			params["fee_option"] = opts.FeeOption
		}
		if opts.Aggregation != nil {
			params["aggregation"] = map[bool]string{true: "yes", false: "no"}[*opts.Aggregation]
		}
	}

	statusCode, body, err := c.request("create", http.MethodGet, params)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %d, response: %s", statusCode, string(body))
	}

	var result *CreateOrderResposnse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return result, nil
}

// GetOrder fetches order details.
func (c *Client) GetOrder(id string) (*OrderResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	params := map[string]string{"orderid": id}

	statusCode, body, err := c.request("order", http.MethodGet, params)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %d, response: %s", statusCode, string(body))
	}

	var result *OrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return result, nil
}

// Refund initiates a refund for an order.
func (c *Client) Refund(id string) (*ResultResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	params := map[string]string{"orderid": id}

	statusCode, body, err := c.request("order/refund", http.MethodGet, params)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %d, response: %s", statusCode, string(body))
	}

	var result *ResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return result, nil
}

// ConfirmRefund confirms a refund.
func (c *Client) ConfirmRefund(id string) (*ResultResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	params := map[string]string{"orderid": id}

	statusCode, body, err := c.request("order/refund_confirm", http.MethodGet, params)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %d, response: %s", statusCode, string(body))
	}

	var result *ResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return result, nil
}

// RevalidateAddress revalidates an address.
func (c *Client) RevalidateAddress(id, address string) (*ResultResponse, error) {
	if id == "" || address == "" {
		return nil, fmt.Errorf("id and address are required")
	}

	// Required parameters
	params := map[string]string{
		"orderid":    id,
		"to_address": address,
	}

	// Make the request
	statusCode, body, err := c.request("order/revalidate_address", http.MethodGet, params)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d, response: %s", statusCode, string(body))
	}

	var result *ResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

// Remove deletes order data.
func (c *Client) Remove(id string) (*ResultResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("order id is empty")
	}

	// Required parameters
	params := map[string]string{
		"orderid": id,
	}

	// Make the request
	statusCode, body, err := c.request("order/remove", http.MethodGet, params)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d, response: %s", statusCode, string(body))
	}

	var result *ResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}
