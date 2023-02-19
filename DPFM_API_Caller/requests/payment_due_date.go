package requests

type PaymentDueDate struct {
	InvoiceDocumentDate string  `json:"invoiceDocumentDate"`
	PaymentDueDate      *string `json:"PaymentDueDate"`
}
