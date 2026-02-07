-- Mini App Config API Database Schema
-- Version: 1.0.0

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- PAGES TABLE
-- ============================================
-- Represents a screen/page in the mobile app (e.g., Home, Collection, Product Detail)
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    route VARCHAR(255) NOT NULL UNIQUE,
    is_home BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster route lookups
CREATE INDEX IF NOT EXISTS idx_pages_route ON pages(route);

-- Index for home page lookup
CREATE INDEX IF NOT EXISTS idx_pages_is_home ON pages(is_home) WHERE is_home = TRUE;

-- ============================================
-- WIDGETS TABLE
-- ============================================
-- Represents UI components on a page (e.g., Banner, ProductGrid, Text)
CREATE TABLE IF NOT EXISTS widgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('banner', 'product_grid', 'text', 'image', 'spacer')),
    position INTEGER NOT NULL DEFAULT 0,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster page_id lookups
CREATE INDEX IF NOT EXISTS idx_widgets_page_id ON widgets(page_id);

-- Index for position ordering
CREATE INDEX IF NOT EXISTS idx_widgets_position ON widgets(page_id, position);

-- Index for type filtering
CREATE INDEX IF NOT EXISTS idx_widgets_type ON widgets(type);

-- ============================================
-- TRIGGER FUNCTION FOR updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for pages table
DROP TRIGGER IF EXISTS update_pages_updated_at ON pages;
CREATE TRIGGER update_pages_updated_at
    BEFORE UPDATE ON pages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for widgets table
DROP TRIGGER IF EXISTS update_widgets_updated_at ON widgets;
CREATE TRIGGER update_widgets_updated_at
    BEFORE UPDATE ON widgets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- SAMPLE DATA (Optional - for testing)
-- ============================================
-- Uncomment the following to insert sample data:

-- INSERT INTO pages (name, route, is_home) VALUES 
--     ('Home', '/home', TRUE),
--     ('Collections', '/collection', FALSE),
--     ('Product Detail', '/product/:id', FALSE);

-- INSERT INTO widgets (page_id, type, position, config) VALUES 
--     ((SELECT id FROM pages WHERE route = '/home'), 'banner', 1, '{"title": "Welcome Banner", "image_url": "https://example.com/banner.jpg"}'),
--     ((SELECT id FROM pages WHERE route = '/home'), 'product_grid', 2, '{"columns": 2, "limit": 10}'),
--     ((SELECT id FROM pages WHERE route = '/home'), 'text', 3, '{"content": "Featured Products", "style": "heading"}');
