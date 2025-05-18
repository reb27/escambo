# Escambo

# Backend
## ✅ Requisitos

- Go 1.24 ou superior  
- Docker e Docker Compose  
- `make` instalado  
- CLI do [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) instalada  
  > No Windows, recomenda-se instalar com [Scoop](https://scoop.sh/)

---

## 🚀 Iniciando o Projeto

### Passos para subir o ambiente

```bash
go mod tidy                # Instala as dependências Go
make run                   # Executa a aplicação
```
## 🧱 Migrations

### Criar nova migration

```bash
make migrate-create name=your_migration_name
```

### Rodar migrations
- Substitua DATABASE_URL pela URL do banco presente no .env.example

```bash
migrate -path=database/migrations -database DATABASE_URL -verbose up
```

Eli
