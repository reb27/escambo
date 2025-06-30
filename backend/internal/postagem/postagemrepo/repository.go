package postagemrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (r Repository) InsertPostagem(ctx context.Context, post Postagem) (string, error) {
	query := `
		INSERT INTO postagens (id, titulo, descricao, user_id, categoria)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT unique_titulo_user DO NOTHING
		RETURNING id;
	`

	var id string
	err := r.DB.QueryRowContext(ctx,
		query,
		post.Titulo,
		post.Descricao,
		post.UserID,
		post.Categoria,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return "", nil
	}

	return id, err
}

func (r Repository) GetPostagemByID(ctx context.Context, postID string, usuarioID string) (Postagem, error) {
	query := `
		SELECT 
			p.titulo, 
			p.descricao, 
			p.user_id, 
			p.categoria,
			p.created_at, 
			p.imagem_url,
			EXISTS (
				SELECT 1 FROM favoritos f 
				WHERE f.postagem_id = p.id AND f.usuario_id = $2
			) AS eh_favorito
		FROM postagens p
		WHERE p.id = $1;
	`

	var post Postagem
	var imagensJSON []byte

	err := r.DB.QueryRowContext(ctx, query, postID, usuarioID).Scan(
		&post.Titulo,
		&post.Descricao,
		&post.UserID,
		&post.Categoria,
		&post.CreatedAt,
		&imagensJSON,
		&post.MarcadoComoFavorito,
	)
	if err != nil {
		return Postagem{}, err
	}

	if len(imagensJSON) > 0 {
		if err := json.Unmarshal(imagensJSON, &post.Imagens); err != nil {
			return Postagem{}, fmt.Errorf("erro ao decodificar imagem_url: %w", err)
		}
	}

	return post, nil
}

func (r Repository) GetPostagens(ctx context.Context, filtro FiltroPostagem, offset int) ([]Postagem, error) {
	ordenacao := "DESC"
	if strings.ToUpper(filtro.Ordenacao) == "ASC" {
		ordenacao = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT 
			p.id,
			p.ativa,
			p.titulo, 
			p.descricao, 
			p.categoria, 
			p.imagem_url,
			p.user_id,
			p.created_at,
			u.nome AS nome_usuario,
			e.cidade,
			e.estado,
			e.bairro
		FROM postagens p
		JOIN usuarios u ON u.id = p.user_id
		JOIN endereco e ON e.user_id = u.id
		WHERE p.ativa = TRUE
			AND (NULLIF($1, '') IS NULL OR p.categoria = $1)
			AND (NULLIF($2, '') IS NULL OR p.user_id = $2::uuid)
			AND (
				NULLIF($5, '') IS NULL OR
				p.titulo ILIKE '%%' || $5 || '%%' OR
				p.descricao ILIKE '%%' || $5 || '%%'
			)
		ORDER BY p.created_at %s
		LIMIT $3 OFFSET $4
	`, ordenacao)

	rows, err := r.DB.QueryContext(ctx, query,
		filtro.Categoria,  // $1
		filtro.UsuarioID,  // $2
		filtro.Limite,     // $3
		offset,            // $4
		filtro.BuscaTexto, // $5
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var postagens []Postagem

	for rows.Next() {
		var post Postagem
		var imagensJSON *string

		if err := rows.Scan(
			&post.ID,
			&post.Status,
			&post.Titulo,
			&post.Descricao,
			&post.Categoria,
			&imagensJSON,
			&post.UserID,
			&post.CreatedAt,
			&post.NomeUsuario,
			&post.Cidade,
			&post.Estado,
			&post.Bairro,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear linha: %w", err)
		}

		if imagensJSON != nil && *imagensJSON != "" {
			if err := json.Unmarshal([]byte(*imagensJSON), &post.Imagens); err != nil {
				return nil, fmt.Errorf("erro ao decodificar imagem_url: %w", err)
			}
		}

		postagens = append(postagens, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no loop de leitura: %w", err)
	}

	return postagens, nil
}

func (r Repository) UpdatePostagem(ctx context.Context, postagem PostagemEdicao) error {
	queryCheck := `
		SELECT EXISTS (
			SELECT 1 FROM propostas
			WHERE postagem_id = $1 AND status = 'aceita'
		)
	`
	var hasAccepted bool
	err := r.DB.QueryRowContext(ctx, queryCheck, postagem.ID).Scan(&hasAccepted)
	if err != nil {
		return fmt.Errorf("erro ao verificar propostas: %w", err)
	}
	if hasAccepted {
		return fmt.Errorf("não é possível atualizar postagem pois já existe uma proposta aceita")
	}

	query := `UPDATE postagens SET `
	args := []interface{}{}
	updates := []string{}
	i := 1

	if postagem.Titulo != nil {
		updates = append(updates, fmt.Sprintf("titulo = $%d", i))
		args = append(args, *postagem.Titulo)
		i++
	}
	if postagem.Descricao != nil {
		updates = append(updates, fmt.Sprintf("descricao = $%d", i))
		args = append(args, *postagem.Descricao)
		i++
	}
	if postagem.Categoria != nil {
		updates = append(updates, fmt.Sprintf("categoria = $%d", i))
		args = append(args, *postagem.Categoria)
		i++
	}
	if postagem.Status != nil {
		updates = append(updates, fmt.Sprintf("ativa = $%d", i))
		args = append(args, *postagem.Status)
		i++
	}

	if len(updates) == 0 {
		return fmt.Errorf("nenhum campo foi informado para atualização")
	}

	updates = append(updates, "updated_at = NOW()")
	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, postagem.ID)

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar postagem: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r Repository) DeletarPostagem(ctx context.Context, postagemID string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM propostas
		WHERE postagem_id = $1 AND status != 'aceita'
	`, postagemID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao deletar propostas: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM postagens
		WHERE id = $1
	`, postagemID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao deletar postagem: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	return nil
}

func (r Repository) FavoritarPostagem(ctx context.Context, userID, postagemID string) error {
	query := `
		INSERT INTO favoritos (usuario_id, postagem_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := r.DB.ExecContext(ctx, query, userID, postagemID)
	if err != nil {
		return fmt.Errorf("erro ao favoritar postagem: %w", err)
	}

	return nil
}

func (r Repository) DesfavoritarPostagem(ctx context.Context, userID, postagemID string) error {
	query := `
		DELETE FROM favoritos WHERE usuario_id = $1 AND postagem_id = $2;
	`

	_, err := r.DB.ExecContext(ctx, query, userID, postagemID)
	if err != nil {
		return fmt.Errorf("erro ao desfavoritar postagem: %w", err)
	}

	return nil
}

func (r Repository) GetFavoritosByID(ctx context.Context, userID string) ([]Postagem, error) {
	query := `
		SELECT 
			p.id,
			p.ativa,
			p.titulo, 
			p.descricao, 
			p.categoria, 
			p.imagem_url,
			p.user_id,
			p.created_at,
			u.nome AS nome_usuario,
			e.cidade,
			e.estado,
			e.bairro
		FROM favoritos f
		JOIN postagens p ON p.id = f.postagem_id
		JOIN usuarios u ON u.id = p.user_id
		JOIN endereco e ON e.user_id = u.id
		WHERE f.usuario_id = $1;
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favoritos []Postagem

	for rows.Next() {
		var f Postagem
		var imagemURL []byte

		if err := rows.Scan(
			&f.ID,
			&f.Status,
			&f.Titulo,
			&f.Descricao,
			&f.Categoria,
			&imagemURL,
			&f.UserID,
			&f.CreatedAt,
			&f.NomeUsuario,
			&f.Cidade,
			&f.Estado,
			&f.Bairro,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(imagemURL, &f.Imagens); err != nil {
			return nil, err
		}

		favoritos = append(favoritos, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favoritos, nil
}
