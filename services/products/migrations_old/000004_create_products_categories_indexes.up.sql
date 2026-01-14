-- Products
CREATE INDEX idx_products_active ON products(is_active);
CREATE INDEX idx_products_slug ON products(slug);

-- Variants
CREATE INDEX idx_variants_product_id ON products_variants(product_id);
CREATE INDEX idx_variants_active ON products_variants(is_active);

-- Categories
CREATE INDEX idx_categories_active ON categories(is_active);
