package requests

type OrderItemKey struct {
	OrderID                       []int  `json:"OrderID"`
	ItemCompleteDeliveryIsDefined bool   `json:"ItemCompleteDeliveryIsDefined"`
	ItemDeliveryStatus            string `json:"ItemDeliveryStatus"`
	ItemBillingStatus             string `json:"ItemBillingStatus"`
	ItemBillingBlockStatus        bool   `json:"ItemBillingBlockStatus"`
	IsCancelled                   *bool  `json:"IsCancelled"`
	IsMarkedForDeletion           *bool  `json:"IsMarkedForDeletion"`
}
