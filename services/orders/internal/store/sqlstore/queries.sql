-- name: GetOrder :one
SELECT
    id,
    user_id,
    status,
    subtotal_amount,
    tax_amount,
    shipping_amount,
    discount_amount,
    total_amount,
    created_at,
    updated_at
FROM
    orders
WHERE
    id = $1;

-- name: InsertOrder :one
INSERT INTO orders (user_id, status, subtotal_amount, tax_amount, shipping_amount, discount_amount, total_amount)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id, user_id, status, subtotal_amount, tax_amount, shipping_amount, discount_amount, total_amount, created_at, updated_at;

-- name: GetOrderItem :one
SELECT
    id,
    order_id,
    product_id,
    variant_id,
    product_name,
    variant_name,
    sku,
    unit_price,
    quantity,
    total_price,
    created_at,
    updated_at
FROM
    order_items
WHERE
    id = $1;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, variant_id, product_name, variant_name, sku, unit_price, quantity, total_price)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
    id, order_id, product_id, variant_id, product_name, variant_name, sku, unit_price, quantity, total_price, created_at, updated_at;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

