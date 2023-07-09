CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE document_t AS ENUM ('CPF', 'CNPJ');
CREATE TYPE account_t AS ENUM ('user', 'admin');

CREATE TABLE "account"(
    "id" BIGSERIAL PRIMARY KEY,
    "email" citext NOT NULL UNIQUE,
    "document_type" document_t NOT NULL,
    "document_number" VARCHAR NOT NULL,
    "password_hash" VARCHAR,
    "access_type" account_t DEFAULT('user'),
    "name" VARCHAR(150),
    "birth_date" date,
    UNIQUE (document_type,document_number)
)