package postagemhandler

type PostagemResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type FavoritarPostagemRequest struct {
	UserID     string `json:"user_id"`
	PostagemID string `json:"postagem_id"`
}

type DesfavoritarPostagemRequest struct {
	UserID     string `json:"user_id"`
	PostagemID string `json:"postagem_id"`
}
