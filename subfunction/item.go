package subfunction

import (
	api_input_reader "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Input_Reader"
	api_processing_data_formatter "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Processing_Data_Formatter"
	"strings"
)

func (f *SubFunction) DeliveryDocumentItemData(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentItemData, error) {
	args := make([]interface{}, 0)

	deliveryDocument := psdc.DeliveryDocumentItem
	repeat := strings.Repeat("(?, ?),", len(deliveryDocument)-1) + "(?, ?)"
	for _, tag := range deliveryDocument {
		args = append(args, tag.DeliveryDocument, tag.DeliveryDocumentItem)
	}

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, DeliveryDocumentItem, DeliveryDocumentItemCategory, SupplyChainRelationshipID, SupplyChainRelationshipDeliveryID,
		SupplyChainRelationshipDeliveryPlantID, SupplyChainRelationshipBillingID, SupplyChainRelationshipPaymentID, Buyer, Seller, DeliverToParty,
		DeliverFromParty, DeliverToPlant, DeliverFromPlant, BillToParty, BillFromParty, BillToCountry, BillFromCountry, Payer, Payee, 
		DeliverToPlantStorageLocation, DeliverFromPlantStorageLocation, ProductionPlantBusinessPartner, ProductionPlant,
		ProductionPlantStorageLocation, DeliveryDocumentItemText, DeliveryDocumentItemTextByBuyer, DeliveryDocumentItemTextBySeller, Product,
		ProductStandardID, ProductGroup, BaseUnit, ActualGoodsIssueDate, ActualGoodsIssueTime, ActualGoodsReceiptDate, ActualGoodsReceiptTime,
		ActualGoodsIssueQuantity, ActualGoodsIssueQtyInBaseUnit, ActualGoodsReceiptQuantity, ActualGoodsReceiptQtyInBaseUnit, ItemGrossWeight,
		ItemNetWeight, ItemWeightUnit, NetAmount, TaxAmount, GrossAmount, OrderID, OrderItem, OrderType, ContractType, OrderValidityStartDate,
		OrderValidityEndDate, PaymentTerms, PaymentMethod, InvoicePeriodStartDate, InvoicePeriodEndDate, Project, ReferenceDocument, ReferenceDocumentItem,
		TransactionTaxClassification, ProductTaxClassificationBillToCountry, ProductTaxClassificationBillFromCountry, DefinedTaxClassification,
		TaxCode, TaxRate, CountryOfOrigin, CountryOfOriginLanguage
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_item_data
		WHERE (DeliveryDocument, DeliveryDocumentItem)  IN ( `+repeat+` );`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocumentItemData(sdc, rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrdersItem(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrdersItem, error) {
	args := make([]interface{}, 0)

	orderID := psdc.OrderItem

	repeat := strings.Repeat("(?, ?),", len(orderID)-1) + "(?, ?)"
	for _, v := range orderID {
		args = append(args, v.OrderID, v.OrderItem)
	}

	rows, err := f.db.Query(
		`SELECT OrderID, OrderItem, OrderItemCategory, SupplyChainRelationshipID, OrderItemText, OrderItemTextByBuyer,
		OrderItemTextBySeller,Product, ProductStandardID, ProductGroup, BaseUnit, PricingDate, OrderQuantityInBaseUnit,
		OrderQuantityInDeliveryUnit, NetAmount, TaxAmount, GrossAmount, Incoterms, TransactionTaxClassification,
		ProductTaxClassificationBillToCountry, ProductTaxClassificationBillFromCountry, DefinedTaxClassification,
		PaymentTerms, PaymentMethod, Project, ReferenceDocument, ReferenceDocumentItem, TaxCode, TaxRate
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_item_data
		WHERE (OrderID, OrderItem) IN ( `+repeat+` );`, args...,
	)
	if err != nil {
		return nil, err
	}

	data, err := psdc.ConvertToOrdersItem(sdc, rows)
	if err != nil {
		return nil, err
	}

	return data, err
}
