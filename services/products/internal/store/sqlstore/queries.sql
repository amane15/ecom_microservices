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
INSERT INTO products (name, slug, description, short_description, meta_title, meta_description, status)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id, name, slug, description, short_description, meta_title, meta_description, status, default_variant_id, created_at, updated_at;

-- name: UpdateProduct :one
UPDATE
    products
SET
    name = COALESCE(sqlc.narg (name), name),
    description = COALESCE(sqlc.narg (description), description),
    short_description = COALESCE(sqlc.narg (short_description), short_description),
    meta_title = COALESCE(sqlc.narg (meta_title), meta_title),
    meta_description = COALESCE(sqlc.narg (meta_description), meta_description),
    status = COALESCE(sqlc.narg (status), status),
    default_variant_id = COALESCE(sqlc.narg (default_variant_id), default_variant_id),
    updated_at = now()
WHERE
    id = $1
RETURNING
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
    updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

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
    id, product_id, slug, name, price, is_active, created_at, updated_at, deleted_at;

-- name: UpdateVariant :one
UPDATE
    products_variants
SET
    name = COALESCE(sqlc.narg (name), name),
    price = COALESCE(sqlc.narg (price), price),
    is_active = COALESCE(sqlc.narg (is_active), is_active)
WHERE
    id = $1
RETURNING
    id,
    product_id,
    slug,
    name,
    price,
    is_active,
    created_at,
    updated_at,
    deleted_at;

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

-- name: ListProductVariants :many
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
    product_id = $1;

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
    id, name, slug, description, is_active, created_at, updated_at;

-- name: UpdateCategory :one
UPDATE
    categories
SET
    name = COALESCE(sqlc.narg (name), name),
    description = COALESCE(sqlc.narg (description), description),
    is_active = COALESCE(sqlc.narg (is_active), is_active)
WHERE
    id = $1
RETURNING
    id,
    name,
    slug,
    description,
    is_active,
    created_at,
    updated_at,
    deleted_at;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

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

