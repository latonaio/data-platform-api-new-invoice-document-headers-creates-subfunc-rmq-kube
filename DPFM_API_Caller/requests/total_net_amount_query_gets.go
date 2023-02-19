package requests

type TotalNetAmountQueryGets struct {
	InvoiceDocument *int     `json:"InvoiceDocument"`
	TotalNetAmount  *float32 `json:"TotalNetAmount"`
}
