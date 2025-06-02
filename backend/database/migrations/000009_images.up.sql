CREATE TABLE post_images (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  operacao_id UUID,
  operacao VARCHAR(30),
  url TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  CHECK (operacao IN ('troca', 'postagem'))
);
