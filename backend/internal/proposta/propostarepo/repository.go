package propostarepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

func (r Repository) UpdatePropostaStatus(ctx context.Context, propostaID, status string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	updateQuery := `
		UPDATE propostas
		SET status = $1
		WHERE id = $2;
	`
	if _, err := tx.ExecContext(ctx, updateQuery, status, propostaID); err != nil {
		tx.Rollback()
		return err
	}

	var remetenteID, destinatarioID, postagemID string
	selectQuery := `
		SELECT interessado_id, dono_postagem_id, postagem_id
		FROM propostas
		WHERE id = $1;
	`
	if err := tx.QueryRowContext(ctx, selectQuery, propostaID).Scan(&remetenteID, &destinatarioID, &postagemID); err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao buscar dados da proposta: %w", err)
	}

	insertNotificacao := `
		INSERT INTO notificacoes (remetente_id, destinatario_id, postagem_id, proposta_status)
		VALUES ($1, $2, $3, $4);
	`
	if _, err := tx.ExecContext(ctx, insertNotificacao, remetenteID, destinatarioID, postagemID, status); err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir notificação: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

func (r Repository) InsertProposta(ctx context.Context, proposta PropostaWriteModel) (*Proposta, error) {
	checkQuery := `
		SELECT 1 FROM propostas
		WHERE postagem_id = $1 AND interessado_id = $2 AND nome = $3
	`
	var exists int
	err := r.DB.QueryRowContext(ctx, checkQuery,
		proposta.PostagemID,
		proposta.RemetenteID,
		proposta.Nome,
	).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("erro ao verificar proposta existente: %w", err)
	}

	if err == nil {
		return nil, fmt.Errorf("uma proposta com esse título já foi enviada para essa postagem")
	}

	insertQuery := `
		INSERT INTO propostas (
			postagem_id, interessado_id, dono_postagem_id, descricao, nome, categoria
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, postagem_id, interessado_id, dono_postagem_id, descricao, nome, categoria, status, created_at;
	`

	var inserted Proposta
	err = r.DB.QueryRowContext(ctx, insertQuery,
		proposta.PostagemID,
		proposta.RemetenteID,
		proposta.DestinatarioID,
		proposta.Descricao,
		proposta.Nome,
		proposta.Categoria,
	).Scan(
		&inserted.ID,
		&inserted.PostagemID,
		&inserted.InteressadoID,
		&inserted.DonoPostagemID,
		&inserted.Descricao,
		&inserted.Nome,
		&inserted.Categoria,
		&inserted.Status,
		&inserted.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar proposta: %w", err)
	}

	evento := `
		INSERT INTO notificacoes (
			remetente_id, destinatario_id, postagem_id, proposta_status
		) VALUES ($1, $2, $3, $4);
	`

	_, err = r.DB.ExecContext(ctx, evento, proposta.RemetenteID, proposta.DestinatarioID, proposta.PostagemID, "pendente")
	if err != nil {
		// Loga o erro, mas não falha a operação principal
		log.Printf("erro ao salvar evento na tabela de notificacoes: %v", err)
	}

	return &inserted, nil
}

func (r Repository) GetPropostas(ctx context.Context, filter PropostasQueryFilter) ([]PropostaFormatada, error) {
	query := `
		SELECT 
			p.status,

			po.titulo AS postagem_titulo,
			po.descricao AS postagem_descricao,
			po.categoria AS postagem_categoria,
			uo.nome AS postagem_usuario,
			po.imagem_url AS postagem_imagens,

			p.nome AS proposta_titulo,
			p.categoria AS proposta_categoria,
			p.descricao AS proposta_descricao,
			ui.nome AS proposta_usuario,
			p.imagem_url AS proposta_imagens

		FROM propostas p
		JOIN postagens po ON p.postagem_id = po.id
		JOIN usuarios uo ON po.user_id = uo.id
		JOIN usuarios ui ON p.interessado_id = ui.id
	`

	filteredQuery, args := buildQuery(query, filter)
	rows, err := r.DB.QueryContext(ctx, filteredQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var propostas []PropostaFormatada
	for rows.Next() {
		var p PropostaFormatada
		var imagensPostagemJSON, imagensPropostaJSON []byte

		err := rows.Scan(
			&p.Status,

			&p.ProdutoPostagem.Nome,
			&p.ProdutoPostagem.Descricao,
			&p.ProdutoPostagem.Categoria,
			&p.ProdutoPostagem.Usuario,
			&imagensPostagemJSON,

			&p.ProdutoPropostaTroca.Nome,
			&p.ProdutoPropostaTroca.Categoria,
			&p.ProdutoPropostaTroca.Descricao,
			&p.ProdutoPropostaTroca.Usuario,
			&imagensPropostaJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(imagensPostagemJSON) > 0 && string(imagensPostagemJSON) != "null" {
			if err := json.Unmarshal(imagensPostagemJSON, &p.ProdutoPostagem.Imagens); err != nil {
				return nil, fmt.Errorf("erro ao decodificar imagem_url da postagem: %w", err)
			}
		}

		if len(imagensPropostaJSON) > 0 && string(imagensPropostaJSON) != "null" {
			if err := json.Unmarshal(imagensPropostaJSON, &p.ProdutoPropostaTroca.Imagens); err != nil {
				return nil, fmt.Errorf("erro ao decodificar imagem_url da proposta: %w", err)
			}
		}

		propostas = append(propostas, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(propostas) == 0 {
		return nil, fmt.Errorf("nenhum dado encontrado")
	}

	return propostas, nil
}

func buildQuery(baseSQL string, filter PropostasQueryFilter) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.Tipo == "enviadas" {
		conditions = append(conditions, fmt.Sprintf("interessado_id = $%d", argIndex))
		args = append(args, filter.UsuarioID)
		argIndex++
	} else if filter.Tipo == "recebidas" {
		conditions = append(conditions, fmt.Sprintf("dono_postagem_id = $%d", argIndex))
		args = append(args, filter.UsuarioID)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if len(conditions) > 0 {
		baseSQL += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseSQL += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	return baseSQL, args
}
