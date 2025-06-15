package usuariorepo

type WriteUsuario struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Senha    string `json:"senha"`
	Telefone string `json:"telefone"`
	Endereco
}

type Endereco struct {
	Rua    string `json:"rua"`
	Cidade string `json:"cidade"`
	Estado string `json:"estado"`
	Bairro string `json:"bairro"`
	Numero int    `json:"numero"`
	CEP    string `json:"cep"`
}

type ReadUsuario struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}
