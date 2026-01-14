CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_variants_slug ON products_variants(slug);
CREATE INDEX idx_products_variants_is_active ON products_variants(is_active);

