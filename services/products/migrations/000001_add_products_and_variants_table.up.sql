CREATE TYPE product_status AS ENUM ('draft', 'active', 'archived');

CREATE TABLE products (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    name VARCHAR(255) NOT NULL,
    slug VARCHAR(128) NOT NULL UNIQUE,

    description TEXT,
    short_description TEXT,

    meta_title TEXT,
    meta_description TEXT,

    default_variant_id BIGINT,
    status product_status NOT NULL DEFAULT 'draft',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE products_variants (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    product_id BIGINT,

    slug VARCHAR(128) NOT NULL UNIQUE,
    name VARCHAR(255),

    price DECIMAL(6, 2),

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
