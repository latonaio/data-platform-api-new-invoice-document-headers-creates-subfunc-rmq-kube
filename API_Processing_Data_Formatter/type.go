package api_processing_data_formatter

type SDC struct {
	MetaData                   *MetaData                     `json:"MetaData"`
	ProcessType                *ProcessType                  `json:"ProcessType"`
	ReferenceType              *ReferenceType                `json:"ReferenceType"`
	OrderID                    []*OrderID                    `json:"OrderID"`
	OrderItem                  []*OrderItem                  `json:"OrderItem"`
	DeliveryDocumentHeader     []*DeliveryDocumentHeader     `json:"DeliveryDocumentHeader"`
	DeliveryDocumentItem       []*DeliveryDocumentItem       `json:"DeliveryDocumentItem"`
	OrdersHeader               []*OrdersHeader               `json:"OrdersHeader"`
	OrdersItem                 []*OrdersItem                 `json:"OrdersItem"`
	DeliveryDocumentHeaderData []*DeliveryDocumentHeaderData `json:"DeliveryDocumentHeaderData"`
	DeliveryDocumentItemData   []*DeliveryDocumentItemData   `json:"DeliveryDocumentItemData"`
	CalculateInvoiceDocument   []*CalculateInvoiceDocument   `json:"CalculateInvoiceDocument"`
	TotalNetAmount             *TotalNetAmount               `json:"TotalNetAmount"`
	TotalTaxAmount             *TotalTaxAmount               `json:"TotalTaxAmount"`
	TotalGrossAmount           *TotalGrossAmount             `json:"TotalGrossAmount"`
	InvoiceDocumentDate        *InvoiceDocumentDate          `json:"InvoiceDocumentDate"`
	PaymentTerms               []*PaymentTerms               `json:"PaymentTerms"`
	PaymentDueDate             *PaymentDueDate               `json:"PaymentDueDate"`
	NetPaymentDays             *NetPaymentDays               `json:"NetPaymentDays"`
	CreationDateHeader         *CreationDate                 `json:"CreationDateHeader"`
	LastChangeDateHeader       *LastChangeDate               `json:"LastChangeDateHeader"`
	CreationTimeHeader         *CreationTime                 `json:"CreationTimeHeader"`
	LastChangeTimeHeader       *LastChangeTime               `json:"LastChangeTimeHeader"`
}

type MetaData struct {
	BusinessPartnerID *int   `json:"business_partner"`
	ServiceLabel      string `json:"service_label"`
}

type ProcessType struct {
	BulkProcess       bool `json:"BulkProcess"`
	IndividualProcess bool `json:"IndividualProcess"`
	ArraySpec         bool `json:"ArraySpec"`
	RangeSpec         bool `json:"RangeSpec"`
}

type ReferenceType struct {
	OrderID          bool `json:"OrderID"`
	DeliveryDocument bool `json:"DeliveryDocument"`
}

type OrderIDKey struct {
	BillFromParty                   []*int `json:"BillFromParty"`
	BillFromPartyFrom               *int   `json:"BillFromPartyFrom"`
	BillFromPartyTo                 *int   `json:"BillFromPartyTo"`
	BillToParty                     []*int `json:"BillToParty"`
	BillToPartyFrom                 *int   `json:"BillToPartyFrom"`
	BillToPartyTo                   *int   `json:"BillToPartyTo"`
	HeaderCompleteDeliveryIsDefined bool   `json:"HeaderCompleteDeliveryIsDefined"`
	HeaderDeliveryStatus            string `json:"HeaderDeliveryStatus"`
	HeaderBillingStatus             string `json:"HeaderBillingStatus"`
	HeaderBillingBlockStatus        bool   `json:"HeaderBillingBlockStatus"`
	IsCancelled                     bool   `json:"IsCancelled"`
	IsMarkedForDeletion             bool   `json:"IsMarkedForDeletion"`
	ReferenceDocument               int    `json:"ReferenceDocument"`
}

type OrderID struct {
	OrderID                         int    `json:"OrderID"`
	BillFromParty                   *int   `json:"BillFromParty"`
	BillToParty                     *int   `json:"BillToParty"`
	HeaderCompleteDeliveryIsDefined *bool  `json:"HeaderCompleteDeliveryIsDefined"`
	HeaderDeliveryStatus            string `json:"HeaderDeliveryStatus"`
	HeaderBillingStatus             string `json:"HeaderBillingStatus"`
	HeaderBillingBlockStatus        *bool  `json:"HeaderBillingBlockStatus"`
	IsCancelled                     *bool  `json:"IsCancelled"`
	IsMarkedForDeletion             *bool  `json:"IsMarkedForDeletion"`
}

type OrderItemKey struct {
	OrderID                       []int  `json:"OrderID"`
	ItemCompleteDeliveryIsDefined bool   `json:"ItemCompleteDeliveryIsDefined"`
	ItemDeliveryStatus            string `json:"ItemDeliveryStatus"`
	ItemBillingStatus             string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus        bool   `json:"ItemBillingBlockStatus"`
	IsCancelled                   bool   `json:"IsCancelled"`
	IsMarkedForDeletion           bool   `json:"IsMarkedForDeletion"`
}

type OrderItem struct {
	OrderID                       int     `json:"OrderID"`
	OrderItem                     int     `json:"OrderItem"`
	ItemCompleteDeliveryIsDefined *bool   `json:"ItemCompleteDeliveryIsDefined"`
	ItemDeliveryStatus            *string `json:"ItemDeliveryStatus"`
	ItemBillingStatus             *string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus        *bool   `json:"ItemBillingBlockStatus"`
	IsCancelled                   *bool   `json:"IsCancelled"`
	IsMarkedForDeletion           *bool   `json:"IsMarkedForDeletion"`
}

type DeliveryDocumentHeaderKey struct {
	BillFromParty                   []*int `json:"BillFromParty"`
	BillFromPartyFrom               *int   `json:"BillFromPartyFrom"`
	BillFromPartyTo                 *int   `json:"BillFromPartyTo"`
	BillToParty                     []*int `json:"BillToParty"`
	BillToPartyFrom                 *int   `json:"BillToPartyFrom"`
	BillToPartyTo                   *int   `json:"BillToPartyTo"`
	HeaderCompleteDeliveryIsDefined bool   `json:"HeaderCompleteDeliveryIsDefined"`
	HeaderDeliveryStatus            string `json:"HeaderDeliveryStatus"`
	HeaderBillingStatus             string `json:"HeaderBillingStatus"`
	HeaderBillingBlockStatus        bool   `json:"HeaderBillingBlockStatus"`
	IsCancelled                     bool   `json:"IsCancelled"`
	IsMarkedForDeletion             bool   `json:"IsMarkedForDeletion"`
	ReferenceDocument               int    `json:"ReferenceDocument"`
}

type DeliveryDocumentHeader struct {
	DeliveryDocument                int    `json:"DeliveryDocument"`
	BillFromParty                   *int   `json:"BillFromParty"`
	BillToParty                     *int   `json:"BillToParty"`
	HeaderCompleteDeliveryIsDefined *bool  `json:"HeaderCompleteDeliveryIsDefined"`
	HeaderDeliveryStatus            string `json:"HeaderDeliveryStatus"`
	HeaderBillingStatus             string `json:"HeaderBillingStatus"`
	HeaderBillingBlockStatus        *bool  `json:"HeaderBillingBlockStatus"`
	IsCancelled                     *bool  `json:"IsCancelled"`
	IsMarkedForDeletion             *bool  `json:"IsMarkedForDeletion"`
}

type DeliveryDocumentItemKey struct {
	DeliveryDocument              []int     `json:"DeliveryDocument"`
	ConfirmedDeliveryDate         []*string `json:"ConfirmedDeliveryDate"`
	ConfirmedDeliveryDateFrom     *string   `json:"BillFromPartyFrom"`
	ConfirmedDeliveryDateTo       *string   `json:"BillFromPartyTo"`
	ActualGoodsIssueDate          []*string `json:"ActualGoodsIssueDate"`
	ActualGoodsIssueDateFrom      *string   `json:"BillToPartyFrom"`
	ActualGoodsIssueDateTo        *string   `json:"BillToPartyTo"`
	ItemCompleteDeliveryIsDefined bool      `json:"ItemCompleteDeliveryIsDefined"`
	// ItemDeliveryStatus            string    `json:"ItemDeliveryStatus"`
	ItemBillingStatus      string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus bool   `json:"ItemBillingBlockStatus"`
	IsCancelled            bool   `json:"IsCancelled"`
	IsMarkedForDeletion    bool   `json:"IsMarkedForDeletion"`
}

type DeliveryDocumentItem struct {
	DeliveryDocument              int     `json:"DeliveryDocument"`
	DeliveryDocumentItem          int     `json:"DeliveryDocumentItem"`
	ConfirmedDeliveryDate         *string `json:"ConfirmedDeliveryDate"`
	ActualGoodsIssueDate          *string `json:"ActualGoodsIssueDate"`
	ItemCompleteDeliveryIsDefined *bool   `json:"ItemCompleteDeliveryIsDefined"`
	// ItemDeliveryStatus            string  `json:"ItemDeliveryStatus"`
	ItemBillingStatus      *string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus *bool   `json:"ItemBillingBlockStatus"`
	IsCancelled            *bool   `json:"IsCancelled"`
	IsMarkedForDeletion    *bool   `json:"IsMarkedForDeletion"`
}

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

type BillParty struct {
	BillFromParty        int `json:"BillFromParty"`
	BillToParty          int `json:"BillToParty"`
	OrderID              int `json:"OrderID"`
	OrderItem            int `json:"OrderItem"`
	DeliveryDocument     int `json:"DeliveryDocument"`
	DeliveryDocumentItem int `json:"DeliveryDocumentItem"`
}

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

type TotalNetAmountQueryGets struct {
	InvoiceDocument *int     `json:"InvoiceDocument"`
	TotalNetAmount  *float32 `json:"TotalNetAmount"`
}

type TotalNetAmount struct {
	InvoiceDocument *int     `json:"InvoiceDocument"`
	TotalNetAmount  *float32 `json:"TotalNetAmount"`
}

type TotalTaxAmount struct {
	InvoiceDocument *int     `json:"invoiceDocument"`
	TotalTaxAmount  *float32 `json:"TotalTaxAmount"`
}

type TotalGrossAmount struct {
	InvoiceDocument  *int     `json:"invoiceDocument"`
	TotalGrossAmount *float32 `json:"TotalGrossAmount"`
}

type InvoiceDocumentDate struct {
	InvoiceDocument     *int    `json:"invoiceDocument"`
	InvoiceDocumentDate *string `json:"invoiceDocumentDate"`
}

type PaymentDueDate struct {
	InvoiceDocumentDate string  `json:"invoiceDocumentDate"`
	PaymentDueDate      *string `json:"PaymentDueDate"`
}

type NetPaymentDays struct {
	InvoiceDocumentDate string  `json:"InvoiceDocumentDate"`
	PaymentDueDate      *string `json:"PaymentDueDate"`
	NetPaymentDays      *int    `json:"NetPaymentDays"`
}

type PaymentTerms struct {
	PaymentTerms                string `json:"PaymentTerms"`
	BaseDate                    int    `json:"BaseDate"`
	BaseDateCalcAddMonth        *int   `json:"BaseDateCalcAddMonth"`
	BaseDateCalcFixedDate       *int   `json:"BaseDateCalcFixedDate"`
	PaymentDueDateCalcAddMonth  *int   `json:"PaymentDueDateCalcAddMonth"`
	PaymentDueDateCalcFixedDate *int   `json:"PaymentDueDateCalcFixedDate"`
}

//日付等の処理

type CreationDate struct {
	CreationDate string `json:"CreationDate"`
}

type LastChangeDate struct {
	LastChangeDate string `json:"LastChangeDate"`
}

type CreationTime struct {
	CreationTime string `json:"CreationTime"`
}

type LastChangeTime struct {
	LastChangeTime string `json:"LastChangeTime"`
}
