package subfunction

import (
	"context"
	api_input_reader "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Output_Formatter"
	api_processing_data_formatter "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Processing_Data_Formatter"
	"strings"

	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	"golang.org/x/xerrors"
)

type SubFunction struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewSubFunction(ctx context.Context, db *database.Mysql, l *logger.Logger) *SubFunction {
	return &SubFunction{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (f *SubFunction) MetaData(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.MetaData {
	metaData := new(api_processing_data_formatter.MetaData)

	metaData = psdc.ConvertToMetaData(sdc)

	return metaData
}

func (f *SubFunction) ProcessType(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.ProcessType {
	processType := psdc.ConvertToProcessType()

	processType.BulkProcess = true
	// processType.IndividualProcess = true

	return processType
}

func (f *SubFunction) ReferenceType(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.ReferenceType {
	referenceType := psdc.ConvertToReferenceType()

	referenceType.OrderID = true
	// referenceType.DeliveryDocument = true

	return referenceType
}

func (f *SubFunction) OrderIDByArraySpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderID, error) {
	args := make([]interface{}, 0)

	billFromParty := sdc.InputParameters.BillFromParty
	billToParty := sdc.InputParameters.BillToParty

	if len(*billFromParty) != len(*billToParty) {
		return nil, nil
	}

	dataKey := psdc.ConvertToOrderIDByArraySpecKey(len(*billFromParty))

	for i := range *billFromParty {
		dataKey.BillFromParty[i] = (*billFromParty)[i]
		dataKey.BillToParty[i] = (*billToParty)[i]
	}

	repeat := strings.Repeat("(?,?),", len(dataKey.BillFromParty)-1) + "(?,?)"
	for i := range dataKey.BillFromParty {
		args = append(args, dataKey.BillFromParty[i], dataKey.BillToParty[i])
	}

	args = append(
		args,
		dataKey.HeaderCompleteDeliveryIsDefined,
		dataKey.HeaderDeliveryStatus,
		dataKey.HeaderBillingBlockStatus,
		dataKey.HeaderBillingStatus,
	)

	var count *int
	err := f.db.QueryRow(
		`SELECT COUNT(*)
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE (BillFromParty, BillToParty) IN ( `+repeat+` )
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus) = (?, ?, ?)
		AND HeaderBillingStatus <> ?;`, args...,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if *count == 0 || *count > 1000 {
		return nil, xerrors.Errorf("OrderIDの検索結果がゼロ件または1,000件超です。")
	}

	rows, err := f.db.Query(
		`SELECT OrderID, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingStatus, HeaderBillingBlockStatus
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE (BillFromParty, BillToParty) IN ( `+repeat+` )
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus) = (?, ?, ?)
		AND HeaderBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToOrderIDByArraySpec(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrderIDByRangeSpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderID, error) {
	dataKey := psdc.ConvertToOrderIDByRangeSpecKey()

	dataKey.BillFromPartyFrom = sdc.InputParameters.BillFromPartyFrom
	dataKey.BillFromPartyTo = sdc.InputParameters.BillFromPartyTo
	dataKey.BillToPartyFrom = sdc.InputParameters.BillToPartyFrom
	dataKey.BillToPartyTo = sdc.InputParameters.BillToPartyTo

	count := new(int)
	err := f.db.QueryRow(
		`SELECT COUNT(*)
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE BillFromParty BETWEEN ? AND ?
		AND BillToParty BETWEEN ? AND ?
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, dataKey.BillFromPartyFrom, dataKey.BillFromPartyTo, dataKey.BillToPartyFrom, dataKey.BillToPartyTo, dataKey.HeaderCompleteDeliveryIsDefined, dataKey.HeaderDeliveryStatus, dataKey.HeaderBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.HeaderBillingStatus,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if *count == 0 || *count > 1000 {
		return nil, xerrors.Errorf("OrderIDの検索結果がゼロ件または1,000件超です。")
	}

	rows, err := f.db.Query(
		`SELECT OrderID, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingStatus,IsCancelled, IsMarkedForDeletion, HeaderBillingBlockStatus
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE BillFromParty BETWEEN ? AND ?
		AND BillToParty BETWEEN ? AND ?
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, dataKey.BillFromPartyFrom, dataKey.BillFromPartyTo, dataKey.BillToPartyFrom, dataKey.BillToPartyTo, dataKey.HeaderCompleteDeliveryIsDefined, dataKey.HeaderDeliveryStatus, dataKey.HeaderBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.HeaderBillingStatus,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToOrderIDByRangeSpec(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrderItemInBulkProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderItem, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToOrderItemInBulkProcessKey()

	orderID := psdc.OrderID

	for i := range orderID {
		dataKey.OrderID = append(dataKey.OrderID, (orderID)[i].OrderID)
	}

	repeat := strings.Repeat("?,", len(dataKey.OrderID)-1) + "?"
	for _, v := range dataKey.OrderID {
		args = append(args, v)
	}

	args = append(args, dataKey.ItemCompleteDeliveryIsDefined, dataKey.ItemDeliveryStatus, dataKey.ItemBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.ItemBillingStatus)

	rows, err := f.db.Query(
		`SELECT OrderID, OrderItem, ItemCompleteDeliveryIsDefined, ItemDeliveryStatus, ItemBillingStatus, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_item_data
		WHERE OrderID IN ( `+repeat+` )
		AND (ItemCompleteDeliveryIsDefined, ItemDeliveryStatus, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND ItemBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToOrderItemInBulkProcess(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentItemInBulkProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentItem, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToDeliveryDocumentItemInBulkProcessKey()

	dataKey.ConfirmedDeliveryDateFrom = sdc.InputParameters.ConfirmedDeliveryDateFrom
	dataKey.ConfirmedDeliveryDateTo = sdc.InputParameters.ConfirmedDeliveryDateTo
	dataKey.ActualGoodsIssueDateFrom = sdc.InputParameters.ActualGoodsIssueDateFrom
	dataKey.ActualGoodsIssueDateTo = sdc.InputParameters.ActualGoodsIssueDateTo

	deliveryDocumentItem := psdc.DeliveryDocumentHeader

	for i := range deliveryDocumentItem {
		dataKey.DeliveryDocument = append(dataKey.DeliveryDocument, *(deliveryDocumentItem)[i].DeliveryDocument)
	}

	repeat := strings.Repeat("?,", len(dataKey.DeliveryDocument)-1) + "?"
	for _, v := range dataKey.DeliveryDocument {
		args = append(args, v)
	}

	args = append(args, dataKey.ConfirmedDeliveryDateFrom, dataKey.ConfirmedDeliveryDateTo, dataKey.ActualGoodsIssueDateFrom, dataKey.ActualGoodsIssueDateTo, dataKey.ItemCompleteDeliveryIsDefined, dataKey.ItemBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.ItemBillingStatus)

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, DeliveryDocumentItem, ConfirmedDeliveryDate, ActualGoodsIssueDate, ItemCompleteDeliveryIsDefined, ItemBillingStatus, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_item_data
		WHERE DeliveryDocument IN ( `+repeat+` )
		AND ConfirmedDeliveryDate BETWEEN ? AND ?
		AND ActualGoodsIssueDate BETWEEN ? AND ?
		AND (ItemCompleteDeliveryIsDefined, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?)
		AND ItemBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocumentItemInBulkProcess(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrderIDByReferenceDocument(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderID, error) {
	dataKey := psdc.ConvertToOrderIDByReferenceDocumentKey()

	dataKey.ReferenceDocument = sdc.InputParameters.ReferenceDocument

	rows, err := f.db.Query(
		`SELECT OrderID, HeaderCompleteDeliveryIsDefined, HeaderBillingStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE (OrderID, HeaderCompleteDeliveryIsDefined, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, dataKey.ReferenceDocument, dataKey.HeaderCompleteDeliveryIsDefined, dataKey.HeaderBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.HeaderBillingStatus,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToOrderIDByReferenceDocument(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrderItemInIndividualProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderItem, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToOrderItemInIndividualProcessKey()

	orderID := psdc.OrderID

	for i := range orderID {
		dataKey.OrderID = append(dataKey.OrderID, (orderID)[i].OrderID)
	}

	repeat := strings.Repeat("?,", len(dataKey.OrderID)-1) + "?"
	for _, v := range dataKey.OrderID {
		args = append(args, v)
	}

	args = append(args, dataKey.ItemCompleteDeliveryIsDefined, dataKey.ItemDeliveryStatus, dataKey.ItemBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.ItemBillingStatus)

	rows, err := f.db.Query(
		`SELECT OrderID, OrderItem, ItemCompleteDeliveryIsDefined, ItemDeliveryStatus, ItemBillingStatus, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_item_data
		WHERE OrderID IN ( `+repeat+` )
		AND (ItemCompleteDeliveryIsDefined, ItemDeliveryStatus, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND ItemBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToOrderItemInIndividualProcess(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentItemInIndividualProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentItem, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToDeliveryDocumentItemInIndividualProcessKey()

	dataKey.ConfirmedDeliveryDateFrom = sdc.InputParameters.ConfirmedDeliveryDateFrom
	dataKey.ConfirmedDeliveryDateTo = sdc.InputParameters.ConfirmedDeliveryDateTo
	dataKey.ActualGoodsIssueDateFrom = sdc.InputParameters.ActualGoodsIssueDateFrom
	dataKey.ActualGoodsIssueDateTo = sdc.InputParameters.ActualGoodsIssueDateTo

	deliveryDocumentItem := psdc.DeliveryDocumentHeader

	for i := range deliveryDocumentItem {
		dataKey.DeliveryDocument = append(dataKey.DeliveryDocument, *(deliveryDocumentItem)[i].DeliveryDocument)
	}

	repeat := strings.Repeat("?,", len(dataKey.DeliveryDocument)-1) + "?"
	for _, v := range dataKey.DeliveryDocument {
		args = append(args, v)
	}

	args = append(args, dataKey.ConfirmedDeliveryDateFrom, dataKey.ConfirmedDeliveryDateTo, dataKey.ActualGoodsIssueDateFrom, dataKey.ActualGoodsIssueDateTo, dataKey.ItemCompleteDeliveryIsDefined, dataKey.ItemBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.ItemBillingStatus)

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, DeliveryDocumentItem, ConfirmedDeliveryDate, ActualGoodsIssueDate, ItemCompleteDeliveryIsDefined, ItemBillingStatus, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_item_data
		WHERE DeliveryDocument IN ( `+repeat+` )
		AND ConfirmedDeliveryDate BETWEEN ? AND ?
		AND ActualGoodsIssueDate BETWEEN ? AND ?
		AND (ItemCompleteDeliveryIsDefined, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?)
		AND ItemBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocumentItemInIndividualProcess(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

// func (f *SubFunction) OrdersHeaderPartner(
// 	sdc *api_input_reader.SDC,
// 	psdc *api_processing_data_formatter.SDC,
// ) (*[]api_processing_data_formatter.OrdersHeaderPartner, error) {
// 	var args []interface{}

// 	orderID := psdc.OrderID

// 	repeat := strings.Repeat("?,", len(*orderID)-1) + "?"
// 	for _, tag := range *orderID {
// 		args = append(args, tag.OrderID)
// 	}

// 	rows, err := f.db.Query(
// 		`SELECT OrderID, PartnerFunction, BusinessPartner
// 		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_partner_data
// 		WHERE OrderID IN ( `+repeat+` );`, args...,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data, err := psdc.ConvertToOrdersHeaderPartner(sdc, rows)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return data, err
// }

// func (f *SubFunction) DeliveryDocumentByNumberSpecification(
// 	sdc *api_input_reader.SDC,
// 	psdc *api_processing_data_formatter.SDC,
// ) (*[]api_processing_data_formatter.DeliveryDocument, error) {
// 	var args []interface{}

// 	billFromParty := sdc.InputParameters.BillFromParty
// 	billToParty := sdc.InputParameters.BillToParty

// 	if len(*billFromParty) != len(*billToParty) {
// 		return nil, nil
// 	}

// 	dataKey := psdc.ConvertToDeliveryDocumentByNumberSpecificationKey(sdc, len(*billFromParty))

// 	for i := range *billFromParty {
// 		dataKey.BillFromParty[i] = (*billFromParty)[i]
// 		dataKey.BillToParty[i] = (*billToParty)[i]
// 	}

// 	repeat := strings.Repeat("(?,?),", len(dataKey.BillFromParty)-1) + "(?,?)"
// 	for i := range dataKey.BillFromParty {
// 		args = append(args, dataKey.BillFromParty[i], dataKey.BillToParty[i])
// 	}

// 	args = append(
// 		args,
// 		dataKey.HeaderCompleteDeliveryIsDefined,
// 		dataKey.HeaderDeliveryStatus,
// 		dataKey.HeaderBillingBlockStatus,
// 		dataKey.HeaderBillingStatus,
// 	)

// 	var count *int
// 	err := f.db.QueryRow(
// 		`SELECT COUNT(*)
// 		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
// 		WHERE (BillFromParty, BillToParty) IN ( `+repeat+` )
// 		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus) = (?, ?, ?)
// 		AND HeaderBillingStatus <> ?;`, args...,
// 	).Scan(&count)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if *count == 0 || *count > 1000 {
// 		return nil, xerrors.Errorf("OrderIDの検索結果がゼロ件または1,000件超です。")
// 	}

// 	rows, err := f.db.Query(
// 		`SELECT DeliveryDocument, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingStatus, HeaderBillingBlockStatus
// 		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
// 		WHERE (BillFromParty, BillToParty) IN ( `+repeat+` )
// 		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus) = (?, ?, ?)
// 		AND HeaderBillingStatus <> ?;`, args...,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data, err := psdc.ConvertToDeliveryDocumentByNumberSpecification(sdc, rows)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return data, err
// }

func (f *SubFunction) DeliveryDocumentByRangeSpecification(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentHeader, error) {
	dataKey := psdc.ConvertToDeliveryDocumentByRangeSpecificationKey()

	dataKey.BillFromPartyFrom = sdc.InputParameters.BillFromPartyFrom
	dataKey.BillFromPartyTo = sdc.InputParameters.BillFromPartyTo
	dataKey.BillToPartyFrom = sdc.InputParameters.BillToPartyFrom
	dataKey.BillToPartyTo = sdc.InputParameters.BillToPartyTo

	count := new(int)
	err := f.db.QueryRow(
		`SELECT COUNT(*)
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
		WHERE BillFromParty BETWEEN ? AND ?
		AND BillToParty BETWEEN ? AND ?
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, dataKey.BillFromPartyFrom, dataKey.BillFromPartyTo, dataKey.BillToPartyFrom, dataKey.BillToPartyTo, dataKey.HeaderCompleteDeliveryIsDefined, dataKey.HeaderDeliveryStatus, dataKey.HeaderBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.HeaderBillingStatus,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if *count == 0 || *count > 1000 {
		return nil, xerrors.Errorf("OrderIDの検索結果がゼロ件または1,000件超です。")
	}

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
		WHERE BillFromParty BETWEEN ? AND ?
		AND BillToParty BETWEEN ? AND ?
		AND  (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, dataKey.BillFromPartyFrom, dataKey.BillFromPartyTo, dataKey.BillToPartyFrom, dataKey.BillToPartyTo, dataKey.HeaderCompleteDeliveryIsDefined, dataKey.HeaderDeliveryStatus, dataKey.HeaderBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.HeaderBillingStatus,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocumentByRangeSpecification(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentByReferenceDocument(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentHeader, error) {
	dataKey := psdc.ConvertToDeliveryDocumentByReferenceDocumentKey()

	dataKey.ReferenceDocument = sdc.InputParameters.ReferenceDocument

	count := new(int)
	err := f.db.QueryRow(
		`SELECT COUNT(*)
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
		WHERE (DeliveryDocument, HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, dataKey.ReferenceDocument, dataKey.HeaderCompleteDeliveryIsDefined, dataKey.HeaderDeliveryStatus, dataKey.HeaderBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.HeaderBillingStatus,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if *count == 0 {
		return nil, xerrors.Errorf("OrderIDの検索結果がゼロ件です。")
	}

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus,  HeaderBillingStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
		WHERE (DeliveryDocument, HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, dataKey.ReferenceDocument, dataKey.HeaderCompleteDeliveryIsDefined, dataKey.HeaderDeliveryStatus, dataKey.HeaderBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.HeaderBillingStatus,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocumentByReferenceDocument(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

// func (f *SubFunction) DeliveryDocumentHeaderPartner(
// 	sdc *api_input_reader.SDC,
// 	psdc *api_processing_data_formatter.SDC,
// ) (*[]api_processing_data_formatter.DeliveryDocumentHeaderPartner, error) {
// 	var args []interface{}

// 	deliveryDocument := psdc.DeliveryDocument

// 	repeat := strings.Repeat("?,", len(*deliveryDocument)-1) + "?"
// 	for _, tag := range *deliveryDocument {
// 		args = append(args, tag.DeliveryDocument)
// 	}

// 	rows, err := f.db.Query(
// 		`SELECT DeliveryDocument, PartnerFunction, BusinessPartner
// 		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_partner_data
// 		WHERE DeliveryDocument IN ( `+repeat+` );`, args...,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data, err := psdc.ConvertToDeliveryDocumentHeaderPartner(sdc, rows)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return data, err
// }

func (f *SubFunction) CreateSdc(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
	osdc *dpfm_api_output_formatter.SDC,
) error {
	var err error
	var e error

	wg := sync.WaitGroup{}
	wg.Add(3)

	psdc.MetaData = f.MetaData(sdc, psdc)
	psdc.ProcessType = f.ProcessType(sdc, psdc)
	psdc.ReferenceType = f.ReferenceType(sdc, psdc)

	processType := psdc.ProcessType

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if processType.BulkProcess {

			// // I-1-1. OrderIDの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
			// psdc.OrderID, e = f.OrderIDByArraySpec(sdc, psdc)
			// if e != nil {
			// 	err = e
			// 	return
			// }

			// I-1-1. OrderIDの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
			psdc.OrderID, e = f.OrderIDByRangeSpec(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			//I-1-2. OrderItemの絞り込み
			psdc.OrderItem, e = f.OrderItemInBulkProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			// I-2-1. Delivery Document Headerの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
			psdc.DeliveryDocumentHeader, e = f.DeliveryDocumentByRangeSpecification(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			//I-2-2. Delivery Document Itemの絞り込み
			psdc.DeliveryDocumentItem, e = f.DeliveryDocumentItemInBulkProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}

		} else if processType.IndividualProcess {

			// II-1-1. OrderIDが未請求対象であることの確認
			psdc.OrderID, e = f.OrderIDByReferenceDocument(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			//II-1-2. OrderItemの絞り込み
			psdc.OrderItem, e = f.OrderItemInIndividualProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			// II-2-1. Delivery Document Headerの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
			psdc.DeliveryDocumentHeader, e = f.DeliveryDocumentByReferenceDocument(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			// II-2-2. Delivery Document Itemの絞り込み
			psdc.DeliveryDocumentItem, e = f.DeliveryDocumentItemInIndividualProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}

		}
		// 1-1.オーダー参照レコード・値の取得（オーダーヘッダ）
		psdc.OrdersHeader, e = f.OrdersHeader(sdc, psdc)
		if e != nil {
			err = e
			return
		}

		// 1-2. オーダー参照レコード・値の取得（オーダー明細）
		psdc.OrdersItem, e = f.OrdersItem(sdc, psdc)
		if e != nil {
			err = e
			return
		}

		// 2-1. 入出荷伝票参照レコード・値の取得（入出荷伝票ヘッダ）
		psdc.DeliveryDocumentHeaderData, e = f.DeliveryDocumentHeaderData(sdc, psdc)
		if e != nil {
			err = e
			return
		}

		// 2-2. 入出荷伝票参照レコード・値の取得（入出荷伝票明細）
		psdc.DeliveryDocumentItemData, e = f.DeliveryDocumentItemData(sdc, psdc)
		if e != nil {
			err = e
			return
		}

		//3-1. InvoiceDocumentHeader
		psdc.CalculateInvoiceDocument, e = f.CalculateInvoiceDocument(sdc, psdc)
		if e != nil {
			err = e
			return
		}
		// f.l.Info(psdc.CalculateInvoiceDocument)

		// // 2-5. TotalNetAmount
		// psdc.TotalNetAmount, e = f.TotalNetAmount(sdc, psdc)
		// if e != nil {
		// 	err = e
		// 	return
		// }

		// // I-1-2. ヘッダパートナのデータ取得
		// psdc.OrdersHeaderPartner, e = f.OrdersHeaderPartner(sdc, psdc)
		// if e != nil {
		// 	err = e
		// 	return
		// }

		// // 1-I-2. オーダー参照レコード・値の取得（オーダーヘッダパートナ）
		// psdc.HeaderOrdersHeaderPartner, e = f.HeaderOrdersHeaderPartner(sdc, psdc)
		// if e != nil {
		// 	err = e
		// 	return
		// }

	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// // // I-2-1. Delivery Document Headerの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
		// // psdc.DeliveryDocument, e = f.DeliveryDocumentByNumberSpecification(sdc, psdc)
		// // if e != nil {
		// // 	err = e
		// // 	return
		// // }

		// // I-2-2. ヘッダパートナのデータ取得
		// psdc.DeliveryDocumentHeaderPartner, e = f.DeliveryDocumentHeaderPartner(sdc, psdc)
		// if e != nil {
		// 	err = e
		// 	return
		// }

		// // 1-II-2. 入出荷伝票参照レコード・値の取得（入出荷伝票ヘッダパートナ）
		// psdc.HeaderDeliveryDocumentHeaderPartner, e = f.HeaderDeliveryDocumentHeaderPartner(sdc, psdc)
		// if e != nil {
		// 	err = e
		// 	return
		// }

		// // // 2-5. TotalNetAmount
		// // psdc.TotalNetAmount, e = f.TotalNetAmount(sdc, psdc)
		// // if e != nil {
		// // 	err = e
		// // 	return
		// // }
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// // 1-1 InvoiceDocument
		// psdc.CalculateInvoiceDocument, e = f.CalculateInvoiceDocument(sdc, psdc)
		// if e != nil {
		// 	err = e
		// 	return
		// }
	}(&wg)

	wg.Wait()
	if err != nil {
		return err
	}

	f.l.Info(psdc)
	err = f.SetValue(sdc, psdc, osdc)
	if err != nil {
		return err
	}

	return nil
}
