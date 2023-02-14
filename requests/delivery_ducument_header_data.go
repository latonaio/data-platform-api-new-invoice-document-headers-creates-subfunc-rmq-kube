package requests

type DeliveryDocumentHeaderData struct {
	DeliveryDocument                       int     `json:"DeliveryDocument"`
	SupplyChainRelationshipID              int     `json:"SupplyChainRelationshipID"`
	SupplyChainRelationshipDeliveryID      int     `json:"SupplyChainRelationshipDeliveryID"`
	SupplyChainRelationshipDeliveryPlantID int     `json:"SupplyChainRelationshipDeliveryPlantID"`
	SupplyChainRelationshipBillingID       *int    `json:"SupplyChainRelationshipBillingID"`
	SupplyChainRelationshipPaymentID       *int    `json:"SupplyChainRelationshipPaymentID"`
	Buyer                                  int     `json:"Buyer"`
	Seller                                 int     `json:"Seller"`
	DeliverToParty                         int     `json:"DeliverToParty"`
	DeliverFromParty                       int     `json:"DeliverFromParty"`
	DeliverToPlant                         string  `json:"DeliverToPlant"`
	DeliverFromPlant                       string  `json:"DeliverFromPlant"`
	BillToParty                            *int    `json:"BillToParty"`
	BillFromParty                          *int    `json:"BillFromParty"`
	BillToCountry                          *string `json:"BillToCountry"`
	BillFromCountry                        *string `json:"BillFromCountry"`
	Payer                                  *int    `json:"Payer"`
	Payee                                  *int    `json:"Payee"`
	IsExportImport                         *bool   `json:"IsExportImport"`
	OrderID                                *int    `json:"OrderID"`
	OrderItem                              *int    `json:"OrderItem"`
	ContractType                           *string `json:"ContractType"`
	OrderValidityStartDate                 *string `json:"OrderValidityStartDate"`
	OrderValidityEndDate                   *string `json:"OrderValidityEndDate"`
	GoodsIssueOrReceiptSlipNumber          *string `json:"GoodsIssueOrReceiptSlipNumber"`
	Incoterms                              *string `json:"Incoterms"`
	TransactionCurrency                    *string `json:"TransactionCurrency"`
}
