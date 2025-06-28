package propostahandler

type GetPropostasParams struct {
	UsuarioID string
	Status    *string
	Tipo      string
	Limit     int
	Offset    int
}

type StatusUpdatePayload struct {
	Status string `json:"status"`
}
