package api_processing_data_formatter

import (
	api_input_reader "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Input_Reader"
	"data-platform-api-invoice-document-headers-creates-subfunc-rmq/DPFM_API_Caller/requests"
	"database/sql"
	"fmt"
)

// initializer
func (psdc *SDC) ConvertToMetaData(sdc *api_input_reader.SDC) *MetaData {
	pm := &requests.MetaData{
		BusinessPartnerID: sdc.BusinessPartnerID,
		ServiceLabel:      sdc.ServiceLabel,
	}

	data := pm
	res := MetaData{
		BusinessPartnerID: data.BusinessPartnerID,
		ServiceLabel:      data.ServiceLabel,
	}

	return &res
}

func (psdc *SDC) ConvertToProcessType() *ProcessType {
	pm := &requests.ProcessType{}
	data := pm

	processType := ProcessType{
		BulkProcess:       data.BulkProcess,
		IndividualProcess: data.IndividualProcess,
	}

	return &processType
}

func (psdc *SDC) ConvertToReferenceType() *ReferenceType {
	pm := &requests.ReferenceType{}
	data := pm

	referenceType := ReferenceType{
		OrderID:          data.OrderID,
		DeliveryDocument: data.DeliveryDocument,
	}

	return &referenceType
}

func (psdc *SDC) ConvertToOrderIDByArraySpecKey(length int) *OrderIDKey {
	pm := &requests.OrderIDKey{
		HeaderCompleteDeliveryIsDefined: getBoolPtr(true),
		HeaderDeliveryStatus:            "CL",
		HeaderBillingStatus:             "CL",
		HeaderBillingBlockStatus:        getBoolPtr(false),
	}

	for i := 0; i < length; i++ {
		pm.BillFromParty = append(pm.BillFromParty, nil)
		pm.BillToParty = append(pm.BillToParty, nil)
	}

	data := pm
	res := OrderIDKey{
		BillFromPartyFrom:               data.BillFromPartyFrom,
		BillFromPartyTo:                 data.BillFromPartyTo,
		BillToPartyFrom:                 data.BillToPartyFrom,
		BillToPartyTo:                   data.BillToPartyTo,
		BillFromParty:                   data.BillFromParty,
		BillToParty:                     data.BillToParty,
		HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
		HeaderBillingStatus:             data.HeaderBillingStatus,
		HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
	}

	return &res
}

func (psdc *SDC) ConvertToOrderIDByArraySpec(rows *sql.Rows) ([]*OrderID, error) {
	defer rows.Close()
	res := make([]*OrderID, 0)

	for i := 0; true; i++ {
		pm := &requests.OrderID{}

		if !rows.Next() {
			if i == 0 {
				return nil, fmt.Errorf("'data_platform_orders_header_data'テーブルに対象のレコードが存在しません。")
			} else {
				break
			}
		}
		err := rows.Scan(
			&pm.OrderID,
			&pm.BillFromParty,
			&pm.BillToParty,
			&pm.HeaderCompleteDeliveryIsDefined,
			&pm.HeaderDeliveryStatus,
			&pm.HeaderBillingStatus,
			&pm.HeaderBillingBlockStatus,
		)
		if err != nil {
			fmt.Printf("err = %+v \n", err)
			return nil, err
		}

		data := pm
		res = append(res, &OrderID{
			OrderID:                         data.OrderID,
			BillFromParty:                   data.BillFromParty,
			BillToParty:                     data.BillToParty,
			HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
			HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
			HeaderBillingStatus:             data.HeaderBillingStatus,
			HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
		})
	}

	return res, nil
}

func (psdc *SDC) ConvertToOrderIDByRangeSpecKey() *OrderIDKey {
	pm := &requests.OrderIDKey{
		HeaderCompleteDeliveryIsDefined: getBoolPtr(true),
		HeaderDeliveryStatus:            "CL",
		HeaderBillingStatus:             "CL",
		HeaderBillingBlockStatus:        getBoolPtr(false),
		IsCancelled:                     getBoolPtr(false),
		IsMarkedForDeletion:             getBoolPtr(false),
	}

	data := pm
	res := OrderIDKey{
		BillFromPartyFrom:               data.BillFromPartyFrom,
		BillFromPartyTo:                 data.BillFromPartyTo,
		BillToPartyFrom:                 data.BillToPartyFrom,
		BillToPartyTo:                   data.BillToPartyTo,
		BillFromParty:                   data.BillFromParty,
		BillToParty:                     data.BillToParty,
		HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
		HeaderBillingStatus:             data.HeaderBillingStatus,
		HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
		IsCancelled:                     data.IsCancelled,
		IsMarkedForDeletion:             data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToOrderIDByRangeSpec(rows *sql.Rows) ([]*OrderID, error) {
	defer rows.Close()
	res := make([]*OrderID, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.OrderID{}

		err := rows.Scan(
			&pm.OrderID,
			&pm.BillFromParty,
			&pm.BillToParty,
			&pm.HeaderCompleteDeliveryIsDefined,
			&pm.HeaderDeliveryStatus,
			&pm.HeaderBillingStatus,
			&pm.HeaderBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &OrderID{
			OrderID:                         data.OrderID,
			BillFromParty:                   data.BillFromParty,
			BillToParty:                     data.BillToParty,
			HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
			HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
			HeaderBillingStatus:             data.HeaderBillingStatus,
			HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
			IsCancelled:                     data.IsCancelled,
			IsMarkedForDeletion:             data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_orders_header_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToOrderItemInBulkProcessKey() *OrderItemKey {
	pm := &requests.OrderItemKey{
		ItemCompleteDeliveryIsDefined: true,
		ItemDeliveryStatus:            "CL",
		ItemBillingStatus:             "CL",
		ItemBillingBlockStatus:        false,
		IsCancelled:                   getBoolPtr(false),
		IsMarkedForDeletion:           getBoolPtr(false),
	}

	data := pm
	res := OrderItemKey{
		OrderID:                       data.OrderID,
		ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
		ItemDeliveryStatus:            data.ItemDeliveryStatus,
		ItemBillingStatus:             data.ItemBillingStatus,
		ItemBillingBlockStatus:        data.ItemBillingBlockStatus,
		IsCancelled:                   data.IsCancelled,
		IsMarkedForDeletion:           data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToOrderItemInBulkProcess(rows *sql.Rows) ([]*OrderItem, error) {
	defer rows.Close()
	res := make([]*OrderItem, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.OrderItem{}

		err := rows.Scan(
			&pm.OrderID,
			&pm.OrderItem,
			&pm.ItemCompleteDeliveryIsDefined,
			&pm.ItemDeliveryStatus,
			&pm.ItemBillingStatus,
			&pm.ItemBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &OrderItem{
			OrderID:                       data.OrderID,
			OrderItem:                     data.OrderItem,
			ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
			ItemDeliveryStatus:            data.ItemDeliveryStatus,
			ItemBillingStatus:             data.ItemBillingStatus,
			ItemBillingBlockStatus:        data.ItemBillingBlockStatus,
			IsCancelled:                   data.IsCancelled,
			IsMarkedForDeletion:           data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_orders_item_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToDeliveryDocumentItemInBulkProcessKey() *DeliveryDocumentItemKey {
	pm := &requests.DeliveryDocumentItemKey{
		ItemCompleteDeliveryIsDefined: true,
		// ItemDeliveryStatus:            "CL",
		ItemBillingStatus:      "CL",
		ItemBillingBlockStatus: false,
		IsCancelled:            getBoolPtr(false),
		IsMarkedForDeletion:    getBoolPtr(false),
	}

	data := pm
	res := DeliveryDocumentItemKey{
		DeliveryDocument:              data.DeliveryDocument,
		ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
		// ItemDeliveryStatus:            data.ItemDeliveryStatus,
		ItemBillingStatus:      data.ItemBillingStatus,
		ItemBillingBlockStatus: data.ItemBillingBlockStatus,
		IsCancelled:            data.IsCancelled,
		IsMarkedForDeletion:    data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToDeliveryDocumentItemInBulkProcess(rows *sql.Rows) ([]*DeliveryDocumentItem, error) {
	defer rows.Close()
	res := make([]*DeliveryDocumentItem, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.DeliveryDocumentItem{}

		err := rows.Scan(
			&pm.DeliveryDocument,
			&pm.DeliveryDocumentItem,
			&pm.ConfirmedDeliveryDate,
			&pm.ActualGoodsIssueDate,
			&pm.ItemCompleteDeliveryIsDefined,
			&pm.ItemBillingStatus,
			&pm.ItemBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &DeliveryDocumentItem{
			DeliveryDocument:              data.DeliveryDocument,
			DeliveryDocumentItem:          data.DeliveryDocumentItem,
			ConfirmedDeliveryDate:         data.ConfirmedDeliveryDate,
			ActualGoodsIssueDate:          data.ActualGoodsIssueDate,
			ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
			ItemBillingStatus:             data.ItemBillingStatus,
			ItemBillingBlockStatus:        data.ItemBillingBlockStatus,
			IsCancelled:                   data.IsCancelled,
			IsMarkedForDeletion:           data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_delivery_document_item_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToOrderIDByReferenceDocumentKey() *OrderIDKey {
	pm := &requests.OrderIDKey{
		HeaderCompleteDeliveryIsDefined: getBoolPtr(true),
		HeaderDeliveryStatus:            "CL",
		HeaderBillingBlockStatus:        getBoolPtr(false),
		IsCancelled:                     getBoolPtr(false),
		IsMarkedForDeletion:             getBoolPtr(false),
	}

	data := pm
	res := OrderIDKey{
		ReferenceDocument:               data.ReferenceDocument,
		HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
		HeaderBillingStatus:             data.HeaderBillingStatus,
		HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
		IsCancelled:                     data.IsCancelled,
		IsMarkedForDeletion:             data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToOrderIDByReferenceDocument(rows *sql.Rows) ([]*OrderID, error) {
	defer rows.Close()
	res := make([]*OrderID, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.OrderID{}

		err := rows.Scan(
			&pm.OrderID,
			&pm.HeaderCompleteDeliveryIsDefined,
			&pm.HeaderBillingStatus,
			&pm.HeaderBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &OrderID{
			OrderID:                         data.OrderID,
			HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
			HeaderBillingStatus:             data.HeaderBillingStatus,
			HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
			IsCancelled:                     data.IsCancelled,
			IsMarkedForDeletion:             data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_orders_header_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToOrderItemInIndividualProcessKey() *OrderItemKey {
	pm := &requests.OrderItemKey{
		ItemCompleteDeliveryIsDefined: true,
		ItemDeliveryStatus:            "CL",
		ItemBillingStatus:             "CL",
		ItemBillingBlockStatus:        false,
		IsCancelled:                   getBoolPtr(false),
		IsMarkedForDeletion:           getBoolPtr(false),
	}

	data := pm
	res := OrderItemKey{
		OrderID:                       data.OrderID,
		ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
		ItemDeliveryStatus:            data.ItemDeliveryStatus,
		ItemBillingStatus:             data.ItemBillingStatus,
		ItemBillingBlockStatus:        data.ItemBillingBlockStatus,
		IsCancelled:                   data.IsCancelled,
		IsMarkedForDeletion:           data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToOrderItemInIndividualProcess(rows *sql.Rows) ([]*OrderItem, error) {
	defer rows.Close()
	res := make([]*OrderItem, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.OrderItem{}

		err := rows.Scan(
			&pm.OrderID,
			&pm.OrderItem,
			&pm.ItemCompleteDeliveryIsDefined,
			&pm.ItemDeliveryStatus,
			&pm.ItemBillingStatus,
			&pm.ItemBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &OrderItem{
			OrderID:                       data.OrderID,
			OrderItem:                     data.OrderItem,
			ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
			ItemDeliveryStatus:            data.ItemDeliveryStatus,
			ItemBillingStatus:             data.ItemBillingStatus,
			ItemBillingBlockStatus:        data.ItemBillingBlockStatus,
			IsCancelled:                   data.IsCancelled,
			IsMarkedForDeletion:           data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_orders_item_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToDeliveryDocumentItemInIndividualProcessKey() *DeliveryDocumentItemKey {
	pm := &requests.DeliveryDocumentItemKey{
		ItemCompleteDeliveryIsDefined: true,
		// ItemDeliveryStatus:            "CL",
		ItemBillingStatus:      "CL",
		ItemBillingBlockStatus: false,
		IsCancelled:            getBoolPtr(false),
		IsMarkedForDeletion:    getBoolPtr(false),
	}

	data := pm
	res := DeliveryDocumentItemKey{
		DeliveryDocument:              data.DeliveryDocument,
		ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
		// ItemDeliveryStatus:            data.ItemDeliveryStatus,
		ItemBillingStatus:      data.ItemBillingStatus,
		ItemBillingBlockStatus: data.ItemBillingBlockStatus,
		IsCancelled:            data.IsCancelled,
		IsMarkedForDeletion:    data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToDeliveryDocumentItemInIndividualProcess(rows *sql.Rows) ([]*DeliveryDocumentItem, error) {
	defer rows.Close()
	res := make([]*DeliveryDocumentItem, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.DeliveryDocumentItem{}

		err := rows.Scan(
			&pm.DeliveryDocument,
			&pm.DeliveryDocumentItem,
			&pm.ConfirmedDeliveryDate,
			&pm.ActualGoodsIssueDate,
			&pm.ItemCompleteDeliveryIsDefined,
			&pm.ItemBillingStatus,
			&pm.ItemBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &DeliveryDocumentItem{
			DeliveryDocument:              data.DeliveryDocument,
			DeliveryDocumentItem:          data.DeliveryDocumentItem,
			ConfirmedDeliveryDate:         data.ConfirmedDeliveryDate,
			ActualGoodsIssueDate:          data.ActualGoodsIssueDate,
			ItemCompleteDeliveryIsDefined: data.ItemCompleteDeliveryIsDefined,
			ItemBillingStatus:             data.ItemBillingStatus,
			ItemBillingBlockStatus:        data.ItemBillingBlockStatus,
			IsCancelled:                   data.IsCancelled,
			IsMarkedForDeletion:           data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_delivery_document_item_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

// func (psdc *SDC) ConvertToOrdersHeaderPartner(
// 	sdc *api_input_reader.SDC,
// 	rows *sql.Rows,
// ) (*[]OrdersHeaderPartner, error) {
// 	var ordersHeaderPartner []OrdersHeaderPartner

// 	for i := 0; true; i++ {
// 		pm := &requests.OrdersHeaderPartner{}

// 		if !rows.Next() {
// 			if i == 0 {
// 				return nil, fmt.Errorf("'data_platform_orders_header_partner_data'テーブルに対象のレコードが存在しません。")
// 			} else {
// 				break
// 			}
// 		}
// 		err := rows.Scan(
// 			&pm.OrderID,
// 			&pm.PartnerFunction,
// 			&pm.BusinessPartner,
// 		)
// 		if err != nil {
// 			fmt.Printf("err = %+v \n", err)
// 			return nil, err
// 		}

// 		data := pm
// 		ordersHeaderPartner = append(ordersHeaderPartner, OrdersHeaderPartner{
// 			InvoiceDocument: data.InvoiceDocument,
// 			OrderID:         data.OrderID,
// 			PartnerFunction: data.PartnerFunction,
// 			BusinessPartner: data.BusinessPartner,
// 		})
// 	}

// 	return &ordersHeaderPartner, nil
// }

// func (psdc *SDC) ConvertToDeliveryDocumentByNumberSpecificationKey(sdc *api_input_reader.SDC, length int) *DeliveryDocumentKey {
// 	pm := &requests.DeliveryDocumentKey{
// 		HeaderCompleteDeliveryIsDefined: getBoolPtr(true),
// 		HeaderDeliveryStatus:            "CL",
// 		HeaderBillingStatus:             "CL",
// 		HeaderBillingBlockStatus:        getBoolPtr(false),
// 	}

// 	for i := 0; i < length; i++ {
// 		pm.BillFromParty = append(pm.BillFromParty, nil)
// 		pm.BillToParty = append(pm.BillToParty, nil)
// 	}

// 	data := pm
// 	deliveryDocumentKey := DeliveryDocumentKey{
// 		DeliveryDocument:                data.DeliveryDocument,
// 		BillFromPartyFrom:               data.BillFromPartyFrom,
// 		BillFromPartyTo:                 data.BillFromPartyTo,
// 		BillToPartyFrom:                 data.BillToPartyFrom,
// 		BillToPartyTo:                   data.BillToPartyTo,
// 		BillFromParty:                   data.BillFromParty,
// 		BillToParty:                     data.BillToParty,
// 		HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
// 		HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
// 		HeaderBillingStatus:             data.HeaderBillingStatus,
// 		HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
// 	}

// 	return &deliveryDocumentKey
// }

// func (psdc *SDC) ConvertToDeliveryDocumentByNumberSpecification(
// 	sdc *api_input_reader.SDC,
// 	rows *sql.Rows,
// ) (*[]DeliveryDocument, error) {
// 	var deliveryDocument []DeliveryDocument

// 	for i := 0; true; i++ {
// 		pm := &requests.DeliveryDocument{}

// 		if !rows.Next() {
// 			if i == 0 {
// 				return nil, fmt.Errorf("'data_platform_delivery_document_header_data'テーブルに対象のレコードが存在しません。")
// 			} else {
// 				break
// 			}
// 		}
// 		err := rows.Scan(
// 			&pm.DeliveryDocument,
// 			&pm.BillFromParty,
// 			&pm.BillToParty,
// 			&pm.HeaderCompleteDeliveryIsDefined,
// 			&pm.HeaderDeliveryStatus,
// 			&pm.HeaderBillingStatus,
// 			&pm.HeaderBillingBlockStatus,
// 		)
// 		if err != nil {
// 			fmt.Printf("err = %+v \n", err)
// 			return nil, err
// 		}

// 		data := pm
// 		deliveryDocument = append(deliveryDocument, DeliveryDocument{
// 			InvoiceDocument:                 data.InvoiceDocument,
// 			DeliveryDocument:                data.DeliveryDocument,
// 			BillFromPartyFrom:               data.BillFromPartyFrom,
// 			BillFromPartyTo:                 data.BillFromPartyTo,
// 			BillToPartyFrom:                 data.BillToPartyFrom,
// 			BillToPartyTo:                   data.BillToPartyTo,
// 			BillFromParty:                   data.BillFromParty,
// 			BillToParty:                     data.BillToParty,
// 			HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
// 			HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
// 			HeaderBillingStatus:             data.HeaderBillingStatus,
// 			HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
// 		})
// 	}

// 	return &deliveryDocument, nil
// }

func (psdc *SDC) ConvertToDeliveryDocumentByRangeSpecificationKey() *DeliveryDocumentHeaderKey {
	pm := &requests.DeliveryDocumentHeaderKey{
		HeaderCompleteDeliveryIsDefined: getBoolPtr(true),
		HeaderDeliveryStatus:            "CL",
		HeaderBillingStatus:             "CL",
		HeaderBillingBlockStatus:        getBoolPtr(false),
		IsCancelled:                     getBoolPtr(false),
		IsMarkedForDeletion:             getBoolPtr(false),
	}

	data := pm
	res := DeliveryDocumentHeaderKey{
		BillFromPartyFrom:               data.BillFromPartyFrom,
		BillFromPartyTo:                 data.BillFromPartyTo,
		BillToPartyFrom:                 data.BillToPartyFrom,
		BillToPartyTo:                   data.BillToPartyTo,
		BillFromParty:                   data.BillFromParty,
		BillToParty:                     data.BillToParty,
		HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
		HeaderBillingStatus:             data.HeaderBillingStatus,
		HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
		IsCancelled:                     data.IsCancelled,
		IsMarkedForDeletion:             data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToDeliveryDocumentByRangeSpecification(rows *sql.Rows) ([]*DeliveryDocumentHeader, error) {
	defer rows.Close()
	res := make([]*DeliveryDocumentHeader, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.DeliveryDocumentHeader{}

		err := rows.Scan(
			&pm.DeliveryDocument,
			&pm.BillFromParty,
			&pm.BillToParty,
			&pm.HeaderCompleteDeliveryIsDefined,
			&pm.HeaderDeliveryStatus,
			&pm.HeaderBillingStatus,
			&pm.HeaderBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &DeliveryDocumentHeader{
			DeliveryDocument:                data.DeliveryDocument,
			BillFromParty:                   data.BillFromParty,
			BillToParty:                     data.BillToParty,
			HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
			HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
			HeaderBillingStatus:             data.HeaderBillingStatus,
			HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
			IsCancelled:                     data.IsCancelled,
			IsMarkedForDeletion:             data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_delivery_document_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToDeliveryDocumentByReferenceDocumentKey() *DeliveryDocumentHeaderKey {
	pm := &requests.DeliveryDocumentHeaderKey{
		HeaderCompleteDeliveryIsDefined: getBoolPtr(true),
		HeaderDeliveryStatus:            "CL",
		HeaderBillingStatus:             "CL",
		HeaderBillingBlockStatus:        getBoolPtr(false),
		IsCancelled:                     getBoolPtr(false),
		IsMarkedForDeletion:             getBoolPtr(false),
	}

	data := pm
	res := DeliveryDocumentHeaderKey{
		ReferenceDocument:               data.ReferenceDocument,
		HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
		HeaderBillingStatus:             data.HeaderBillingStatus,
		HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
		IsCancelled:                     data.IsCancelled,
		IsMarkedForDeletion:             data.IsMarkedForDeletion,
	}

	return &res
}

func (psdc *SDC) ConvertToDeliveryDocumentByReferenceDocument(rows *sql.Rows) ([]*DeliveryDocumentHeader, error) {
	defer rows.Close()
	res := make([]*DeliveryDocumentHeader, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.DeliveryDocumentHeader{}

		err := rows.Scan(
			&pm.DeliveryDocument,
			&pm.HeaderCompleteDeliveryIsDefined,
			&pm.HeaderDeliveryStatus,
			&pm.HeaderBillingStatus,
			&pm.HeaderBillingBlockStatus,
			&pm.IsCancelled,
			&pm.IsMarkedForDeletion,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &DeliveryDocumentHeader{
			DeliveryDocument:                data.DeliveryDocument,
			HeaderCompleteDeliveryIsDefined: data.HeaderCompleteDeliveryIsDefined,
			HeaderDeliveryStatus:            data.HeaderDeliveryStatus,
			HeaderBillingStatus:             data.HeaderBillingStatus,
			HeaderBillingBlockStatus:        data.HeaderBillingBlockStatus,
			IsCancelled:                     data.IsCancelled,
			IsMarkedForDeletion:             data.IsMarkedForDeletion,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_delivery_document_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

// Header
func (psdc *SDC) ConvertToOrdersHeader(rows *sql.Rows) ([]*OrdersHeader, error) {
	defer rows.Close()
	res := make([]*OrdersHeader, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.OrdersHeader{}

		err := rows.Scan(
			&pm.OrderID,
			&pm.OrderType,
			&pm.SupplyChainRelationshipID,
			&pm.SupplyChainRelationshipBillingID,
			&pm.SupplyChainRelationshipPaymentID,
			&pm.Buyer,
			&pm.Seller,
			&pm.BillToParty,
			&pm.BillFromParty,
			&pm.BillToCountry,
			&pm.BillFromCountry,
			&pm.Payer,
			&pm.Payee,
			&pm.ContractType,
			&pm.OrderValidityStartDate,
			&pm.OrderValidityEndDate,
			&pm.InvoicePeriodStartDate,
			&pm.InvoicePeriodEndDate,
			&pm.TotalNetAmount,
			&pm.TotalTaxAmount,
			&pm.TotalGrossAmount,
			&pm.TransactionCurrency,
			&pm.PricingDate,
			&pm.Incoterms,
			&pm.PaymentTerms,
			&pm.PaymentMethod,
			&pm.IsExportImport,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &OrdersHeader{
			OrderID:                          data.OrderID,
			OrderType:                        data.OrderType,
			SupplyChainRelationshipID:        data.SupplyChainRelationshipID,
			SupplyChainRelationshipBillingID: data.SupplyChainRelationshipBillingID,
			SupplyChainRelationshipPaymentID: data.SupplyChainRelationshipPaymentID,
			Buyer:                            data.Buyer,
			Seller:                           data.Seller,
			BillToParty:                      data.BillToParty,
			BillFromParty:                    data.BillFromParty,
			BillToCountry:                    data.BillToCountry,
			BillFromCountry:                  data.BillFromCountry,
			Payer:                            data.Payer,
			Payee:                            data.Payee,
			ContractType:                     data.ContractType,
			OrderValidityStartDate:           data.OrderValidityStartDate,
			OrderValidityEndDate:             data.OrderValidityEndDate,
			InvoicePeriodStartDate:           data.InvoicePeriodStartDate,
			InvoicePeriodEndDate:             data.InvoicePeriodEndDate,
			TotalNetAmount:                   data.TotalNetAmount,
			TotalTaxAmount:                   data.TotalTaxAmount,
			TotalGrossAmount:                 data.TotalGrossAmount,
			TransactionCurrency:              data.TransactionCurrency,
			PricingDate:                      data.PricingDate,
			Incoterms:                        data.Incoterms,
			PaymentTerms:                     data.PaymentTerms,
			PaymentMethod:                    data.PaymentMethod,
			IsExportImport:                   data.IsExportImport,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_orders_header_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToDeliveryDocumentHeaderData(
	sdc *api_input_reader.SDC,
	rows *sql.Rows,
) ([]*DeliveryDocumentHeaderData, error) {

	defer rows.Close()
	res := make([]*DeliveryDocumentHeaderData, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.DeliveryDocumentHeaderData{}

		err := rows.Scan(
			&pm.DeliveryDocument,
			&pm.SupplyChainRelationshipID,
			&pm.SupplyChainRelationshipDeliveryID,
			&pm.SupplyChainRelationshipDeliveryPlantID,
			&pm.SupplyChainRelationshipBillingID,
			&pm.SupplyChainRelationshipPaymentID,
			&pm.Buyer,
			&pm.Seller,
			&pm.DeliverToParty,
			&pm.DeliverFromParty,
			&pm.DeliverToPlant,
			&pm.DeliverFromPlant,
			&pm.BillToParty,
			&pm.BillFromParty,
			&pm.BillToCountry,
			&pm.BillFromCountry,
			&pm.Payer,
			&pm.Payee,
			&pm.IsExportImport,
			&pm.OrderID,
			&pm.OrderItem,
			&pm.ContractType,
			&pm.OrderValidityStartDate,
			&pm.OrderValidityEndDate,
			&pm.GoodsIssueOrReceiptSlipNumber,
			&pm.Incoterms,
			&pm.TransactionCurrency,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &DeliveryDocumentHeaderData{
			DeliveryDocument:                       data.DeliveryDocument,
			SupplyChainRelationshipID:              data.SupplyChainRelationshipID,
			SupplyChainRelationshipDeliveryID:      data.SupplyChainRelationshipDeliveryID,
			SupplyChainRelationshipDeliveryPlantID: data.SupplyChainRelationshipDeliveryPlantID,
			SupplyChainRelationshipBillingID:       data.SupplyChainRelationshipBillingID,
			SupplyChainRelationshipPaymentID:       data.SupplyChainRelationshipPaymentID,
			Buyer:                                  data.Buyer,
			Seller:                                 data.Seller,
			DeliverToParty:                         data.DeliverToParty,
			DeliverFromParty:                       data.DeliverFromParty,
			DeliverToPlant:                         data.DeliverToPlant,
			DeliverFromPlant:                       data.DeliverFromPlant,
			BillToParty:                            data.BillToParty,
			BillFromParty:                          data.BillFromParty,
			BillToCountry:                          data.BillToCountry,
			BillFromCountry:                        data.BillFromCountry,
			Payer:                                  data.Payer,
			Payee:                                  data.Payee,
			IsExportImport:                         data.IsExportImport,
			OrderID:                                data.OrderID,
			OrderItem:                              data.OrderItem,
			ContractType:                           data.ContractType,
			OrderValidityStartDate:                 data.OrderValidityStartDate,
			OrderValidityEndDate:                   data.OrderValidityEndDate,
			GoodsIssueOrReceiptSlipNumber:          data.GoodsIssueOrReceiptSlipNumber,
			Incoterms:                              data.Incoterms,
			TransactionCurrency:                    data.TransactionCurrency,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_delivery_document_header_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToInvoiceDocumentHeaderKey() *CalculateInvoiceDocumentKey {
	pm := &requests.CalculateInvoiceDocumentKey{
		FieldNameWithNumberRange: "InvoiceDocument",
	}

	data := pm
	res := CalculateInvoiceDocumentKey{
		ServiceLabel:             data.ServiceLabel,
		FieldNameWithNumberRange: data.FieldNameWithNumberRange,
	}

	return &res
}

func (psdc *SDC) ConvertToInvoiceDocumentHeaderQueryGets(
	sdc *api_input_reader.SDC,
	rows *sql.Rows,
) (*CalculateInvoiceDocumentQueryGets, error) {
	defer rows.Close()
	var res *CalculateInvoiceDocumentQueryGets

	i := 0
	for rows.Next() {
		i++
		pm := &requests.CalculateInvoiceDocumentQueryGets{}

		err := rows.Scan(
			&pm.ServiceLabel,
			&pm.FieldNameWithNumberRange,
			&pm.InvoiceDocumentLatestNumber,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = &CalculateInvoiceDocumentQueryGets{
			ServiceLabel:                data.ServiceLabel,
			FieldNameWithNumberRange:    data.FieldNameWithNumberRange,
			InvoiceDocumentLatestNumber: data.InvoiceDocumentLatestNumber,
		}

	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_number_range_latest_number_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToCalculateInvoiceDocument(invoiceDocumentLatestNumber *int) *CalculateInvoiceDocument {
	pm := &requests.CalculateInvoiceDocument{}

	pm.InvoiceDocumentLatestNumber = invoiceDocumentLatestNumber

	data := pm
	res := CalculateInvoiceDocument{
		InvoiceDocumentLatestNumber: data.InvoiceDocumentLatestNumber,
		InvoiceDocument:             data.InvoiceDocument,
	}

	return &res
}

// func (psdc *SDC) ConvertToTotalNetAmountQueryGets(
// 	sdc *api_input_reader.SDC,
// 	rows *sql.Rows,
// ) (*TotalNetAmountQueryGets, error) {
// 	pm := &requests.TotalNetAmountQueryGets{}

// 	for i := 0; true; i++ {
// 		if !rows.Next() {
// 			if i == 0 {
// 				return nil, fmt.Errorf("'data_platform_invoice_document_header_data'テーブルに対象のレコードが存在しません。")
// 			} else {
// 				break
// 			}
// 		}
// 		err := rows.Scan(
// 			&pm.InvoiceDocument,
// 			&pm.TotalNetAmount,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	data := pm
// 	totalNetAmountQueryGets := TotalNetAmountQueryGets{
// 		InvoiceDocument: data.InvoiceDocument,
// 		TotalNetAmount:  data.TotalNetAmount,
// 	}

// 	return &totalNetAmountQueryGets, nil
// }

// func (psdc *SDC) ConvertToTotalNetAmount(
// 	inputTotalNetAmount *float32,
// ) *TotalNetAmount {
// 	pm := &requests.TotalNetAmount{}

// 	pm.TotalNetAmount = inputTotalNetAmount

// 	data := pm
// 	totalNetAmount := TotalNetAmount{
// 		InvoiceDocument: data.InvoiceDocument,
// 		TotalNetAmount:  data.TotalNetAmount,
// 	}

// 	return &totalNetAmount
// }

// // HeaderPartner
// func (psdc *SDC) ConvertToHeaderOrdersHeaderPartner(
// 	sdc *api_input_reader.SDC,
// 	rows *sql.Rows,
// ) (*[]HeaderOrdersHeaderPartner, error) {
// 	var headerOrdersHeaderPartner []HeaderOrdersHeaderPartner
// 	pm := &requests.HeaderOrdersHeaderPartner{}

// 	for i := 0; true; i++ {
// 		if !rows.Next() {
// 			if i == 0 {
// 				return nil, fmt.Errorf("'data_platform_orders_header_partner_data'テーブルに対象のレコードが存在しません。")
// 			} else {
// 				break
// 			}
// 		}
// 		err := rows.Scan(
// 			&pm.OrderID,
// 			&pm.PartnerFunction,
// 			&pm.BusinessPartner,
// 			&pm.BusinessPartnerFullName,
// 			&pm.BusinessPartnerName,
// 			&pm.Organization,
// 			&pm.Country,
// 			&pm.Language,
// 			&pm.Currency,
// 			&pm.ExternalDocumentID,
// 			&pm.AddressID,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		data := pm
// 		headerOrdersHeaderPartner = append(headerOrdersHeaderPartner, HeaderOrdersHeaderPartner{
// 			InvoiceDocument:         data.InvoiceDocument,
// 			OrderID:                 data.OrderID,
// 			PartnerFunction:         data.PartnerFunction,
// 			BusinessPartner:         data.BusinessPartner,
// 			BusinessPartnerFullName: data.BusinessPartnerFullName,
// 			BusinessPartnerName:     data.BusinessPartnerName,
// 			Organization:            data.Organization,
// 			Country:                 data.Country,
// 			Language:                data.Language,
// 			Currency:                data.Currency,
// 			ExternalDocumentID:      data.ExternalDocumentID,
// 			AddressID:               data.AddressID,
// 		})
// 	}

// 	return &headerOrdersHeaderPartner, nil
// }

// func (psdc *SDC) ConvertToHeaderDeliveryDocumentHeaderPartner(
// 	sdc *api_input_reader.SDC,
// 	rows *sql.Rows,
// ) (*[]HeaderDeliveryDocumentHeaderPartner, error) {
// 	var headerDeliveryDocumentHeaderPartner []HeaderDeliveryDocumentHeaderPartner
// 	pm := &requests.HeaderDeliveryDocumentHeaderPartner{}

// 	for i := 0; true; i++ {
// 		if !rows.Next() {
// 			if i == 0 {
// 				return nil, fmt.Errorf("`data_platform_delivery_document_header_partner_data`テーブルに対象のレコードが存在しません。")
// 			} else {
// 				break
// 			}
// 		}
// 		err := rows.Scan(
// 			&pm.DeliveryDocument,
// 			&pm.PartnerFunction,
// 			&pm.BusinessPartner,
// 			&pm.BusinessPartnerFullName,
// 			&pm.BusinessPartnerName,
// 			&pm.Organization,
// 			&pm.Country,
// 			&pm.Language,
// 			&pm.Currency,
// 			&pm.ExternalDocumentID,
// 			&pm.AddressID,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		data := pm
// 		headerDeliveryDocumentHeaderPartner = append(headerDeliveryDocumentHeaderPartner, HeaderDeliveryDocumentHeaderPartner{
// 			InvoiceDocument:         data.InvoiceDocument,
// 			DeliveryDocument:        data.DeliveryDocument,
// 			PartnerFunction:         data.PartnerFunction,
// 			BusinessPartner:         data.BusinessPartner,
// 			BusinessPartnerFullName: data.BusinessPartnerFullName,
// 			BusinessPartnerName:     data.BusinessPartnerName,
// 			Organization:            data.Organization,
// 			Country:                 data.Country,
// 			Language:                data.Language,
// 			Currency:                data.Currency,
// 			ExternalDocumentID:      data.ExternalDocumentID,
// 			AddressID:               data.AddressID,
// 		})
// 	}

// 	return &headerDeliveryDocumentHeaderPartner, nil
// }

// Item
func (psdc *SDC) ConvertToDeliveryDocumentItemData(
	sdc *api_input_reader.SDC,
	rows *sql.Rows,
) ([]*DeliveryDocumentItemData, error) {

	defer rows.Close()
	res := make([]*DeliveryDocumentItemData, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.DeliveryDocumentItemData{}

		err := rows.Scan(
			&pm.DeliveryDocument,
			&pm.DeliveryDocumentItem,
			&pm.DeliveryDocumentItemCategory,
			&pm.SupplyChainRelationshipID,
			&pm.SupplyChainRelationshipDeliveryID,
			&pm.SupplyChainRelationshipDeliveryPlantID,
			&pm.SupplyChainRelationshipBillingID,
			&pm.SupplyChainRelationshipPaymentID,
			&pm.Buyer,
			&pm.Seller,
			&pm.DeliverToParty,
			&pm.DeliverFromParty,
			&pm.DeliverToPlant,
			&pm.DeliverFromPlant,
			&pm.BillToParty,
			&pm.BillFromParty,
			&pm.BillToCountry,
			&pm.BillFromCountry,
			&pm.Payer,
			&pm.Payee,
			&pm.DeliverToPlantStorageLocation,
			&pm.DeliverFromPlantStorageLocation,
			&pm.ProductionPlantBusinessPartner,
			&pm.ProductionPlant,
			&pm.ProductionPlantStorageLocation,
			&pm.DeliveryDocumentItemText,
			&pm.DeliveryDocumentItemTextByBuyer,
			&pm.DeliveryDocumentItemTextBySeller,
			&pm.Product,
			&pm.ProductStandardID,
			&pm.ProductGroup,
			&pm.BaseUnit,
			&pm.ActualGoodsIssueDate,
			&pm.ActualGoodsIssueTime,
			&pm.ActualGoodsReceiptDate,
			&pm.ActualGoodsReceiptTime,
			&pm.ActualGoodsIssueQuantity,
			&pm.ActualGoodsIssueQtyInBaseUnit,
			&pm.ActualGoodsReceiptQuantity,
			&pm.ActualGoodsReceiptQtyInBaseUnit,
			&pm.ItemGrossWeight,
			&pm.ItemNetWeight,
			&pm.ItemWeightUnit,
			&pm.NetAmount,
			&pm.TaxAmount,
			&pm.GrossAmount,
			&pm.OrderID,
			&pm.OrderItem,
			&pm.OrderType,
			&pm.ContractType,
			&pm.OrderValidityStartDate,
			&pm.OrderValidityEndDate,
			&pm.PaymentTerms,
			&pm.PaymentMethod,
			&pm.InvoicePeriodStartDate,
			&pm.InvoicePeriodEndDate,
			&pm.Project,
			&pm.ReferenceDocument,
			&pm.ReferenceDocumentItem,
			&pm.TransactionTaxClassification,
			&pm.ProductTaxClassificationBillToCountry,
			&pm.ProductTaxClassificationBillFromCountry,
			&pm.DefinedTaxClassifications,
			&pm.TaxCode,
			&pm.TaxRate,
			&pm.CountryOfOrigin,
			&pm.CountryOfOriginLanguage,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &DeliveryDocumentItemData{
			DeliveryDocument:                        data.DeliveryDocument,
			DeliveryDocumentItem:                    data.DeliveryDocumentItem,
			DeliveryDocumentItemCategory:            data.DeliveryDocumentItemCategory,
			SupplyChainRelationshipID:               data.SupplyChainRelationshipID,
			SupplyChainRelationshipDeliveryID:       data.SupplyChainRelationshipDeliveryID,
			SupplyChainRelationshipDeliveryPlantID:  data.SupplyChainRelationshipDeliveryPlantID,
			SupplyChainRelationshipBillingID:        data.SupplyChainRelationshipBillingID,
			SupplyChainRelationshipPaymentID:        data.SupplyChainRelationshipPaymentID,
			Buyer:                                   data.Buyer,
			Seller:                                  data.Seller,
			DeliverToParty:                          data.DeliverToParty,
			DeliverFromParty:                        data.DeliverFromParty,
			DeliverToPlant:                          data.DeliverToPlant,
			DeliverFromPlant:                        data.DeliverFromPlant,
			BillToParty:                             data.BillToParty,
			BillFromParty:                           data.BillFromParty,
			BillToCountry:                           data.BillToCountry,
			BillFromCountry:                         data.BillFromCountry,
			Payer:                                   data.Payer,
			Payee:                                   data.Payee,
			DeliverToPlantStorageLocation:           data.DeliverToPlantStorageLocation,
			DeliverFromPlantStorageLocation:         data.DeliverFromPlantStorageLocation,
			ProductionPlantBusinessPartner:          data.ProductionPlantBusinessPartner,
			ProductionPlant:                         data.ProductionPlant,
			ProductionPlantStorageLocation:          data.ProductionPlantStorageLocation,
			DeliveryDocumentItemText:                data.DeliveryDocumentItemText,
			DeliveryDocumentItemTextByBuyer:         data.DeliveryDocumentItemTextByBuyer,
			DeliveryDocumentItemTextBySeller:        data.DeliveryDocumentItemTextBySeller,
			Product:                                 data.Product,
			ProductStandardID:                       data.ProductStandardID,
			ProductGroup:                            data.ProductGroup,
			BaseUnit:                                data.BaseUnit,
			ActualGoodsIssueDate:                    data.ActualGoodsIssueDate,
			ActualGoodsIssueTime:                    data.ActualGoodsIssueTime,
			ActualGoodsReceiptDate:                  data.ActualGoodsReceiptDate,
			ActualGoodsReceiptTime:                  data.ActualGoodsReceiptTime,
			ActualGoodsIssueQuantity:                data.ActualGoodsIssueQuantity,
			ActualGoodsIssueQtyInBaseUnit:           data.ActualGoodsIssueQtyInBaseUnit,
			ActualGoodsReceiptQuantity:              data.ActualGoodsReceiptQuantity,
			ActualGoodsReceiptQtyInBaseUnit:         data.ActualGoodsReceiptQtyInBaseUnit,
			ItemGrossWeight:                         data.ItemGrossWeight,
			ItemNetWeight:                           data.ItemNetWeight,
			ItemWeightUnit:                          data.ItemWeightUnit,
			NetAmount:                               data.NetAmount,
			TaxAmount:                               data.TaxAmount,
			GrossAmount:                             data.GrossAmount,
			OrderID:                                 data.OrderID,
			OrderItem:                               data.OrderItem,
			OrderType:                               data.OrderType,
			ContractType:                            data.ContractType,
			OrderValidityStartDate:                  data.OrderValidityStartDate,
			OrderValidityEndDate:                    data.OrderValidityEndDate,
			PaymentTerms:                            data.PaymentTerms,
			PaymentMethod:                           data.PaymentMethod,
			InvoicePeriodStartDate:                  data.InvoicePeriodStartDate,
			InvoicePeriodEndDate:                    data.InvoicePeriodEndDate,
			Project:                                 data.Project,
			ReferenceDocument:                       data.ReferenceDocument,
			ReferenceDocumentItem:                   data.ReferenceDocumentItem,
			TransactionTaxClassification:            data.TransactionTaxClassification,
			ProductTaxClassificationBillToCountry:   data.ProductTaxClassificationBillToCountry,
			ProductTaxClassificationBillFromCountry: data.ProductTaxClassificationBillFromCountry,
			DefinedTaxClassifications:               data.DefinedTaxClassifications,
			TaxCode:                                 data.TaxCode,
			TaxRate:                                 data.TaxRate,
			CountryOfOrigin:                         data.CountryOfOrigin,
			CountryOfOriginLanguage:                 data.CountryOfOriginLanguage,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_delivery_document_item_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func (psdc *SDC) ConvertToOrdersItem(
	sdc *api_input_reader.SDC,
	rows *sql.Rows,
) ([]*OrdersItem, error) {

	defer rows.Close()
	res := make([]*OrdersItem, 0)

	i := 0
	for rows.Next() {
		i++
		pm := &requests.OrdersItem{}

		err := rows.Scan(
			&pm.OrderID,
			&pm.OrderItem,
			&pm.OrderItemCategory,
			&pm.SupplyChainRelationshipID,
			&pm.OrderItemText,
			&pm.OrderItemTextByBuyer,
			&pm.OrderItemTextBySeller,
			&pm.Product,
			&pm.ProductStandardID,
			&pm.ProductGroup,
			&pm.BaseUnit,
			&pm.PricingDate,
			&pm.OrderQuantityInBaseUnit,
			&pm.OrderQuantityInDeliveryUnit,
			&pm.NetAmount,
			&pm.TaxAmount,
			&pm.GrossAmount,
			&pm.Incoterms,
			&pm.TransactionTaxClassification,
			&pm.ProductTaxClassificationBillToCountry,
			&pm.ProductTaxClassificationBillFromCountry,
			&pm.DefinedTaxClassification,
			&pm.PaymentTerms,
			&pm.PaymentMethod,
			&pm.Project,
			&pm.ReferenceDocument,
			&pm.ReferenceDocumentItem,
			&pm.TaxCode,
			&pm.TaxRate,
		)
		if err != nil {
			return nil, err
		}

		data := pm
		res = append(res, &OrdersItem{
			OrderID:                                 data.OrderID,
			OrderItem:                               data.OrderItem,
			OrderItemCategory:                       data.OrderItemCategory,
			SupplyChainRelationshipID:               data.SupplyChainRelationshipID,
			OrderItemText:                           data.OrderItemText,
			OrderItemTextByBuyer:                    data.OrderItemTextByBuyer,
			OrderItemTextBySeller:                   data.OrderItemTextBySeller,
			Product:                                 data.Product,
			ProductStandardID:                       data.ProductStandardID,
			ProductGroup:                            data.ProductGroup,
			BaseUnit:                                data.BaseUnit,
			PricingDate:                             data.PricingDate,
			OrderQuantityInBaseUnit:                 data.OrderQuantityInBaseUnit,
			OrderQuantityInDeliveryUnit:             data.OrderQuantityInDeliveryUnit,
			NetAmount:                               data.NetAmount,
			TaxAmount:                               data.TaxAmount,
			GrossAmount:                             data.GrossAmount,
			Incoterms:                               data.Incoterms,
			TransactionTaxClassification:            data.TransactionTaxClassification,
			ProductTaxClassificationBillToCountry:   data.ProductTaxClassificationBillToCountry,
			ProductTaxClassificationBillFromCountry: data.ProductTaxClassificationBillFromCountry,
			DefinedTaxClassification:                data.DefinedTaxClassification,
			PaymentTerms:                            data.PaymentTerms,
			PaymentMethod:                           data.PaymentMethod,
			Project:                                 data.Project,
			ReferenceDocument:                       data.ReferenceDocument,
			ReferenceDocumentItem:                   data.ReferenceDocumentItem,
			TaxCode:                                 data.TaxCode,
			TaxRate:                                 data.TaxRate,
		})
	}
	if i == 0 {
		return nil, fmt.Errorf("'data_platform_orders_item_data'テーブルに対象のレコードが存在しません。")
	}

	return res, nil
}

func getBoolPtr(b bool) *bool {
	return &b
}
