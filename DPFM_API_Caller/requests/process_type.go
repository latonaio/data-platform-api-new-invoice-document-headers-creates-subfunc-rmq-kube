package requests

type ProcessType struct {
	BulkProcess       bool `json:"BulkProcess"`
	IndividualProcess bool `json:"IndividualProcess"`
	ArraySpec         bool `json:"ArraySpec"`
	RangeSpec         bool `json:"RangeSpec"`
}
