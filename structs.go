package goexch

import "time"

type CryptoCurrency string

const (
	Monero           CryptoCurrency = "XMR"
	Litecoin         CryptoCurrency = "LTC"
	Ethereum         CryptoCurrency = "ETH"
	Dash             CryptoCurrency = "DASH"
	BitcoinLightning CryptoCurrency = "BTCLN"
	Bitcoin          CryptoCurrency = "BTC"
	USDCoinErc20     CryptoCurrency = "USDC"
	TetherErc20      CryptoCurrency = "USDT"
	Dai              CryptoCurrency = "DAI"
)

// CreateOrderOptional holds optional parameters for creating an order.
type OrderOptions struct {
	// RefundAddress is the address for refunds in case of a failed exchange (Optional; used in REFUND_REQUEST state).
	RefundAddress string `json:"refund_address,omitempty"`
	// RateMode specifies the rate type, either "flat" or "dynamic" (Optional; default is "dynamic").
	RateMode string `json:"rate_mode,omitempty"`
	// ReferrerID is an identifier for referrals (Optional).
	ReferrerID string `json:"ref,omitempty"`
	// FeeOption specifies the network fee option: "s" for slow, "m" for medium, "f" for quick (Optional; default is "f").
	FeeOption string `json:"fee_option,omitempty"`
	// Aggregation indicates BTC aggregation preference: true for aggregated (receive/send), false for mixed, and omitted for default behavior (Optional).
	Aggregation *bool `json:"aggregation,omitempty"`
}

type CreateOrderResposnse struct {
	OrderID string `json:"orderid"`
}

type GetVolumeResponse struct {
	Bitcoin  *Volume `json:"BTC"`
	Btcln    *Volume `json:"BTCLN"`
	Dai      *Volume `json:"DAI"`
	Dash     *Volume `json:"DASH"`
	Eth      *Volume `json:"ETH"`
	Litecoin *Volume `json:"LTC"`
	Usdc     *Volume `json:"USDC"`
	Usdt     *Volume `json:"USDT"`
	Monero   *Volume `json:"XMR"`
}
type Volume struct {
	Volume string `json:"volume"`
}

type OrderResponse struct {
	Created        int            `json:"created"`
	FromAddr       string         `json:"from_addr"`
	AmountReceived *string        `json:"from_amount_received"`
	FromCurrency   CryptoCurrency `json:"from_currency"`
	MaxInput       string         `json:"max_input"`
	MinInput       string         `json:"min_input"`
	NetworkFee     int            `json:"network_fee"`
	Orderid        string         `json:"orderid"`
	Rate           string         `json:"rate"`
	RateMode       string         `json:"rate_mode"`
	State          string         `json:"state"`
	SvcFee         string         `json:"svc_fee"`
	ToAddress      string         `json:"to_address"`
	AmountSent     *string        `json:"to_amount"`
	ToCurrency     CryptoCurrency `json:"to_currency"`
	ReceivedID     *string        `json:"transaction_id_received"`
	SentID         *string        `json:"transaction_id_sent"`
}

func (od *OrderResponse) Date() time.Time {
	return time.Unix(int64(od.Created), 0)
}

type ResultResponse struct {
	Error  string `json:"error"`
	Result bool   `json:"result"`
}
