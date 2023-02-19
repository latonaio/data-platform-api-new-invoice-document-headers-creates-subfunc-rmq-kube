package requests

type DeliveryDocumentItemData struct {
	DeliveryDocument                        int      `json:"DeliveryDocument"`
	DeliveryDocumentItem                    int      `json:"DeliveryDocumentItem"`
	DeliveryDocumentItemCategory            *string  `json:"DeliveryDocumentItemCategory"`
	SupplyChainRelationshipID               int      `json:"SupplyChainRelationshipID"`
	SupplyChainRelationshipDeliveryID       int      `json:"SupplyChainRelationshipDeliveryID"`
	SupplyChainRelationshipDeliveryPlantID  int      `json:"SupplyChainRelationshipDeliveryPlantID"`
	SupplyChainRelationshipBillingID        *int     `json:"SupplyChainRelationshipBillingID"`
	SupplyChainRelationshipPaymentID        *int     `json:"SupplyChainRelationshipPaymentID"`
	Buyer                                   int      `json:"Buyer"`
	Seller                                  int      `json:"Seller"`
	DeliverToParty                          int      `json:"DeliverToParty"`
	DeliverFromParty                        int      `json:"DeliverFromParty"`
	DeliverToPlant                          string   `json:"DeliverToPlant"`
	DeliverFromPlant                        string   `json:"DeliverFromPlant"`
	BillToParty                             *int     `json:"BillToParty"`
	BillFromParty                           *int     `json:"BillFromParty"`
	BillToCountry                           *string  `json:"BillToCountry"`
	BillFromCountry                         *string  `json:"BillFromCountry"`
	Payer                                   *int     `json:"Payer"`
	Payee                                   *int     `json:"Payee"`
	DeliverToPlantStorageLocation           *string  `json:"DeliverToPlantStorageLocation"`
	DeliverFromPlantStorageLocation         *string  `json:"DeliverFromPlantStorageLocation"`
	ProductionPlantBusinessPartner          *int     `json:"ProductionPlantBusinessPartner"`
	ProductionPlant                         *string  `json:"ProductionPlant"`
	ProductionPlantStorageLocation          *string  `json:"ProductionPlantStorageLocation"`
	DeliveryDocumentItemText                *string  `json:"DeliveryDocumentItemText"`
	DeliveryDocumentItemTextByBuyer         string   `json:"DeliveryDocumentItemTextByBuyer"`
	DeliveryDocumentItemTextBySeller        string   `json:"DeliveryDocumentItemTextBySeller"`
	Product                                 *string  `json:"Product"`
	ProductStandardID                       *string  `json:"ProductStandardID"`
	ProductGroup                            *string  `json:"ProductGroup"`
	BaseUnit                                *string  `json:"BaseUnit"`
	ActualGoodsIssueDate                    *string  `json:"ActualGoodsIssueDate"`
	ActualGoodsIssueTime                    *string  `json:"ActualGoodsIssueTime"`
	ActualGoodsReceiptDate                  *string  `json:"ActualGoodsReceiptDate"`
	ActualGoodsReceiptTime                  *string  `json:"ActualGoodsReceiptTime"`
	ActualGoodsIssueQuantity                *float32 `json:"ActualGoodsIssueQuantity"`
	ActualGoodsIssueQtyInBaseUnit           *float32 `json:"ActualGoodsIssueQtyInBaseUnit"`
	ActualGoodsReceiptQuantity              *float32 `json:"ActualGoodsReceiptQuantity"`
	ActualGoodsReceiptQtyInBaseUnit         *float32 `json:"ActualGoodsReceiptQtyInBaseUnit"`
	ItemGrossWeight                         *float32 `json:"ItemGrossWeight"`
	ItemNetWeight                           *float32 `json:"ItemNetWeight"`
	ItemWeightUnit                          *string  `json:"ItemWeightUnit"`
	NetAmount                               *float32 `json:"NetAmount"`
	TaxAmount                               *float32 `json:"TaxAmount"`
	GrossAmount                             *float32 `json:"GrossAmount"`
	OrderID                                 *int     `json:"OrderID"`
	OrderItem                               *int     `json:"OrderItem"`
	OrderType                               *string  `json:"OrderType"`
	ContractType                            *string  `json:"ContractType"`
	OrderValidityStartDate                  *string  `json:"OrderValidityStartDate"`
	OrderValidityEndDate                    *string  `json:"OrderValidityEndDate"`
	PaymentTerms                            *string  `json:"PaymentTerms"`
	PaymentMethod                           *string  `json:"PaymentMethod"`
	InvoicePeriodStartDate                  *string  `json:"InvoicePeriodStartDate"`
	InvoicePeriodEndDate                    *string  `json:"InvoicePeriodEndDate"`
	Project                                 *string  `json:"Project"`
	ReferenceDocument                       *int     `json:"ReferenceDocument"`
	ReferenceDocumentItem                   *int     `json:"ReferenceDocumentItem"`
	TransactionTaxClassification            string   `json:"TransactionTaxClassification"`
	ProductTaxClassificationBillToCountry   string   `json:"ProductTaxClassificationBillToCountry"`
	ProductTaxClassificationBillFromCountry string   `json:"ProductTaxClassificationBillFromCountry"`
	DefinedTaxClassifications               string   `json:"DefinedTaxClassification"`
	TaxCode                                 *string  `json:"TaxCode"`
	TaxRate                                 *float32 `json:"TaxRate"`
	CountryOfOrigin                         *string  `json:"CountryOfOrigin"`
	CountryOfOriginLanguage                 *string  `json:"CountryOfOriginLanguage"`
}