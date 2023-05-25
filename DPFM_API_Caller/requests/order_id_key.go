package requests

type OrderIDKey struct {
	BillFromParty                   []*int `json:"BillFromParty"`
	BillFromPartyFrom               *int   `json:"BillFromPartyFrom"`
	BillFromPartyTo                 *int   `json:"BillFromPartyTo"`
	BillToParty                     []*int `json:"BillToParty"`
	BillToPartyFrom                 *int   `json:"BillToPartyFrom"`
	BillToPartyTo                   *int   `json:"BillToPartyTo"`
	HeaderCompleteDeliveryIsDefined bool   `json:"HeaderCompleteDeliveryIsDefined"`
	HeaderDeliveryStatus            string `json:"HeaderDeliveryStatus"`
	HeaderBillingStatus             string `json:"HeaderBillingStatus"`
	HeaderBillingBlockStatus        bool   `json:"HeaderBillingBlockStatus"`
	IsCancelled                     bool   `json:"IsCancelled"`
	IsMarkedForDeletion             bool   `json:"IsMarkedForDeletion"`
	ReferenceDocument               int    `json:"ReferenceDocument"`
}
