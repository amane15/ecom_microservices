CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_role AS ENUM (
    'user',
    'admin'
);

CREATE TABLE users (
    id bigserial PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    first_name text NOT NULL,
    last_name text NOT NULL,
    ROLE user_role NOT NULL DEFAULT 'user',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

