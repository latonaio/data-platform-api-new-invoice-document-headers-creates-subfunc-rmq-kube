package dpfm_api_output_formatter

import (
	api_input_reader "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Input_Reader"
	api_processing_data_formatter "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Processing_Data_Formatter"
	"encoding/json"
	"reflect"

	"golang.org/x/xerrors"
)

func ConvertToHeader(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*Header, error) {
	var err error
	ordersHeaderMap := StructArrayToMap(psdc.OrdersHeader, "OrderID")
	deliveryDocumentHeaderMap := StructArrayToMap(psdc.DeliveryDocumentHeaderData, "DeliveryDocument")
	referenceType := psdc.ReferenceType

	headers := make([]*Header, 0)
	for i, invoiceDocument := range psdc.CalculateInvoiceDocument {
		header := &Header{}
		inputHeader := sdc.Header

		// 入力ファイル
		header, err = jsonTypeConversion(header, inputHeader)
		if err != nil {
			return nil, xerrors.Errorf("request create error: %w", err)
		}

		if referenceType.OrderID {
			orderID := invoiceDocument.OrderID

			// 1-1
			header, err = jsonTypeConversion(header, ordersHeaderMap[orderID])
			if err != nil {
				return nil, xerrors.Errorf("request create error: %w", err)
			}
		} else if referenceType.DeliveryDocument {
			deliveryDocument := invoiceDocument.DeliveryDocument
			deliveryDocumentItem := invoiceDocument.DeliveryDocumentItem
			//2-1
			header, err = jsonTypeConversion(header, deliveryDocumentHeaderMap[deliveryDocument])
			if err != nil {
				return nil, xerrors.Errorf("request create error: %w", err)
			}
			//2-2
			for _, deliveryDocumentItemData := range psdc.DeliveryDocumentItemData {
				if deliveryDocumentItemData.DeliveryDocument == deliveryDocument || deliveryDocumentItemData.DeliveryDocumentItem == deliveryDocumentItem {
					header.SupplyChainRelationshipBillingID = *deliveryDocumentItemData.SupplyChainRelationshipBillingID
					header.SupplyChainRelationshipPaymentID = *deliveryDocumentItemData.SupplyChainRelationshipPaymentID
					header.BillToParty = *deliveryDocumentItemData.BillToParty
					header.BillFromParty = *deliveryDocumentItemData.BillFromParty
					header.BillToCountry = *deliveryDocumentItemData.BillToCountry
					header.BillFromCountry = *deliveryDocumentItemData.BillFromCountry
					header.Payer = *deliveryDocumentItemData.Payer
					header.Payee = *deliveryDocumentItemData.Payee
					header.InvoicePeriodStartDate = *deliveryDocumentItemData.InvoicePeriodStartDate
					header.InvoicePeriodEndDate = *deliveryDocumentItemData.InvoicePeriodEndDate
					header.PaymentTerms = deliveryDocumentItemData.PaymentTerms
					header.PaymentMethod = deliveryDocumentItemData.PaymentMethod
				}
			}
		}

		header.InvoiceDocument = psdc.CalculateInvoiceDocument[i].InvoiceDocument
		header.CreationDate = psdc.CreationDateHeader.CreationDate
		header.CreationTime = psdc.CreationTimeHeader.CreationTime
		header.LastChangeDate = psdc.LastChangeDateHeader.LastChangeDate
		header.LastChangeTime = psdc.LastChangeTimeHeader.LastChangeTime
		header.InvoiceDocumentDate = *psdc.InvoiceDocumentDate.InvoiceDocumentDate
		//header.InvoiceDocumentTime
		//header.AccountingPostingDate
		header.HeaderBillingIsConfirmed = getBoolPtr(false)
		header.HeaderBillingConfStatus = getStringPtr("NP")
		header.TotalNetAmount = psdc.TotalNetAmount.TotalNetAmount
		header.TotalTaxAmount = psdc.TotalTaxAmount.TotalTaxAmount
		header.TotalGrossAmount = psdc.TotalGrossAmount.TotalGrossAmount
		// header.DueCalculationBaseDate = //TBD
		header.PaymentDueDate = psdc.PaymentDueDate.PaymentDueDate
		header.NetPaymentDays = psdc.NetPaymentDays.NetPaymentDays
		header.ExternalReferenceDocument = sdc.Header.ExternalReferenceDocument
		header.DocumentHeaderText = sdc.Header.DocumentHeaderText
		header.HeaderIsCleared = getBoolPtr(false)
		header.HeaderPaymentBlockStatus = getBoolPtr(false)
		header.HeaderPaymentRequisitionIsCreated = getBoolPtr(false)
		header.IsCancelled = getBoolPtr(false)

		headers = append(headers, header)
	}

	return headers, nil
}

func getBoolPtr(b bool) *bool {
	return &b
}

func getStringPtr(s string) *string {
	return &s
}

func StructArrayToMap[T any](data []T, key string) map[any]T {
	res := make(map[any]T, len(data))

	for _, value := range data {
		m := StructToMap[T](&value, key)
		for k, v := range m {
			res[k] = v
		}
	}

	return res
}

func StructToMap[T any](data interface{}, key string) map[any]T {
	res := make(map[any]T)
	elem := reflect.Indirect(reflect.ValueOf(data).Elem())
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Name
		if field == key {
			rv := reflect.ValueOf(elem.Field(i).Interface())
			if rv.Kind() == reflect.Ptr {
				if rv.IsNil() {
					return nil
				}
			}
			value := reflect.Indirect(elem.Field(i)).Interface()
			var dist T
			res[value], _ = jsonTypeConversion(dist, elem.Interface())
			break
		}
	}

	return res
}

func jsonTypeConversion[T any](dist T, data interface{}) (T, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return dist, xerrors.Errorf("Marshal error: %w", err)
	}
	err = json.Unmarshal(b, &dist)
	if err != nil {
		return dist, xerrors.Errorf("Unmarshal error: %w", err)
	}
	return dist, nil
}
