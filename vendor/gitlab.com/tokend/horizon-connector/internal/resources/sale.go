package resources

type CoreSale struct {
	ID uint64 `json:"id"`
}

type Sale struct {
	ID      string      `json:"id"`
	Details SaleDetails `json:"details"`
}

type SaleDetails struct {
	Name string `json:"name"`
}

func (s *Sale) Name() string {
	return s.Details.Name
}
