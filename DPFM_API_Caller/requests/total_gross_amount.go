package requests

type TotalGrossAmount struct {
	InvoiceDocument  *int     `json:"InvoiceDocument"`
	TotalGrossAmount *float32 `json:"TotalGrossAmount"`
}
