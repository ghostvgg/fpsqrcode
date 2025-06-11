package qr

type PaymentOperator struct {
	GlobalUniqueIdentifier string            `json:"global_unique_identifier" binding:"max=32"`
	ExtraFields            map[string]string `json:"extra_fields"` // optional nested key-values
}

type QRRequest struct {
	FPSID           string          `json:"fps_id" binding:"required"`
	MerchantName    string          `json:"merchant_name" binding:"max=25"`
	City            string          `json:"city" binding:"required,max=15"`
	Dynamic         bool            `json:"dynamic"`                           // true for dynamic QR, false for static
	Amount          string          `json:"amount" binding:"max=13"`           // optional unless dynamic
	Currency        string          `json:"currency" binding:"required,len=3"` // optional; default HKD
	MerchantTimeout string          `json:"merchant_timeout"`
	ReferenceLabel  string          `json:"reference_label,omitempty" binding:"max=25"` // optional
	BillNumber      string          `json:"bill_number,omitempty" binding:"max=25"`     // optional
	PaymentOperator PaymentOperator `json:"payment_operator,omitempty"`                 // optional
}
