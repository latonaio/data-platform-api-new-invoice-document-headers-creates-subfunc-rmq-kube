package requests

type DeliveryDocumentHeader struct {
	DeliveryDocument                int    `json:"DeliveryDocument"`
	BillFromParty                   *int   `json:"BillFromParty"`
	BillToParty                     *int   `json:"BillToParty"`
	HeaderCompleteDeliveryIsDefined *bool  `json:"HeaderCompleteDeliveryIsDefined"`
	HeaderDeliveryStatus            string `json:"HeaderDeliveryStatus"`
	HeaderBillingStatus             string `json:"HeaderBillingStatus"`
	HeaderBillingBlockStatus        *bool  `json:"HeaderBillingBlockStatus"`
	IsCancelled                     *bool  `json:"IsCancelled"`
	IsMarkedForDeletion             *bool  `json:"IsMarkedForDeletion"`
}
