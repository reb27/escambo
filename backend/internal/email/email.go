package email

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	ProposalUpdateTemplatePath   = "internal/email/templates/proposta_feita.html"
	ProposalResponseTemplatePath = "internal/email/templates/proposta_respondida.html"
)

type EmailData struct {
	UserName           string
	UserEmail          string
	OtherUserName      string
	ItemName           string
	ItemCategory       string
	ProposalDate       string
	ProposalStatus     string
	ProposalLink       string
	ItemImageURL       string
	HeaderMessage      string
	InterestedItemName string
}

func SendEmail(data EmailData, templatePath string, subject string) {
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("Erro ao ler o template: %v", err)
		return
	}

	tmpl, err := template.New("email").Parse(string(tmplContent))
	if err != nil {
		log.Printf("Erro ao processar template: %v", err)
		return
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		log.Printf("Erro ao preencher template: %v", err)
		return
	}

	from := mail.NewEmail("Escambo Online", "recepcaoescambo@gmail.com")
	to := mail.NewEmail(data.UserName, data.UserEmail)
	plainText := "Você recebeu uma atualização sobre uma proposta de troca na plataforma Escambo Online."

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlBody.String())
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println("Erro ao enviar:", err)
	} else {
		log.Printf("E-mail enviado para %s. Status: %d\n", data.UserEmail, response.StatusCode)
	}
}

type Notificacao struct {
	ID             string
	RemetenteID    string
	DestinatarioID string
	PostagemID     string
	Status         string
	CreatedAt      time.Time
}

func BuscarNotificacoesPendentes(db *sql.DB) (*Notificacao, error) {
	log.Println("Buscando notificações pendentes")

	var n Notificacao
	err := db.QueryRow(`
		SELECT id, remetente_id, destinatario_id, postagem_id, proposta_status, created_at
		FROM notificacoes
		WHERE email_enviado = FALSE
		LIMIT 1
	`).Scan(&n.ID, &n.RemetenteID, &n.DestinatarioID, &n.PostagemID, &n.Status, &n.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Nenhuma notificação pendente")
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func MarcarEmailComoEnviado(db *sql.DB, id string) error {
	_, err := db.Exec(`UPDATE notificacoes SET email_enviado = TRUE WHERE id = $1`, id)
	return err
}

func IniciarEnvioPeriodico(db *sql.DB, intervalo time.Duration) {
	log.Println("Início do processo de envio")

	ticker := time.NewTicker(intervalo)
	defer ticker.Stop()

	for range ticker.C {

		notif, err := BuscarNotificacoesPendentes(db)
		if err != nil {
			log.Printf("Erro ao buscar notificações: %v", err)
			continue
		}

		if notif == nil {
			continue
		}

		var (
			userName, userEmail     string
			otherUserName, itemName string
			itemCategory            string
			itemImageRaw            []byte
		)

		err = db.QueryRow(`
			SELECT u.nome, u.email, dono.nome, p.titulo, p.categoria, p.imagem_url
			FROM usuarios u
			JOIN usuarios dono ON dono.id = $1
			JOIN postagens p ON p.id = $2
			WHERE u.id = $3
		`, notif.RemetenteID, notif.PostagemID, notif.DestinatarioID).
			Scan(&userName, &userEmail, &otherUserName, &itemName, &itemCategory, &itemImageRaw)

		if err != nil {
			log.Printf("Erro ao buscar dados do usuário: %v", err)
			continue
		}

		imgURL := ""
		if len(itemImageRaw) > 0 {
			var urls []string
			if err := json.Unmarshal(itemImageRaw, &urls); err != nil {
				log.Printf("Erro ao parsear imagem JSON: %v", err)
			} else if len(urls) > 0 {
				imgURL = urls[0]
			}
		}

		emailData := EmailData{
			UserName:           userName,
			UserEmail:          userEmail,
			OtherUserName:      otherUserName,
			ItemName:           itemName,
			ItemCategory:       itemCategory,
			ProposalDate:       notif.CreatedAt.Format("02/01/2006"),
			ProposalStatus:     notif.Status,
			ProposalLink:       fmt.Sprintf("https://escambo.online/propostas/%s", notif.ID),
			ItemImageURL:       imgURL,
			HeaderMessage:      gerarHeader(notif.Status),
			InterestedItemName: "Item do interessado",
		}

		templatePath := ProposalUpdateTemplatePath
		if notif.Status == "aceita" || notif.Status == "recusada" {
			templatePath = ProposalResponseTemplatePath
		}

		SendEmail(emailData, templatePath, emailData.HeaderMessage)

		if err := MarcarEmailComoEnviado(db, notif.ID); err != nil {
			log.Printf("Erro ao marcar notificação %s como enviada: %v", notif.ID, err)
		}
	}
}

func gerarHeader(status string) string {
	switch status {
	case "pendente":
		return "Você recebeu uma nova proposta de troca!"
	case "aceita":
		return "Sua proposta foi aceita!"
	case "recusada":
		return "Sua proposta foi recusada!"
	default:
		return "Atualização sobre sua proposta!"
	}
}
