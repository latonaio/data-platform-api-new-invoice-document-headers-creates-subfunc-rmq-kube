package requests

type TotalNetAmount struct {
	InvoiceDocument *int     `json:"InvoiceDocument"`
	TotalNetAmount  *float32 `json:"TotalNetAmount"`
}
