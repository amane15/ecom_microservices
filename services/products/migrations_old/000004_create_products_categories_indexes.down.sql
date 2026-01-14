-- Products
DROP INDEX IF EXISTS idx_products_active;
DROP INDEX idx_products_slug;

-- Variants
DROP INDEX idx_variants_product_id;
DROP INDEX idx_variants_active;

-- Categories
DROP INDEX idx_categories_active;
