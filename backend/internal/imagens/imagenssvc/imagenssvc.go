package imagenssvc

import (
	"context"
	"escambo/internal/imagens/imagensrepo"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type ImagemRepository interface {
	SalvarImagem(dados imagensrepo.Metadata) error
	RemoverImagemDaPostagem(ctx context.Context, postagemID string, targetURL string) error
}

type Service struct {
	s3Client *s3.Client
	bucket   string
	repo     ImagemRepository
	region   string
}

func NewService(bucket string, region string, repo ImagemRepository) (*Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração AWS: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	return &Service{
		s3Client: s3Client,
		bucket:   bucket,
		repo:     repo,
		region:   region,
	}, nil
}

func (s *Service) UploadImagem(file multipart.File, header *multipart.FileHeader, ID, operacao string) (string, error) {
	key := operacao + "/" + uuid.New().String() + filepath.Ext(header.Filename)

	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		log.Printf("Erro ao enviar imagem para o S3 (ID: %s, operação: %s, nome do arquivo: %s): %v", ID, operacao, header.Filename, err)
		return "", fmt.Errorf("erro ao enviar para S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)

	metadata := imagensrepo.Metadata{
		ID:       ID,
		Operacao: operacao,
		URL:      url,
	}

	err = s.repo.SalvarImagem(metadata)
	if err != nil {
		log.Printf("Erro ao salvar metadados da imagem no banco (ID: %s, operação: %s, URL: %s): %v", ID, operacao, url, err)
		return "", fmt.Errorf("erro ao salvar no banco: %w", err)
	}

	return url, nil
}

func (s Service) DeletarImagemPostagem(ctx context.Context, postagemID, url string) error {
	err := s.deletarImagem(ctx, url)
	if err != nil {
		return fmt.Errorf("erro ao deletar imagem do S3: %w", err)
	}

	err = s.DeletarImagemPostagem(ctx, postagemID, url)
	if err != nil {
		return fmt.Errorf("erro ao deletar imagem da tabela de postagens: %w", err)
	}

	return nil
}

func (s Service) deletarImagem(ctx context.Context, url string) error {
	prefix := "https://escambo.s3.us-east-2.amazonaws.com/"
	if !strings.HasPrefix(url, prefix) {
		return fmt.Errorf("URL não é compatível com o padrão esperado")
	}

	key := strings.TrimPrefix(url, prefix)

	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar imagem do S3: %w", err)
	}

	return nil
}
