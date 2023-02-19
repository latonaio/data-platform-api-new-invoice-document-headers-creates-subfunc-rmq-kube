package requests

type DeliveryDocumentItemKey struct {
	DeliveryDocument              []int     `json:"DeliveryDocument"`
	ConfirmedDeliveryDateFrom     *string   `json:"BillFromPartyFrom"`
	ConfirmedDeliveryDateTo       *string   `json:"BillFromPartyTo"`
	ActualGoodsIssueDateFrom      *string   `json:"BillToPartyFrom"`
	ActualGoodsIssueDateTo        *string   `json:"BillToPartyTo"`
	ConfirmedDeliveryDate         []*string `json:"ConfirmedDeliveryDate"`
	ActualGoodsIssueDate          []*string `json:"ActualGoodsIssueDate"`
	ItemCompleteDeliveryIsDefined bool      `json:"ItemCompleteDeliveryIsDefined"`
	// ItemDeliveryStatus            string    `json:"ItemDeliveryStatus"`
	ItemBillingStatus      string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus bool   `json:"ItemBillingBlockStatus"`
	IsCancelled            *bool  `json:"IsCancelled"`
	IsMarkedForDeletion    *bool  `json:"IsMarkedForDeletion"`
}
