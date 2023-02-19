package requests

type CalculateInvoiceDocument struct {
	InvoiceDocumentLatestNumber *int `json:"InvoiceDocumentLatestNumber"`
	InvoiceDocument             int  `json:"InvoiceDocument"`
	OrderID                     int  `json:"OrderID"`
	OrderItem                   int  `json:"OrderItem"`
	DeliveryDocument            int  `json:"DeliveryDocument"`
	DeliveryDocumentItem        int  `json:"DeliveryDocumentItem"`
	BillFromParty               int  `json:"BillFromParty"`
	BillToParty                 int  `json:"BillToParty"`
}
