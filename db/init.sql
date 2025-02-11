CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Включаем поддержку UUID

CREATE DOMAIN CODE AS VARCHAR(50);

CREATE TABLE doc_merchs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
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

CREATE TABLE info_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    balance INT NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE doc_merch_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merch_code UUID NOT NULL REFERENCES doc_merchs(id),
    receiver_id UUID NOT NULL REFERENCES info_users(id),
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT now()
);

-- Добавляем сортировку по умолчанию
CREATE INDEX idx_doc_merchs_created_at ON doc_merchs (created_at DESC);
CREATE INDEX idx_auth_users_created_at ON auth.users (created_at DESC);
CREATE INDEX idx_info_users_created_at ON info_users (created_at DESC);
CREATE INDEX idx_doc_merch_transactions_created_at ON doc_merch_transactions (created_at DESC);

-- Вставка товаров
INSERT INTO doc_merchs (name, price) VALUES
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