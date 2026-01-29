CREATE TYPE product_status AS ENUM (
    'draft',
    'active',
    'archived'
);

CREATE TABLE products (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(255) NOT NULL,
    slug varchar(128) NOT NULL UNIQUE,
    description text,
    short_description text,
    meta_title text,
    meta_description text,
    default_variant_id bigint,
    status product_status NOT NULL DEFAULT 'draft',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE products_variants (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    product_id bigint NOT NULL REFERENCES products (id) ON DELETE RESTRICT,
    slug varchar(128) NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    price DECIMAL(6, 2) NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

