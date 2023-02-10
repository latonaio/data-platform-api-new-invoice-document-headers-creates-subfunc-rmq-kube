package requests

type DeliveryDocumentHeaderPartner struct {
	InvoiceDocument  *int   `json:"InvoiceDocument"`
	DeliveryDocument int    `json:"DeliveryDocument"`
	PartnerFunction  string `json:"PartnerFunction"`
	BusinessPartner  int    `json:"BusinessPartner"`
}
