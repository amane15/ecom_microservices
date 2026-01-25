-- name: GetCart :one
SELECT id, user_id, created_at, updated_at
FROM carts
WHERE id = $1;

-- name: InsertCart :one
INSERT INTO carts(user_id)
VALUES ($1)
RETURNING id, user_id, created_at, updated_at;

-- name: DeleteCart :exec
DELETE FROM carts WHERE id = $1;

-- name: GetCartItem :one
SELECT id, cart_id, product_id, variant_id, quantity,
    created_at, updated_at
FROM cart_items
WHERE id = $1;

-- name: InsertCartItem :one
INSERT INTO cart_items(cart_id, product_id, variant_id, quantity)
VALUES ($1, $2,$3, $4)
ON CONFLICT (cart_id, variant_id)
DO UPDATE
SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = now()
RETURNING id, quantity, created_at, updated_at;

-- name: UpdateCartItem :one
INSERT INTO cart_items (cart_id, product_id, variant_id, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cart_id, variant_id)
DO UPDATE 
SET quantity = quantity + EXCLUDED.quantity, updated_at = now()
RETURNING id, cart_id, product_id, variant_id, quantity, created_at, updated_at;

-- name: DeleteCartItem :exec
DELETE FROM cart_items WHERE id = $1;
