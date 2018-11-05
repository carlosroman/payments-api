package payment

type Payment struct {
	Id             string     `json:"id"`
	Type           string     `json:"type"`
	Version        int32      `json:"version"`
	OrganisationId string     `json:"organisation_id"`
	Attributes     Attributes `json:"attributes"`
}

type Attributes struct {
	Amount            string `json:"amount"`
	PaymentId         string `json:"payment_id"`
	PaymentType       string `json:"payment_type"`
	Currency          string `json:"currency"`
	EndToEndReference string `json:"end_to_end_reference"`
	NumericReference  string `json:"numeric_reference"`
	Reference         string `json:"reference"`
}

type Payments struct {
	Payments []Payment `json:"data"`
}
