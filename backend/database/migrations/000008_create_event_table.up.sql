CREATE TABLE notificacoes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    remetente_id UUID NOT NULL, -- quem envia
    destinatario_id UUID NOT NULL, -- quem recebe
    postagem_id UUID NOT NULL, -- postagem
    proposta_status VARCHAR(20) NOT NULL CHECK (proposta_status IN ('pendente', 'aceita', 'recusada')),
    email_enviado BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
