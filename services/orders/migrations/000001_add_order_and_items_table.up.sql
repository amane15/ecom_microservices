CREATE TYPE order_status AS ENUM (
    'pending', 
    'paid', 
    'cancelled', 
    'delivered', 
    'shipped', 
    'refunded'
);

CREATE TABLE orders (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,

    status order_status NOT NULL DEFAULT 'pending',

    subtotal_amount DECIMAL(10, 2) NOT NULL,
    tax_amount DECIMAL(10, 2) NOT NULL,
    shipping_amount DECIMAL(10, 2) NOT NULL,
    discount_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(10, 2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,

    product_id BIGINT NOT NULL,
    variant_id BIGINT NOT NULL,

    product_name TEXT NOT NULL,
    variant_name TEXT NOT NULL,
    sku VARCHAR(128) NOT NULL,

    unit_price DECIMAL(10, 2) NOT NULL,
    quantity BIGINT NOT NULL CHECK(quantity > 0),
    total_price DECIMAL(10, 2) NOT NULL GENERATED ALWAYS AS (unit_price * quantity) STORED,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

