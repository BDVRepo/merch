CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Включаем поддержку UUID

CREATE DOMAIN CODE as TEXT;

CREATE TABLE doc_merchs (
    code TEXT NOT NULL PRIMARY KEY,
    price INT NOT NULL CHECK (price > 0),
    created_at TIMESTAMP DEFAULT now()
);

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE doc_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    balance INT NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE doc_user_merchs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    root_id UUID NOT NULL REFERENCES doc_users(id) ON DELETE CASCADE,
    merch_code CODE NOT NULL REFERENCES doc_merchs(code)
);

CREATE TABLE doc_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL REFERENCES doc_users(id),
    receiver_id UUID REFERENCES doc_users(id),
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT sender_receiver_check CHECK (sender_id != receiver_id) -- Ограничение на то, чтобы сотрудник не отправлял монетки себе
);

-- Добавляем сортировку по умолчанию
CREATE INDEX idx_doc_merchs_created_at ON doc_merchs (created_at DESC);
CREATE INDEX idx_auth_users_created_at ON auth.users (created_at DESC);
CREATE INDEX idx_doc_users_created_at ON doc_users (created_at DESC);
CREATE INDEX idx_doc_transactions_created_at ON doc_transactions (created_at DESC);
CREATE INDEX idx_doc_transactions_sender_receiver ON doc_transactions (sender_id, receiver_id);

-- Вставка товаров
INSERT INTO doc_merchs (code, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);
