package requests

type OrdersHeader struct {
	OrderID                          int     `json:"OrderID"`
	OrderType                        string  `json:"OrderType"`
	SupplyChainRelationshipID        int     `json:"SupplyChainRelationshipID"`
	SupplyChainRelationshipBillingID *int    `json:"SupplyChainRelationshipBillingID"`
	SupplyChainRelationshipPaymentID *int    `json:"SupplyChainRelationshipPaymentID"`
	Buyer                            int     `json:"Buyer"`
	Seller                           int     `json:"Seller"`
	BillToParty                      *int    `json:"BillToParty"`
	BillFromParty                    *int    `json:"BillFromParty"`
	BillToCountry                    *string `json:"BillToCountry"`
	BillFromCountry                  *string `json:"BillFromCountry"`
	Payer                            *int    `json:"Payer"`
	Payee                            *int    `json:"Payee"`
	ContractType                     *string `json:"ContractType"`
	OrderValidityStartDate           *string `json:"OrderVaridityStartDate"`
	OrderValidityEndDate             *string `json:"OrderValidityEndDate"`
	InvoicePeriodStartDate           *string `json:"InvoicePeriodStartDate"`
	InvoicePeriodEndDate             *string `json:"InvoicePeriodEndDate"`
	TotalNetAmount                   float32 `json:"TotalNetAmount"`
	TotalTaxAmount                   float32 `json:"TotalTaxAmount"`
	TotalGrossAmount                 float32 `json:"TotalGrossAmount"`
	TransactionCurrency              string  `json:"TransactionCurrency"`
	PricingDate                      string  `json:"PricingDate"`
	Incoterms                        *string `json:"Incoterms"`
	PaymentTerms                     string  `json:"PaymentTerms"`
	PaymentMethod                    string  `json:"PaymentMethod"`
	IsExportImport                   *bool   `json:"IsExportImport"`
}
