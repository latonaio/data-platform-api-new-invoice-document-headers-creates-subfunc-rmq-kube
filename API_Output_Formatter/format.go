package dpfm_api_output_formatter

import (
	api_input_reader "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Input_Reader"
	api_processing_data_formatter "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Processing_Data_Formatter"
	"encoding/json"

	"golang.org/x/xerrors"
)

func ConvertToHeader(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) (*Header, error) {
	var err error

	header := &Header{}
	inputHeader := sdc.InvoiceDocument

	referenceType := psdc.ReferenceType

	// 入力ファイル
	header, err = jsonTypeConversion(header, inputHeader)
	if err != nil {
		return nil, xerrors.Errorf("request create error: %w", err)
	}

	// 1-1

	if referenceType.OrderID {
		header, err = jsonTypeConversion(header, psdc.OrdersHeader[0])
		if err != nil {
			return nil, xerrors.Errorf("request create error: %w", err)
		}
	}

	// header.DeliveryDocument = psdc.CalculateDeliveryDocument.DeliveryDocument
	// header.DocumentDate = psdc.DocumentDate.DocumentDate
	// header.InvoiceDocumentDate = psdc.InvoiceDocumentDate.InvoiceDocumentDate
	// header.HeaderCompleteDeliveryIsDefined = getBoolPtr(false)
	// header.HeaderDeliveryStatus = getStringPtr("NP")
	// header.CreationDate = psdc.CreationDateHeader.CreationDate
	// header.CreationTime = psdc.CreationTimeHeader.CreationTime
	// header.LastChangeDate = psdc.LastChangeDateHeader.LastChangeDate
	// header.LastChangeTime = psdc.LastChangeTimeHeader.LastChangeTime
	// header.HeaderBillingStatus = getStringPtr("NP")
	// header.HeaderBillingConfStatus = getStringPtr("NP")
	// header.HeaderGrossWeight = psdc.HeaderGrossWeight.HeaderGrossWeight
	// header.HeaderNetWeight = psdc.HeaderNetWeight.HeaderNetWeight
	// header.HeaderDeliveryBlockStatus = getBoolPtr(false)
	// header.HeaderIssuingBlockStatus = getBoolPtr(false)
	// header.HeaderReceivingBlockStatus = getBoolPtr(false)
	// header.HeaderBillingBlockStatus = getBoolPtr(false)

	return header, nil
}

func getBoolPtr(b bool) *bool {
	return &b
}

func getStringPtr(s string) *string {
	return &s
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

// func ConvertToHeaderPartner(
// 	sdc *api_input_reader.SDC,
// 	psdc *api_processing_data_formatter.SDC,
// ) (*[]HeaderPartner, error) {
// 	calculateInvoiceDocument := psdc.CalculateInvoiceDocument
// 	headerOrdersHeaderPartner := psdc.HeaderOrdersHeaderPartner
// 	headerPartners := make([]HeaderPartner, 0, len(*headerOrdersHeaderPartner))

// 	for _, v := range *headerOrdersHeaderPartner {
// 		headerPartner := HeaderPartner{}

// 		data, err := json.Marshal(v)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = json.Unmarshal(data, &headerPartner)
// 		if err != nil {
// 			return nil, err
// 		}

// 		headerPartner.InvoiceDocument = calculateInvoiceDocument.InvoiceDocumentLatestNumber
// 		headerPartners = append(headerPartners, headerPartner)
// 	}

// 	return &headerPartners, nil
// }
