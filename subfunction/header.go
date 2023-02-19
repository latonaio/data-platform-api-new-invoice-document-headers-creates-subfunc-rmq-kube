package subfunction

import (
	api_input_reader "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Input_Reader"
	api_processing_data_formatter "data-platform-api-invoice-document-headers-creates-subfunc-rmq/API_Processing_Data_Formatter"
	"sort"
	"strings"
	"time"

	"golang.org/x/xerrors"
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

	if sdc.Header.TotalNetAmount != nil {
		totalNetAmount = *sdc.Header.TotalNetAmount
	}

	data := psdc.ConvertToTotalNetAmount(&totalNetAmount)

	return data
}

func (f *SubFunction) TotalTaxAmount(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.TotalTaxAmount {
	totalTaxAmount := float32(0)

	if psdc.ReferenceType.OrderID {
		for _, v := range psdc.OrdersItem {
			if v.TaxAmount != nil {
				totalTaxAmount += *v.TaxAmount
			}
		}
	} else if psdc.ReferenceType.DeliveryDocument {
		for _, v := range psdc.DeliveryDocumentItemData {
			if v.TaxAmount != nil {
				totalTaxAmount += *v.TaxAmount
			}
		}
	}

	if sdc.Header.TotalTaxAmount != nil {
		totalTaxAmount = *sdc.Header.TotalTaxAmount
	}

	data := psdc.ConvertToTotalTaxAmount(&totalTaxAmount)

	return data
}

func (f *SubFunction) TotalGrossAmount(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.TotalGrossAmount {
	totalGrossAmount := float32(0)

	if psdc.ReferenceType.OrderID {
		for _, v := range psdc.OrdersItem {
			if v.GrossAmount != nil {
				totalGrossAmount += *v.GrossAmount
			}
		}
	} else if psdc.ReferenceType.DeliveryDocument {
		for _, v := range psdc.DeliveryDocumentItemData {
			if v.GrossAmount != nil {
				totalGrossAmount += *v.GrossAmount
			}
		}
	}

	if sdc.Header.TotalGrossAmount != nil {
		totalGrossAmount = *sdc.Header.TotalGrossAmount
	}

	data := psdc.ConvertToTotalGrossAmount(&totalGrossAmount)

	return data
}

func (f *SubFunction) InvoiceDocumentDate(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.InvoiceDocumentDate {
	invoiceDocumentDate := string("")

	if sdc.Header.InvoiceDocumentDate != nil {
		if *sdc.Header.InvoiceDocumentDate != "" {
			invoiceDocumentDate = *sdc.Header.InvoiceDocumentDate
		}
	}

	data := psdc.ConvertToInvoiceDocumentDate(&invoiceDocumentDate)

	return data
}

func (f *SubFunction) PaymentDueDate(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) (*api_processing_data_formatter.PaymentDueDate, error) {
	paymentTerms := ""

	if psdc.ReferenceType.OrderID {
		paymentTerms = psdc.OrdersItem[0].PaymentTerms
	} else if psdc.ReferenceType.DeliveryDocument {
		paymentTerms = *psdc.DeliveryDocumentItemData[0].PaymentTerms
	}

	rows, err := f.db.Query(
		`SELECT PaymentTerms, BaseDate, BaseDateCalcAddMonth, BaseDateCalcFixedDate, PaymentDueDateCalcAddMonth, PaymentDueDateCalcFixedDate
			FROM  DataPlatformMastersAndTransactionsMysqlKube.data_platform_payment_terms_payment_terms_data
			WHERE PaymentTerms = ?`, paymentTerms,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	psdc.PaymentTerms, err = psdc.ConvertToPaymentTerms(rows)
	if err != nil {
		return nil, err
	}

	if *sdc.Header.PaymentDueDate != "" {
		data := psdc.ConvertToPaymentDueDate(sdc)
		return data, nil
	}

	invoiceDocumentDate := psdc.InvoiceDocumentDate

	paymentDueDate, err := caluculatePaymentDueDate(psdc, *invoiceDocumentDate.InvoiceDocumentDate, psdc.PaymentTerms)

	data := psdc.ConvertToCalculatePaymentDueDate(sdc, paymentDueDate)

	return data, err
}

func caluculatePaymentDueDate(
	psdc *api_processing_data_formatter.SDC,
	invoiceDocumentDate string,
	paymentTerms []*api_processing_data_formatter.PaymentTerms,
) (*string, error) {
	format := "2006-01-02"
	t, err := time.Parse(format, invoiceDocumentDate)
	if err != nil {
		return nil, err
	}

	sort.Slice(paymentTerms, func(i, j int) bool {
		return paymentTerms[i].BaseDate < paymentTerms[j].BaseDate
	})

	day := t.Day()
	for i, v := range paymentTerms {
		if day <= v.BaseDate {
			t = time.Date(t.Year(), t.Month()+time.Month(*v.BaseDateCalcAddMonth)+1, 0, 0, 0, 0, 0, time.UTC)
			if *v.BaseDateCalcFixedDate == 31 {
				t = time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			} else {
				t = time.Date(t.Year(), t.Month(), *v.BaseDateCalcFixedDate, 0, 0, 0, 0, time.UTC)
			}
			break
		}
		if i == len(paymentTerms)-1 {
			return nil, xerrors.Errorf("'data_platform_payment_terms_payment_terms_data'テーブルが不適切です。")

		}
	}

	res := getStringPtr(t.Format(format))

	return res, nil
}

func (f *SubFunction) NetPaymentDays(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) (*api_processing_data_formatter.NetPaymentDays, error) {

	if sdc.Header.NetPaymentDays != nil {
		data := psdc.ConvertToNetPaymentDays(sdc)
		return data, nil
	}

	paymentDueDate := psdc.PaymentDueDate.PaymentDueDate

	calculateNetPaymentDays, err := calculateNetPaymentDays(*psdc.InvoiceDocumentDate.InvoiceDocumentDate, *paymentDueDate)

	netPaymentDays := calculateNetPaymentDays

	data := psdc.ConvertToCalculateNetPaymentDays(sdc, paymentDueDate, netPaymentDays)
	return data, err
}

func calculateNetPaymentDays(
	invoiceDocumentDate string,
	paymentDueDate string,
) (int, error) {
	format := "2006-01-02"
	tb, err := time.Parse(format, invoiceDocumentDate)
	if err != nil {
		return 0, err
	}

	tp, err := time.Parse(format, paymentDueDate)
	if err != nil {
		return 0, err
	}

	res := int(tp.Sub(tb).Hours() / 24)

	return res, nil
}

func (f *SubFunction) CalculateInvoiceDocument(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) ([]*api_processing_data_formatter.CalculateInvoiceDocument, error) {
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
	defer rows.Close()

	dataQueryGets, err := psdc.ConvertToInvoiceDocumentHeaderQueryGets(sdc, rows)
	if err != nil {
		return nil, err
	}

	if dataQueryGets.InvoiceDocumentLatestNumber == nil {
		return nil, xerrors.Errorf("'data_platform_number_range_latest_number_data'テーブルのLatestNumberがNULLです。")
	}

	billParties := make([]api_processing_data_formatter.BillParty, 0)
	if psdc.ReferenceType.OrderID {
		for _, ordersHeader := range psdc.OrdersHeader {
			billFromParty := ordersHeader.BillFromParty
			billToParty := ordersHeader.BillToParty
			orderID := ordersHeader.OrderID

			if billFromParty == nil || billToParty == nil {
				continue
			}

			if billPartyContain(billParties, *billFromParty, *billToParty) {
				continue
			}

			billParties = append(billParties, api_processing_data_formatter.BillParty{
				BillFromParty: *billFromParty,
				BillToParty:   *billToParty,
				OrderID:       orderID,
			})
		}
	} else if psdc.ReferenceType.DeliveryDocument {
		for _, deliveryDocumentItemData := range psdc.DeliveryDocumentItemData {
			billFromParty := deliveryDocumentItemData.BillFromParty
			billToParty := deliveryDocumentItemData.BillToParty
			deliveryDocument := deliveryDocumentItemData.DeliveryDocument
			deliveryDocumentItem := deliveryDocumentItemData.DeliveryDocumentItem

			if billFromParty == nil || billToParty == nil {
				continue
			}

			if billPartyContain(billParties, *billFromParty, *billToParty) {
				continue
			}

			billParties = append(billParties, api_processing_data_formatter.BillParty{
				BillFromParty:        *billFromParty,
				BillToParty:          *billToParty,
				DeliveryDocument:     deliveryDocument,
				DeliveryDocumentItem: deliveryDocumentItem,
			})
		}
	}

	data := make([]*api_processing_data_formatter.CalculateInvoiceDocument, 0)
	for i, billParty := range billParties {
		invoiceDocumentLatestNumber := dataQueryGets.InvoiceDocumentLatestNumber
		invoiceDocument := *dataQueryGets.InvoiceDocumentLatestNumber + i + 1
		billFromParty := billParty.BillFromParty
		billToParty := billParty.BillToParty
		orderID := billParty.OrderID
		deliveryDocument := billParty.DeliveryDocument
		deliveryDocumentItem := billParty.DeliveryDocumentItem

		datum := psdc.ConvertToCalculateInvoiceDocument(invoiceDocumentLatestNumber, invoiceDocument, orderID, deliveryDocument, deliveryDocumentItem, billFromParty, billToParty)
		data = append(data, datum)
	}

	return data, err
}

func (f *SubFunction) CreationDateHeader(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.CreationDate {
	data := psdc.ConvertToCreationDateHeader(getSystemDate())

	return data
}

func (f *SubFunction) LastChangeDateHeader(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.LastChangeDate {
	data := psdc.ConvertToLastChangeDateHeader(getSystemDate())

	return data
}

func (f *SubFunction) CreationTimeHeader(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.CreationTime {
	data := psdc.ConvertToCreationTimeHeader(getSystemTime())

	return data
}

func (f *SubFunction) LastChangeTimeHeader(
	sdc *api_input_reader.SDC,
	psdc *api_processing_data_formatter.SDC,
) *api_processing_data_formatter.LastChangeTime {
	data := psdc.ConvertToLastChangeTimeHeader(getSystemTime())

	return data
}

func billPartyContain(billParties []api_processing_data_formatter.BillParty, billFromParty, billToParty int) bool {
	for _, billParty := range billParties {
		if billFromParty == billParty.BillFromParty && billToParty == billParty.BillToParty {
			return true
		}
	}
	return false
}

func getSystemDate() string {
	day := time.Now()
	return day.Format("2006-01-02")
}

func getSystemTime() string {
	day := time.Now()
	return day.Format("15:04:05")
}

func getStringPtr(s string) *string {
	return &s
}
