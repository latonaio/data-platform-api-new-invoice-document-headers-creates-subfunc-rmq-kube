package requests

type OrdersItem struct {
	OrderID                                 int      `json:"OrderID"`
	OrderItem                               int      `json:"OrderItem"`
	OrderItemCategory                       string   `json:"OrderItemCategory"`
	SupplyChainRelationshipID               int      `json:"SupplyChainRelationshipID"`
	OrderItemText                           string   `json:"OrderItemText"`
	OrderItemTextByBuyer                    string   `json:"OrderItemTextByBuyer"`
	OrderItemTextBySeller                   string   `json:"OrderItemTextBySeller"`
	Product                                 string   `json:"Product"`
	ProductStandardID                       string   `json:"ProductStandardID"`
	ProductGroup                            *string  `json:"ProductGroup"`
	BaseUnit                                string   `json:"BaseUnit"`
	PricingDate                             string   `json:"PricingDate"`
	OrderQuantityInBaseUnit                 float32  `json:"OrderQuantityInBaseUnit"`
	OrderQuantityInDeliveryUnit             float32  `json:"OrderQuantityInDeliveryUnit"`
	NetAmount                               *float32 `json:"NetAmount"`
	TaxAmount                               *float32 `json:"TaxAmount"`
	GrossAmount                             *float32 `json:"GrossAmount"`
	Incoterms                               *string  `json:"Incoterms"`
	TransactionTaxClassification            string   `json:"TransactionTaxClassification"`
	ProductTaxClassificationBillToCountry   string   `json:"ProductTaxClassificationBillToCountry"`
	ProductTaxClassificationBillFromCountry string   `json:"ProductTaxClassificationBillFromCountry"`
	DefinedTaxClassification                string   `json:"DefinedTaxClassification"`
	PaymentTerms                            string   `json:"PaymentTerms"`
	PaymentMethod                           string   `json:"PaymentMethod"`
	Project                                 *string  `json:"Project"`
	ReferenceDocument                       *int     `json:"ReferenceDocument"`
	ReferenceDocumentItem                   *int     `json:"ReferenceDocumentItem"`
	TaxCode                                 *string  `json:"TaxCode"`
	TaxRate                                 *float32 `json:"TaxRate"`
}
