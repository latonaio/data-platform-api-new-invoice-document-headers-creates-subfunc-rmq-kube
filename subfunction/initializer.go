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
	metaData := psdc.ConvertToMetaData(sdc)

	return metaData
}

func (f *SubFunction) ReferenceType(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) (*api_processing_data_formatter.ReferenceType, error) {
	referenceType := psdc.ConvertToReferenceType()

	if sdc.InputParameters.ReferenceDocument == nil {
		return nil, xerrors.Errorf("入力のReferenceDocumentがnullです。")
	}

	referenceDocument := *sdc.InputParameters.ReferenceDocument
	if 1 <= referenceDocument && referenceDocument <= 9999999 {
		referenceType.OrderID = true
	} else if 80000000 <= referenceDocument && referenceDocument <= 89999999 {
		referenceType.DeliveryDocument = true
	} else {
		return nil, xerrors.Errorf("入力のReferenceDocumentがOrderIDとDeliveryDocumentの範囲にありません。")
	}

	return referenceType, nil
}

func (f *SubFunction) ProcessType(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) (*api_processing_data_formatter.ProcessType, error) {
	referenceType := psdc.ReferenceType
	processType := psdc.ConvertToProcessType()

	if isBulkProcess(sdc, processType, referenceType) {
		processType.BulkProcess = true
	}

	if isIndividualProcess(sdc, processType, referenceType) {
		processType.IndividualProcess = true
	}

	if !processType.BulkProcess && !processType.IndividualProcess {
		return nil, xerrors.Errorf("一括登録または個別登録に必要な入力パラメータが揃っていません。")
	}

	return processType, nil
}

func isBulkProcess(
	sdc *api_input_reader.SDC,
	processType *api_processing_data_formatter.ProcessType,
	referenceType *api_processing_data_formatter.ReferenceType,
) bool {
	inputParameters := sdc.InputParameters

	if referenceType.OrderID {
		if inputParameters.BillFromParty != nil && inputParameters.BillToParty != nil {
			if (*inputParameters.BillFromParty)[0] != nil && (*inputParameters.BillToParty)[0] != nil {
				processType.ArraySpec = true
				return true
			}
		}
		if inputParameters.BillFromPartyTo != nil && inputParameters.BillFromPartyFrom != nil &&
			inputParameters.BillToPartyTo != nil && inputParameters.BillToPartyFrom != nil {
			processType.RangeSpec = true
			return true
		}
	} else if referenceType.DeliveryDocument {
		if inputParameters.BillFromParty != nil && inputParameters.BillToParty != nil &&
			inputParameters.ConfirmedDeliveryDate != nil && inputParameters.ActualGoodsIssueDate != nil {
			if (*inputParameters.BillFromParty)[0] != nil && (*inputParameters.BillToParty)[0] != nil &&
				(*inputParameters.ConfirmedDeliveryDate)[0] != nil && (*inputParameters.ActualGoodsIssueDate)[0] != nil {
				if len(*(*inputParameters.ConfirmedDeliveryDate)[0]) != 0 && len(*(*inputParameters.ActualGoodsIssueDate)[0]) != 0 {
					processType.ArraySpec = true
					return true
				}
			}
		}
		if inputParameters.BillFromPartyTo != nil && inputParameters.BillFromPartyFrom != nil &&
			inputParameters.BillToPartyTo != nil && inputParameters.BillToPartyFrom != nil &&
			inputParameters.ConfirmedDeliveryDateTo != nil && inputParameters.ConfirmedDeliveryDateFrom != nil &&
			inputParameters.ActualGoodsIssueDateTo != nil && inputParameters.ActualGoodsIssueDateFrom != nil {
			if len(*inputParameters.ConfirmedDeliveryDateTo) != 0 && len(*inputParameters.ConfirmedDeliveryDateFrom) != 0 &&
				len(*inputParameters.ActualGoodsIssueDateTo) != 0 && len(*inputParameters.ActualGoodsIssueDateFrom) != 0 {
				processType.RangeSpec = true
				return true
			}
		}
	}

	return false
}

func isIndividualProcess(
	sdc *api_input_reader.SDC,
	processType *api_processing_data_formatter.ProcessType,
	referenceType *api_processing_data_formatter.ReferenceType,
) bool {
	inputParameters := sdc.InputParameters

	if inputParameters.ReferenceDocument != nil {
		if referenceType.OrderID {
			return true
		} else if referenceType.DeliveryDocument {
			if inputParameters.ConfirmedDeliveryDate != nil && inputParameters.ActualGoodsIssueDate != nil {
				if (*inputParameters.ConfirmedDeliveryDate)[0] != nil && (*inputParameters.ActualGoodsIssueDate)[0] != nil {
					if len(*(*inputParameters.ConfirmedDeliveryDate)[0]) != 0 && len(*(*inputParameters.ActualGoodsIssueDate)[0]) != 0 {
						processType.ArraySpec = true
						return true
					}
				}
			}
			if inputParameters.ConfirmedDeliveryDateTo != nil && inputParameters.ConfirmedDeliveryDateFrom != nil &&
				inputParameters.ActualGoodsIssueDateTo != nil && inputParameters.ActualGoodsIssueDateFrom != nil {
				if len(*inputParameters.ConfirmedDeliveryDateTo) != 0 && len(*inputParameters.ConfirmedDeliveryDateFrom) != 0 &&
					len(*inputParameters.ActualGoodsIssueDateTo) != 0 && len(*inputParameters.ActualGoodsIssueDateFrom) != 0 {
					processType.RangeSpec = true
					return true
				}
			}
		}
	}

	return false
}

func (f *SubFunction) OrderIDInBulkProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderID, error) {
	data := make([]*api_processing_data_formatter.OrderID, 0)
	var err error

	processType := psdc.ProcessType

	if processType.ArraySpec {
		data, err = f.OrderIDByArraySpec(sdc, psdc)
		if err != nil {
			return nil, err
		}
	} else if processType.RangeSpec {
		data, err = f.OrderIDByRangeSpec(sdc, psdc)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, xerrors.Errorf("OrderIDの絞り込み（一括登録）に必要な入力パラメータが揃っていません。")
	}

	return data, nil
}

func (f *SubFunction) OrderIDByArraySpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderID, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToOrderIDKey()

	billFromParty := sdc.InputParameters.BillFromParty
	billToParty := sdc.InputParameters.BillToParty

	dataKey.BillFromParty = append(dataKey.BillFromParty, *billFromParty...)
	dataKey.BillToParty = append(dataKey.BillToParty, *billToParty...)

	repeat1 := strings.Repeat("?,", len(dataKey.BillFromParty)-1) + "?"
	for _, v := range dataKey.BillFromParty {
		args = append(args, v)
	}
	repeat2 := strings.Repeat("?,", len(dataKey.BillToParty)-1) + "?"
	for _, v := range dataKey.BillToParty {
		args = append(args, v)
	}

	args = append(
		args,
		dataKey.HeaderCompleteDeliveryIsDefined,
		dataKey.HeaderDeliveryStatus,
		dataKey.HeaderBillingBlockStatus,
		dataKey.IsCancelled,
		dataKey.IsMarkedForDeletion,
		dataKey.HeaderBillingStatus,
	)

	var count *int
	err := f.db.QueryRow(
		`SELECT COUNT(*)
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE BillFromParty IN ( `+repeat1+` )
		AND BillToParty IN ( `+repeat2+` )
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, args...,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if *count == 0 || *count > 1000 {
		return nil, xerrors.Errorf("OrderIDの検索結果がゼロ件または1,000件超です。")
	}

	rows, err := f.db.Query(
		`SELECT OrderID, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus, HeaderBillingStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_orders_header_data
		WHERE BillFromParty IN ( `+repeat1+` )
		AND BillToParty IN ( `+repeat2+` )
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToOrderID(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrderIDByRangeSpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderID, error) {
	dataKey := psdc.ConvertToOrderIDKey()

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
		`SELECT OrderID, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus, HeaderBillingStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion
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

	data, err := psdc.ConvertToOrderID(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrderIDInIndividualProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderID, error) {
	dataKey := psdc.ConvertToOrderIDInIndividualProcessKey()

	dataKey.ReferenceDocument = *sdc.InputParameters.ReferenceDocument

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

	data, err := psdc.ConvertToOrderIDInIndividualProcess(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) OrderItem(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.OrderItem, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToOrderItemKey()

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

	data, err := psdc.ConvertToOrderItem(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentInBulkProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentHeader, error) {
	data := make([]*api_processing_data_formatter.DeliveryDocumentHeader, 0)
	var err error

	processType := psdc.ProcessType

	if processType.ArraySpec {
		data, err = f.DeliveryDocumentByArraySpec(sdc, psdc)
		if err != nil {
			return nil, err
		}
	} else if processType.RangeSpec {
		data, err = f.DeliveryDocumentByRangeSpec(sdc, psdc)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, xerrors.Errorf("DeliveryDocumentの絞り込み（一括登録）に必要な入力パラメータが揃っていません。")
	}

	return data, nil
}

func (f *SubFunction) DeliveryDocumentByArraySpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentHeader, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToDeliveryDocumentKey()

	billFromParty := sdc.InputParameters.BillFromParty
	billToParty := sdc.InputParameters.BillToParty

	dataKey.BillFromParty = append(dataKey.BillFromParty, *billFromParty...)
	dataKey.BillToParty = append(dataKey.BillToParty, *billToParty...)

	repeat1 := strings.Repeat("?,", len(dataKey.BillFromParty)-1) + "?"
	for _, v := range dataKey.BillFromParty {
		args = append(args, v)
	}
	repeat2 := strings.Repeat("?,", len(dataKey.BillToParty)-1) + "?"
	for _, v := range dataKey.BillToParty {
		args = append(args, v)
	}

	args = append(
		args,
		dataKey.HeaderCompleteDeliveryIsDefined,
		dataKey.HeaderDeliveryStatus,
		dataKey.HeaderBillingBlockStatus,
		dataKey.IsCancelled,
		dataKey.IsMarkedForDeletion,
		dataKey.HeaderBillingStatus,
	)

	count := new(int)
	err := f.db.QueryRow(
		`SELECT COUNT(*)
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
		WHERE BillFromParty IN ( `+repeat1+` )
		AND BillToParty IN ( `+repeat2+` )
		AND (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, args...,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if *count == 0 || *count > 1000 {
		return nil, xerrors.Errorf("DeliveryDocumentの検索結果がゼロ件または1,000件超です。")
	}

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus, HeaderBillingStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_header_data
		WHERE BillFromParty IN ( `+repeat1+` )
		AND BillToParty IN ( `+repeat2+` )
		AND  (HeaderCompleteDeliveryIsDefined, HeaderDeliveryStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?, ?)
		AND HeaderBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocument(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentByRangeSpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentHeader, error) {
	dataKey := psdc.ConvertToDeliveryDocumentKey()

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
		return nil, xerrors.Errorf("DeliveryDocumentの検索結果がゼロ件または1,000件超です。")
	}

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, BillFromParty, BillToParty, HeaderCompleteDeliveryIsDefined,
		HeaderDeliveryStatus, HeaderBillingStatus, HeaderBillingBlockStatus, IsCancelled, IsMarkedForDeletion
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

	data, err := psdc.ConvertToDeliveryDocument(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentInIndividualProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentHeader, error) {
	dataKey := psdc.ConvertToDeliveryDocumentInIndividualProcessKey()

	dataKey.ReferenceDocument = *sdc.InputParameters.ReferenceDocument

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

	data, err := psdc.ConvertToDeliveryDocumentInIndividualProcess(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentItem(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentItem, error) {
	data := make([]*api_processing_data_formatter.DeliveryDocumentItem, 0)
	var err error

	processType := psdc.ProcessType

	if processType.ArraySpec {
		data, err = f.DeliveryDocumentItemByArraySpec(sdc, psdc)
		if err != nil {
			return nil, err
		}
	} else if processType.RangeSpec {
		data, err = f.DeliveryDocumentItemByRangeSpec(sdc, psdc)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, xerrors.Errorf("DeliveryDocumentItemの絞り込み（一括登録または個別登録）に必要な入力パラメータが揃っていません。")
	}

	return data, nil
}

func (f *SubFunction) DeliveryDocumentItemByArraySpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentItem, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToDeliveryDocumentItemKey()

	deliveryDocumentItem := psdc.DeliveryDocumentHeader

	for i := range deliveryDocumentItem {
		dataKey.DeliveryDocument = append(dataKey.DeliveryDocument, deliveryDocumentItem[i].DeliveryDocument)
	}

	confirmedDeliveryDate := sdc.InputParameters.ConfirmedDeliveryDate
	atualGoodsIssueDate := sdc.InputParameters.ActualGoodsIssueDate

	dataKey.ConfirmedDeliveryDate = append(dataKey.ConfirmedDeliveryDate, *confirmedDeliveryDate...)
	dataKey.ActualGoodsIssueDate = append(dataKey.ActualGoodsIssueDate, *atualGoodsIssueDate...)

	repeat1 := strings.Repeat("?,", len(dataKey.DeliveryDocument)-1) + "?"
	for _, v := range dataKey.DeliveryDocument {
		args = append(args, v)
	}
	repeat2 := strings.Repeat("?,", len(dataKey.ConfirmedDeliveryDate)-1) + "?"
	for _, v := range dataKey.ConfirmedDeliveryDate {
		args = append(args, v)
	}
	repeat3 := strings.Repeat("?,", len(dataKey.ActualGoodsIssueDate)-1) + "?"
	for _, v := range dataKey.ActualGoodsIssueDate {
		args = append(args, v)
	}

	args = append(args, dataKey.ItemCompleteDeliveryIsDefined, dataKey.ItemBillingBlockStatus, dataKey.IsCancelled, dataKey.IsMarkedForDeletion, dataKey.ItemBillingStatus)

	rows, err := f.db.Query(
		`SELECT DeliveryDocument, DeliveryDocumentItem, ConfirmedDeliveryDate, ActualGoodsIssueDate,
		ItemCompleteDeliveryIsDefined, ItemBillingStatus, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_delivery_document_item_data
		WHERE DeliveryDocument IN ( `+repeat1+` )
		AND ConfirmedDeliveryDate IN ( `+repeat2+` )
		AND ActualGoodsIssueDate IN ( `+repeat3+` )
		AND (ItemCompleteDeliveryIsDefined, ItemBillingBlockStatus, IsCancelled, IsMarkedForDeletion) = (?, ?, ?, ?)
		AND ItemBillingStatus <> ?;`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := psdc.ConvertToDeliveryDocumentItem(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) DeliveryDocumentItemByRangeSpec(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.DeliveryDocumentItem, error) {
	args := make([]interface{}, 0)

	dataKey := psdc.ConvertToDeliveryDocumentItemKey()

	dataKey.ConfirmedDeliveryDateFrom = sdc.InputParameters.ConfirmedDeliveryDateFrom
	dataKey.ConfirmedDeliveryDateTo = sdc.InputParameters.ConfirmedDeliveryDateTo
	dataKey.ActualGoodsIssueDateFrom = sdc.InputParameters.ActualGoodsIssueDateFrom
	dataKey.ActualGoodsIssueDateTo = sdc.InputParameters.ActualGoodsIssueDateTo

	deliveryDocumentItem := psdc.DeliveryDocumentHeader

	for i := range deliveryDocumentItem {
		dataKey.DeliveryDocument = append(dataKey.DeliveryDocument, deliveryDocumentItem[i].DeliveryDocument)
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

	data, err := psdc.ConvertToDeliveryDocumentItem(rows)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (f *SubFunction) CreateSdc(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
	osdc *dpfm_api_output_formatter.SDC,
) error {
	var err error

	psdc.MetaData = f.MetaData(sdc, psdc)
	psdc.ReferenceType, err = f.ReferenceType(sdc, psdc)
	if err != nil {
		return err
	}
	psdc.ProcessType, err = f.ProcessType(sdc, psdc)
	if err != nil {
		return err
	}

	referenceType := psdc.ReferenceType
	if referenceType.OrderID {
		err = f.OrdersReferenceProcess(sdc, psdc, osdc)
		if err != nil {
			return err
		}
	} else if referenceType.DeliveryDocument {
		err = f.DeliveryDocumentReferenceProcess(sdc, psdc, osdc)
		if err != nil {
			return err
		}
	}

	f.l.Info(psdc)
	err = f.SetValue(sdc, psdc, osdc)
	if err != nil {
		return err
	}

	return nil
}

func (f *SubFunction) OrdersReferenceProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
	osdc *dpfm_api_output_formatter.SDC,
) error {
	var err error
	var e error

	processType := psdc.ProcessType
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if processType.BulkProcess {
			// I-1-1. OrderIDの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
			psdc.OrderID, e = f.OrderIDInBulkProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}
		} else if processType.IndividualProcess {
			// II-1-1. OrderIDが未請求対象であることの確認
			psdc.OrderID, e = f.OrderIDInIndividualProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}
		}

		// II-1-2. OrderItemの絞り込み  //I-1-1またはII-1-1
		psdc.OrderItem, e = f.OrderItem(sdc, psdc)
		if e != nil {
			err = e
			return
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// 1-1.オーダー参照レコード・値の取得（オーダーヘッダ）  //I-1-2, II-1-2
			psdc.OrdersHeader, e = f.OrdersHeader(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			//3-1. InvoiceDocumentHeader //1-1
			psdc.CalculateInvoiceDocument = f.CalculateInvoiceDocument(sdc, psdc)

		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// 1-2. オーダー参照レコード・値の取得（オーダー明細）  //I-1-2, II-1-2
			psdc.OrdersItem, e = f.OrdersItem(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				// 2-5. TotalNetAmount  //1-2
				psdc.TotalNetAmount = f.TotalNetAmount(sdc, psdc)

				//2-6 TotalTaxAmount  //1-2
				psdc.TotalTaxAmount = f.TotalTaxAmount(sdc, psdc)

				//2-7 TotalGrossAmount  //1-2
				psdc.TotalGrossAmount = f.TotalGrossAmount(sdc, psdc)
			}(wg)

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				//2-8 InvoiceDocumentDate
				psdc.InvoiceDocumentDate = f.InvoiceDocumentDate(sdc, psdc)

				//2-9  PaymentDueDate  //1-2, 2-8
				psdc.PaymentDueDate, e = f.PaymentDueDate(sdc, psdc)

				//2-10. NetPaymentDays  //2-8,2-9
				psdc.NetPaymentDays, e = f.NetPaymentDays(sdc, psdc)
			}(wg)
		}(wg)
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// 99-1-1. CreationDate(Header)
		psdc.CreationDateHeader = f.CreationDateHeader(sdc, psdc)

		// 99-2-1. LastChangeDate(Header)
		psdc.LastChangeDateHeader = f.LastChangeDateHeader(sdc, psdc)

		//99-3-1. CreationTime(Header)
		psdc.CreationTimeHeader = f.CreationTimeHeader(sdc, psdc)

		//99-4-1. LastChangeTime(Header)
		psdc.LastChangeTimeHeader = f.LastChangeTimeHeader(sdc, psdc)
	}(&wg)

	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (f *SubFunction) DeliveryDocumentReferenceProcess(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
	osdc *dpfm_api_output_formatter.SDC,
) error {
	var err error
	var e error

	processType := psdc.ProcessType
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if processType.BulkProcess {
			// I-2-1. Delivery Document Headerの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
			psdc.DeliveryDocumentHeader, e = f.DeliveryDocumentInBulkProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}
		} else if processType.IndividualProcess {
			// II-2-1. Delivery Document Headerの絞り込み、および、入力パラメータによる請求元と請求先の絞り込み
			psdc.DeliveryDocumentHeader, e = f.DeliveryDocumentInIndividualProcess(sdc, psdc)
			if e != nil {
				err = e
				return
			}
		}

		// I-2-2. Delivery Document Itemの絞り込み  //I-2-1またはII-2-1
		psdc.DeliveryDocumentItem, e = f.DeliveryDocumentItem(sdc, psdc)
		if e != nil {
			err = e
			return
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// 2-1. 入出荷伝票参照レコード・値の取得（入出荷伝票ヘッダ） //II-2-1,II-2-2
			psdc.DeliveryDocumentHeaderData, e = f.DeliveryDocumentHeaderData(sdc, psdc)
			if e != nil {
				err = e
				return
			}

		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// 2-2. 入出荷伝票参照レコード・値の取得（入出荷伝票明細） ////II-2-1,II-2-2
			psdc.DeliveryDocumentItemData, e = f.DeliveryDocumentItemData(sdc, psdc)
			if e != nil {
				err = e
				return
			}

			//3-1. InvoiceDocumentHeader //2-2
			psdc.CalculateInvoiceDocument = f.CalculateInvoiceDocument(sdc, psdc)

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				// 2-5. TotalNetAmount  //2-2
				psdc.TotalNetAmount = f.TotalNetAmount(sdc, psdc)

				//2-6 TotalTaxAmount  //2-2
				psdc.TotalTaxAmount = f.TotalTaxAmount(sdc, psdc)

				//2-7 TotalGrossAmount  //2-2
				psdc.TotalGrossAmount = f.TotalGrossAmount(sdc, psdc)
			}(wg)

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				//2-8 InvoiceDocumentDate
				psdc.InvoiceDocumentDate = f.InvoiceDocumentDate(sdc, psdc)

				//2-9  PaymentDueDate //2-2, 2-8
				psdc.PaymentDueDate, e = f.PaymentDueDate(sdc, psdc)

				//2-10. NetPaymentDays  //2-8,2-9
				psdc.NetPaymentDays, e = f.NetPaymentDays(sdc, psdc)
			}(wg)

		}(wg)

	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// 99-1-1. CreationDate(Header)
		psdc.CreationDateHeader = f.CreationDateHeader(sdc, psdc)

		// 99-2-1. LastChangeDate(Header)
		psdc.LastChangeDateHeader = f.LastChangeDateHeader(sdc, psdc)

		//99-3-1. CreationTime(Header)
		psdc.CreationTimeHeader = f.CreationTimeHeader(sdc, psdc)

		//99-4-1. LastChangeTime(Header)
		psdc.LastChangeTimeHeader = f.LastChangeTimeHeader(sdc, psdc)

	}(&wg)

	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}
