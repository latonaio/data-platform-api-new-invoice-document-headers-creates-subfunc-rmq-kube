package requests

type DeliveryDocumentItemKey struct {
	DeliveryDocument              []int     `json:"DeliveryDocument"`
	ConfirmedDeliveryDate         []*string `json:"ConfirmedDeliveryDate"`
	ConfirmedDeliveryDateFrom     *string   `json:"BillFromPartyFrom"`
	ConfirmedDeliveryDateTo       *string   `json:"BillFromPartyTo"`
	ActualGoodsIssueDate          []*string `json:"ActualGoodsIssueDate"`
	ActualGoodsIssueDateFrom      *string   `json:"BillToPartyFrom"`
	ActualGoodsIssueDateTo        *string   `json:"BillToPartyTo"`
	ItemCompleteDeliveryIsDefined bool      `json:"ItemCompleteDeliveryIsDefined"`
	// ItemDeliveryStatus            string    `json:"ItemDeliveryStatus"`
	ItemBillingStatus      string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus bool   `json:"ItemBillingBlockStatus"`
	IsCancelled            bool   `json:"IsCancelled"`
	IsMarkedForDeletion    bool   `json:"IsMarkedForDeletion"`
}
