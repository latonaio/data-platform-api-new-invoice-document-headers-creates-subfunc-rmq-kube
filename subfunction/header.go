package subfunction

import (
	api_input_reader "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Input_Reader"
	api_processing_data_formatter "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Processing_Data_Formatter"
	"strings"
)

func (f *SubFunction) DeliveryDocumentHeaderData(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentHeaderData, error) {
	args := make([]interface{}, 0)

	deliveryDocument := psdc.DeliveryDocumentItem
	repeat := strings.Repeat("?,", len(deliveryDocument)-1) + "?"
	for _, tag := range deliveryDocument {
		args = append(args, tag.DeliveryDocument)
	}

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, SupplyChainRelationshipID, SupplyChainRelationshipDeliveryID, SupplyChainRelationshipDeliveryPlantID,
		SupplyChainRelationshipBillingID, SupplyChainRelationshipPaymentID, Buyer, Seller, DeliverToParty, DeliverFromParty, 
		DeliverToPlant, DeliverFromPlant, BillToParty, BillFromParty, BillToCountry, BillFromCountry, Payer, Payee, IsExportImport,
		OrderID, OrderItem, ContractType, OrderValidityStartDate, OrderValidityEndDate, GoodsIssueOrReceiptSlipNumber,
		Incoterms, TransactionCurrency
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
		WHERE DeliveryDocument IN ( `+repeat+` );`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocumentHeaderData(sdc, rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func CalculateInvoiceDocument(latestNumber int) *int {
	res := latestNumber + 1
	return &res
}

func (f *SubFunction) OrdersHeader(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrdersHeader, error) {
	args := make([]interface{}, 0)

	orderID := psdc.OrderID
	repeat := strings.Repeat("?,", len(orderID)-1) + "?"
	for _, tag := range orderID {
		args = append(args, tag.OrderID)
	}

	rows, err := f.db.Query(
		`SELECT OrderID, OrderType, SupplyChainRelationshipID, SupplyChainRelationshipBillingID,
		SupplyChainRelationshipPaymentID, Buyer, Seller, BillToParty, BillFromParty, BillToCountry, BillFromCountry,
		Payer, Payee, ContractType, OrderValidityStartDate, OrderValidityEndDate, InvoicePeriodStartDate,
		InvoicePeriodEndDate, TotalNetAmount, TotalTaxAmount, TotalGrossAmount, TransactionCurrency,
		PricingDate, Incoterms, PaymentTerms, PaymentMethod, IsExportImport
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE OrderID IN ( `+repeat+` );`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToOrdersHeader(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) CalculateInvoiceDocument(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) (*api_processing_data_formatter.CalculateInvoiceDocument, error) {
	metaData := psdc.MetaData
	dataKey := psdc.ConvertToInvoiceDocumentHeaderKey()

	dataKey.ServiceLabel = metaData.ServiceLabel

	rows, err := f.db.Query(
		`SELECT ServiceLabel, FieldNameWithNumberRange, LatestNumber
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_number_range_latest_number_data
		WHERE (ServiceLabel, FieldNameWithNumberRange) = (?, ?);`, dataKey.ServiceLabel, dataKey.FieldNameWithNumberRange,
	)
	if err != nil {
		return nil, err
	}

	dataQueryGets, err := psdc.ConvertToInvoiceDocumentHeaderQueryGets(sdc, rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	calculateInvoiceDocument := CalculateInvoiceDocument(*dataQueryGets.InvoiceDocumentLatestNumber)

	data := psdc.ConvertToCalculateInvoiceDocument(calculateInvoiceDocument)

	return data, err
}

func (f *SubFunction) TotalNetAmount(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.TotalNetAmount {
	totalNetAmount := float32(0)

	if psdc.ReferenceType.OrderID {
		for _, v := range psdc.OrdersItem {
			if v.NetAmount != nil {
				totalNetAmount += *v.NetAmount
			}
		}
	} else if psdc.ReferenceType.DeliveryDocument {
		for _, v := range psdc.DeliveryDocumentItemData {
			if v.NetAmount != nil {
				totalNetAmount += *v.NetAmount
			}
		}
	}

	if sdc.InvoiceDocument.TotalNetAmount != nil {
		totalNetAmount = *sdc.InvoiceDocument.TotalNetAmount
	}

	data := psdc.ConvertToTotalNetAmount(&totalNetAmount)

	return data
}
