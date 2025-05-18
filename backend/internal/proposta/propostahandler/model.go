package propostahandler

type GetPropostasParams struct {
	UsuarioID string
	Status    *string
	Tipo      string
}

type StatusUpdatePayload struct {
	Status string `json:"status"`
}
