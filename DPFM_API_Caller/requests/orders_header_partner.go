package requests

type OrdersHeaderPartner struct {
	InvoiceDocument *int   `json:"InvoiceDocument"`
	OrderID         int    `json:"OrderID"`
	PartnerFunction string `json:"PartnerFunction"`
	BusinessPartner int    `json:"BusinessPartner"`
}
