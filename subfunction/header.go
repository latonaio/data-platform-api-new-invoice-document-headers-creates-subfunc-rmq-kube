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

// func (f *SubFunction) TotalNetAmount(
// 	sdc *api_input_reader.SDC,
// 	psdc *api_processing_data_formatter.SDC,
// ) (*api_processing_data_formatter.TotalNetAmount, error) {
// 	var err error
// 	// オーダー参照
// 	// TODO: nullの場合どうする？
// 	if sdc.InvoiceDocument.TotalNetAmount != nil {
// 		for i, v := range *psdc.HeaderOrdersHeader {
// 			if v.TotalNetAmount == *sdc.InvoiceDocument.TotalNetAmount {
// 				data = psdc.ConvertToTotalNetAmount(&v.TotalNetAmount)
// 				break
// 			}
// 			if i == len(*psdc.HeaderOrdersHeader)-1 {
// 				return nil, xerrors.Errorf("TotalNetAmountが一致しません。")
// 			}
// 		}
// 	}

// 	// // 入出荷伝票参照
// 	// rows, err := f.db.Query(
// 	// 	`SELECT InvoiceDocument, TotalNetAmount
// 	// 	FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_invoice_document_header_data
// 	// 	WHERE InvoiceDocument = ?;`, sdc.InvoiceDocument.InvoiceDocument,
// 	// )
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// dataQueryGets, err := psdc.ConvertToTotalNetAmountQueryGets(sdc, rows)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// if sdc.InvoiceDocument.TotalNetAmount != nil {
// 	// 	if *dataQueryGets.TotalNetAmount == *sdc.InvoiceDocument.TotalNetAmount {
// 	// 		data, err = psdc.ConvertToTotalNetAmount(dataQueryGets.TotalNetAmount)
// 	// 	} else {
// 	// 		return nil, xerrors.Errorf("TotalNetAmountが一致しません。")
// 	// 	}
// 	// }

// 	return data, err
// }
