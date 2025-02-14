CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Включаем поддержку UUID

CREATE DOMAIN CODE as TEXT;

CREATE TABLE doc_merchs (
    code TEXT NOT NULL PRIMARY KEY,
    price INT NOT NULL CHECK (price > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE doc_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    balance INT NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sender_receiver_check CHECK (sender_id != receiver_id) -- Ограничение на то, чтобы сотрудник не отправлял монетки себе
);

-- Добавляем сортировку по умолчанию
CREATE INDEX idx_doc_merchs_created_at ON doc_merchs (created_at DESC);
CREATE INDEX idx_auth_users_created_at ON auth_users (created_at DESC);
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


-- Пользователи для нагрузочных тестов
INSERT INTO public.auth_users (id, login, "password", created_at) VALUES('516a477d-e9f4-4d05-8ee5-d0dba74182ff'::uuid, 'loadtester1', '$2a$10$y9xVCmNFG8XqpqqCKYG.wejWCwg4//YSYo.HLJJQxgrGxtJ6Hwd3G', '2025-02-14 17:40:58.683'); --loader
INSERT INTO public.auth_users (id, login, "password", created_at) VALUES('d79a6503-31c7-44e6-a831-91f892b8c88d'::uuid, 'loadtester2', '$2a$10$/RSh1B7TdvUNVai3caWuAuw1d91/Qoqph9qnhhf/nMMOfxB7YRbzu', '2025-02-14 17:41:03.674'); --loader

INSERT INTO public.doc_users (id, user_id, "name", balance, created_at) VALUES('2c745ab4-fb83-49b4-bbca-a4bf731e672f'::uuid, '516a477d-e9f4-4d05-8ee5-d0dba74182ff'::uuid, 'loadtester1', 2147483640, '2025-02-14 17:40:58.691');
INSERT INTO public.doc_users (id, user_id, "name", balance, created_at) VALUES('b31cfdd0-2678-4408-9a50-fcfe30d99cdc'::uuid, 'd79a6503-31c7-44e6-a831-91f892b8c88d'::uuid, 'loadtester2', 0, '2025-02-14 17:41:03.677');