package requests

type TotalTaxAmount struct {
	InvoiceDocument *int     `json:"InvoiceDocument"`
	TotalTaxAmount  *float32 `json:"TotalTaxAmount"`
}
