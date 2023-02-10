package requests

type DeliveryDocumentItem struct {
	DeliveryDocument              int     `json:"DeliveryDocument"`
	DeliveryDocumentItem          int     `json:"DeliveryDocumentItem"`
	ConfirmedDeliveryDate         *string `json:"ConfirmedDeliveryDate"`
	ActualGoodsIssueDate          *string `json:"ActualGoodsIssueDate"`
	ItemCompleteDeliveryIsDefined *bool   `json:"ItemCompleteDeliveryIsDefined"`
	// ItemDeliveryStatus            string  `json:"ItemDeliveryStatus"`
	ItemBillingStatus      *string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus *bool   `json:"ItemBillingBlockStatus"`
	IsCancelled            *bool   `json:"IsCancelled"`
	IsMarkedForDeletion    *bool   `json:"IsMarkedForDeletion"`
}
