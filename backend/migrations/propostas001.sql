CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE propostas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    titulo VARCHAR(100) NOT NULL,
    nome_empresa VARCHAR(255) NOT NULL,
    nome_cliente VARCHAR(255) NOT NULL,
    prompt TEXT NOT NULL,
    cores TEXT[] NOT NULL,
    logo VARCHAR(255),
    logo_cliente VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    arquivo_final VARCHAR(255),
    data_criacao TIMESTAMPTZ NOT NULL,
    last_update TIMESTAMPTZ NOT NULL
);