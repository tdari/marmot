-- Create connections table for managing data source credentials
CREATE TABLE connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    description TEXT,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    tags TEXT[] DEFAULT ARRAY[]::TEXT[],
    created_by VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT unique_connection_name UNIQUE(name)
);

COMMENT ON COLUMN connections.config IS 'Configuration JSONB where each field can have {value, is_sensitive} structure. Sensitive values are encrypted.';

-- Indexes for common queries
CREATE INDEX idx_connections_type ON connections(type);
CREATE INDEX idx_connections_created_by ON connections(created_by);

-- Trigger to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_connections_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER connections_updated_at
    BEFORE UPDATE ON connections
    FOR EACH ROW
    EXECUTE FUNCTION update_connections_updated_at();

---- create above / drop below ----

-- Drop trigger and function
DROP TRIGGER IF EXISTS connections_updated_at ON connections;
DROP FUNCTION IF EXISTS update_connections_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_connections_created_by;
DROP INDEX IF EXISTS idx_connections_type;

-- Drop table
DROP TABLE IF EXISTS connections;
