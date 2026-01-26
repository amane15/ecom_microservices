-- name: GetProduct :one
SELECT
    id,
    name,
    slug,
    description,
    short_description,
    meta_title,
    meta_description,
    status,
    default_variant_id,
    created_at,
    updated_at,
    deleted_at
FROM
    products
WHERE
    id = $1;

-- name: InsertProduct :one
INSERT INTO products (name, slug, description, short_description, meta_title, meta_description, status, default_variant_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    id, created_at, updated_at;

-- name: ListProducts :many
SELECT
    id,
    name,
    slug,
    description,
    short_description,
    meta_title,
    meta_description,
    status,
    default_variant_id,
    created_at,
    updated_at,
    deleted_at
FROM
    products
WHERE
    status = 'active'
ORDER BY
    id
LIMIT $1 OFFSET $2;

SELECT
    id,
    name,
    slug,
    description,
    short_description,
    meta_title,
    meta_description,
    status,
    default_variant_id,
    created_at,
    updated_at,
    deleted_at
FROM
    products
WHERE
    status = 'active'
ORDER BY
    id
LIMIT $1 OFFSET $2;

-- name: GetVariant :one
SELECT
    id,
    product_id,
    slug,
    name,
    price,
    is_active,
    created_at,
    updated_at,
    deleted_at
FROM
    products_variants
WHERE
    id = $1;

-- name: InsertVariant :one
INSERT INTO products_variants (product_id, slug, name, price, is_active)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    id, created_at, updated_at;

-- name: GetVariantByProduct :one
SELECT
    count(*)
FROM
    products_variants
WHERE
    product_id = $1;

-- name: DeleteVariant :exec
DELETE FROM products_variants
WHERE id = $1;

-- name: GetCategory :one
SELECT
    id,
    name,
    slug,
    description,
    is_active,
    created_at,
    updated_at,
    deleted_at
FROM
    categories
WHERE
    id = $1;

-- name: InsertCategory :one
INSERT INTO categories (name, slug, description, is_active)
    VALUES ($1, $2, $3, $4)
RETURNING
    id, created_at, updated_at;

-- name: ListCategories :many
SELECT
    id,
    name,
    slug,
    description,
    is_active,
    created_at,
    updated_at,
    deleted_at
FROM
    categories
WHERE
    is_active = TRUE
ORDER BY
    id
LIMIT $1 OFFSET $2;

