package payment

type Payment struct {
	Id                string     `json:"id"`
	Version           int32      `json:"version"`
	OrganisationId    string     `json:"organisation_id"`
	Attributes        Attributes `json:"attributes"`
	EndToEndReference string     `json:"end_to_end_reference"`
	NumericReference  string     `json:"numeric_reference"`
	PaymentId         string     `json:"payment_id"`
	Reference         string     `json:"reference"`
}

type Attributes struct {
	Amount string `json:"amount"`
}

type Payments struct {
	Payments []Payment `json:"data"`
}
