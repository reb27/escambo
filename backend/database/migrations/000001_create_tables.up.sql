CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE usuarios (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  nome VARCHAR(150) NOT NULL,
  email VARCHAR(150) NOT NULL,
  senha VARCHAR(100) NOT NULL,
  telefone VARCHAR(20) NOT NULL,
  whatsapp_link VARCHAR(255), -- Alterado: DROP NOT NULL
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_email UNIQUE (email) -- Adicionado: restrição de e-mail único
);

CREATE TABLE postagens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  titulo VARCHAR(150) NOT NULL,
  descricao TEXT NOT NULL,
  ativa BOOLEAN DEFAULT true,
  imagem_url JSONB,
  user_id UUID NOT NULL,
  categoria VARCHAR(150) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_postagens_usuarios
    FOREIGN KEY (user_id) REFERENCES usuarios(id)
    ON DELETE CASCADE,
  CONSTRAINT unique_titulo_user UNIQUE (titulo, user_id) -- Adicionado: combinação única de título + usuário
);

CREATE TABLE endereco (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  cep VARCHAR(15) NOT NULL,
  rua VARCHAR(255) NOT NULL,
  numero INT NOT NULL,
  complemento TEXT,
  bairro VARCHAR(100),
  cidade VARCHAR(100),
  estado VARCHAR(2),
  user_id UUID NOT NULL UNIQUE, 
  CONSTRAINT fk_endereco_usuarios
    FOREIGN KEY (user_id) REFERENCES usuarios(id)
    ON DELETE CASCADE
);

CREATE TABLE propostas (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  postagem_id UUID NOT NULL,
  interessado_id UUID NOT NULL,      -- Usuário que faz a proposta
  dono_postagem_id UUID NOT NULL,    -- Usuário que recebe a proposta

  status VARCHAR(20) NOT NULL DEFAULT 'pendente'
    CHECK (status IN ('pendente', 'aceita', 'recusada')),

  imagem_url JSONB,
  descricao TEXT,

  nome VARCHAR(150) NOT NULL DEFAULT '',              -- Adicionado
  categoria VARCHAR(150) NOT NULL DEFAULT '',         -- Adicionado (ajustei o default para '', como no update)

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + interval '7 days'),

  CONSTRAINT fk_propostas_postagem
    FOREIGN KEY (postagem_id) REFERENCES postagens(id) ON DELETE CASCADE,
  CONSTRAINT fk_propostas_remetente
    FOREIGN KEY (interessado_id) REFERENCES usuarios(id) ON DELETE CASCADE,
  CONSTRAINT fk_propostas_destinatario
    FOREIGN KEY (dono_postagem_id) REFERENCES usuarios(id) ON DELETE CASCADE
);

CREATE TABLE notificacoes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  remetente_id UUID NOT NULL,        -- quem envia
  destinatario_id UUID NOT NULL,     -- quem recebe
  postagem_id UUID NOT NULL,         -- postagem
  proposta_status VARCHAR(20) NOT NULL
    CHECK (proposta_status IN ('pendente', 'aceita', 'recusada')),
  email_enviado BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
