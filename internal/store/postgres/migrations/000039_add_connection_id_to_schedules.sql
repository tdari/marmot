-- Add connection_id column to ingestion_schedules table
-- This migration decouples credential management from schedules by introducing
-- a foreign key reference to the connections table

ALTER TABLE ingestion_schedules
ADD COLUMN connection_id UUID REFERENCES connections(id) ON DELETE SET NULL;

-- Create index for performance
CREATE INDEX idx_ingestion_schedules_connection_id ON ingestion_schedules(connection_id);

COMMENT ON COLUMN ingestion_schedules.connection_id IS 'Reference to Connection entity for credential management (nullable for backward compatibility)';

-- Note: The config column is kept for backward compatibility during migration
-- It will be deprecated once all schedules have been migrated to use connections

---- create above / drop below ----

-- Drop index
DROP INDEX IF EXISTS idx_ingestion_schedules_connection_id;

-- Drop column
ALTER TABLE ingestion_schedules
DROP COLUMN IF EXISTS connection_id;
