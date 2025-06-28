package propostasvc

type PropostasFilter struct {
	UsuarioID string
	Status    *string
	Tipo      string
	Limit     int
	Offset    int
}

var statusMap = map[string]bool{
	"pendente": true,
	"aceita":   true,
	"recusada": true,
	"expirada": true,
}
