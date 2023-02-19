package requests

type NetPaymentDays struct {
	InvoiceDocumentDate string  `json:"InvoiceDocumentDate"`
	PaymentDueDate      *string `json:"PaymentDueDate"`
	NetPaymentDays      *int    `json:"NetPaymentDays"`
}
