-- Data migration: Move schedule configs to Connection entities
-- 
-- This migration extracts credentials from schedule configs and creates Connection entities.
-- Encrypted values are copied as-is (no decryption/re-encryption needed since the same
-- encryption key is used). The JSON structure is transformed to match connection type expectations.

-- Create connections from existing schedules
INSERT INTO connections (id, name, type, description, config, tags, created_by, created_at, updated_at)
SELECT 
    gen_random_uuid() AS id,
    -- Generate unique connection name from schedule name
    LOWER(REGEXP_REPLACE(
        CASE 
            WHEN TRIM(s.name) = '' THEN 'schedule'
            ELSE SUBSTRING(TRIM(s.name), 1, 180)
        END || '-connection',
        '\s+', '-', 'g'
    )) AS name,
    CASE
        WHEN s.plugin_id IN ('s3', 'dynamodb', 'sns', 'sqs', 'glue', 'lambda', 'kinesis') THEN 'aws'
        WHEN s.plugin_id IN ('redpanda', 'confluent') THEN 'kafka'
        WHEN s.plugin_id = 'iceberg' THEN 'iceberg-rest'
        ELSE s.plugin_id
    END AS type,
    'Migrated from schedule: ' || s.name AS description,
    -- Transform config structure based on plugin type
    CASE 
        -- AWS services: extract flat AWSConfig fields (handle both nested credentials object and flat fields)
        WHEN s.plugin_id IN ('s3', 'dynamodb', 'sns', 'sqs', 'glue', 'lambda', 'kinesis') THEN
            jsonb_strip_nulls(jsonb_build_object(
                'region',          COALESCE(s.config->'credentials'->>'region',         s.config->>'region'),
                'id',              COALESCE(s.config->'credentials'->>'id',             s.config->>'access_key_id'),
                'secret',          COALESCE(s.config->'credentials'->>'secret',         s.config->>'secret_access_key'),
                'token',           COALESCE(s.config->'credentials'->>'token',          s.config->>'token'),
                'use_default',     COALESCE(s.config->'credentials'->'use_default',     s.config->'use_iam_role', s.config->'use_default'),
                'profile',         COALESCE(s.config->'credentials'->>'profile',        s.config->>'profile'),
                'role',            COALESCE(s.config->'credentials'->>'role',           s.config->>'role'),
                'role_external_id',COALESCE(s.config->'credentials'->>'role_external_id',s.config->>'role_external_id'),
                'endpoint',        COALESCE(s.config->'credentials'->>'endpoint',       s.config->>'endpoint')
            ))
        
        -- PostgreSQL: extract connection fields only
        WHEN s.plugin_id = 'postgresql' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'database', s.config->>'database',
                'user', s.config->>'user',
                'password', s.config->>'password',
                'ssl_mode', s.config->>'ssl_mode',
                'sslmode', s.config->>'sslmode'
            ))
        
        -- MySQL: extract connection fields only
        WHEN s.plugin_id = 'mysql' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'database', s.config->>'database',
                'user', s.config->>'user',
                'password', s.config->>'password',
                'tls', s.config->>'tls'
            ))
        
        -- BigQuery: extract connection fields only
        WHEN s.plugin_id = 'bigquery' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'project_id', s.config->>'project_id',
                'dataset_id', s.config->>'dataset_id',
                'credentials_path', s.config->>'credentials_path',
                'use_default_credentials', s.config->'use_default_credentials',
                'location', s.config->>'location'
            ))
        
        -- Snowflake: extract connection fields only
        WHEN s.plugin_id = 'snowflake' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'account', s.config->>'account',
                'user', s.config->>'user',
                'password', s.config->>'password',
                'warehouse', s.config->>'warehouse',
                'database', s.config->>'database',
                'schema', s.config->>'schema',
                'role', s.config->>'role'
            ))
        
        -- Redshift: extract connection fields only
        WHEN s.plugin_id = 'redshift' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'database', s.config->>'database',
                'user', s.config->>'user',
                'password', s.config->>'password',
                'ssl_mode', s.config->>'ssl_mode'
            ))
        
        -- Databricks: extract connection fields only
        WHEN s.plugin_id = 'databricks' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'http_path', s.config->>'http_path',
                'access_token', s.config->>'access_token',
                'catalog', s.config->>'catalog',
                'schema', s.config->>'schema'
            ))
        
        -- Kafka: extract connection fields only
        WHEN s.plugin_id = 'kafka' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'bootstrap_servers', s.config->>'bootstrap_servers',
                'client_id', s.config->>'client_id',
                'authentication', s.config->'authentication',
                'tls', s.config->'tls',
                'schema_registry', s.config->'schema_registry'
            ))
        
        -- MongoDB: extract connection fields only
        WHEN s.plugin_id = 'mongodb' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'database', s.config->>'database',
                'user', s.config->>'user',
                'password', s.config->>'password',
                'auth_source', s.config->>'auth_source',
                'replica_set', s.config->>'replica_set',
                'tls', s.config->'tls'
            ))
        
        -- ClickHouse: extract connection fields only
        WHEN s.plugin_id = 'clickhouse' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'database', s.config->>'database',
                'user', s.config->>'user',
                'password', s.config->>'password',
                'secure', s.config->'secure'
            ))
        
        -- GCS: extract connection fields only
        WHEN s.plugin_id = 'gcs' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'project_id', s.config->>'project_id',
                'bucket_name', s.config->>'bucket_name',
                'credentials_file', s.config->>'credentials_file',
                'credentials_json', s.config->>'credentials_json',
                'disable_auth', s.config->'disable_auth'
            ))
        
        -- Azure Blob: extract connection fields only
        WHEN s.plugin_id = 'azureblob' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'connection_string', s.config->>'connection_string',
                'account_name', s.config->>'account_name',
                'account_key', s.config->>'account_key',
                'container_name', s.config->>'container_name',
                'endpoint', s.config->>'endpoint'
            ))

        -- Trino: extract connection fields only
        WHEN s.plugin_id = 'trino' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'user', s.config->>'user',
                'password', s.config->>'password',
                'secure', s.config->'secure',
                'ssl_cert_path', s.config->>'ssl_cert_path',
                'access_token', s.config->>'access_token'
            ))

        -- Iceberg REST catalog: extract connection fields only
        WHEN s.plugin_id = 'iceberg' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'uri', s.config->>'uri',
                'warehouse', s.config->>'warehouse',
                'credential', s.config->>'credential',
                'token', s.config->>'token',
                'properties', s.config->'properties',
                'prefix', s.config->>'prefix'
            ))

        -- Redpanda / Confluent: reuse kafka connection fields
        WHEN s.plugin_id IN ('redpanda', 'confluent') THEN
            jsonb_strip_nulls(jsonb_build_object(
                'bootstrap_servers', s.config->>'bootstrap_servers',
                'client_id', s.config->>'client_id',
                'authentication', s.config->'authentication',
                'tls', s.config->'tls'
            ))

        -- NATS: extract connection fields only
        WHEN s.plugin_id = 'nats' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'token', s.config->>'token',
                'username', s.config->>'username',
                'password', s.config->>'password',
                'credentials_file', s.config->>'credentials_file',
                'tls', s.config->'tls',
                'tls_insecure', s.config->'tls_insecure'
            ))

        -- Airflow: extract connection fields only
        WHEN s.plugin_id = 'airflow' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'username', s.config->>'username',
                'password', s.config->>'password',
                'api_token', s.config->>'api_token'
            ))

        -- Redis: extract connection fields only
        WHEN s.plugin_id = 'redis' THEN
            jsonb_strip_nulls(jsonb_build_object(
                'host', s.config->>'host',
                'port', s.config->'port',
                'password', s.config->>'password',
                'username', s.config->>'username',
                'db', s.config->'db',
                'tls', s.config->'tls',
                'tls_insecure', s.config->'tls_insecure'
            ))

        -- Fallback: empty config if plugin type not recognized
        ELSE '{}'::jsonb
    END AS config,
    ARRAY['migrated', 'schedule']::TEXT[] AS tags,
    COALESCE(s.created_by, 'system-migration') AS created_by,
    s.created_at AS created_at,
    now() AS updated_at
FROM ingestion_schedules s
WHERE s.connection_id IS NULL
  AND s.plugin_id IN (
      'postgresql', 'mysql', 's3', 'bigquery', 'snowflake', 'redshift',
      'databricks', 'dynamodb', 'sns', 'sqs', 'glue', 'lambda', 'kinesis',
      'kafka', 'gcs', 'mongodb', 'clickhouse', 'azureblob',
      'trino', 'iceberg', 'nats', 'airflow', 'redis', 'redpanda', 'confluent'
  );

-- Update schedules to reference newly created connections
-- Match by name pattern and plugin type
UPDATE ingestion_schedules s
SET connection_id = c.id,
    updated_at = now()
FROM connections c
WHERE s.connection_id IS NULL
  AND c.name = LOWER(REGEXP_REPLACE(
      CASE 
          WHEN TRIM(s.name) = '' THEN 'schedule'
          ELSE SUBSTRING(TRIM(s.name), 1, 180)
      END || '-connection',
      '\s+', '-', 'g'
  ))
  AND c.type = CASE
      WHEN s.plugin_id IN ('s3', 'dynamodb', 'sns', 'sqs', 'glue', 'lambda', 'kinesis') THEN 'aws'
      WHEN s.plugin_id IN ('redpanda', 'confluent') THEN 'kafka'
      WHEN s.plugin_id = 'iceberg' THEN 'iceberg-rest'
      ELSE s.plugin_id
  END
  AND c.tags @> ARRAY['migrated', 'schedule']::TEXT[];

-- Clean up schedule configs: remove connection fields, keep only discovery fields
UPDATE ingestion_schedules s
SET config = CASE 
    -- AWS services: remove credentials object and all flat credential fields
    WHEN s.plugin_id IN ('s3', 'dynamodb', 'sns', 'sqs', 'glue', 'lambda', 'kinesis') THEN
        s.config - ARRAY[
            'credentials',
            'region', 'access_key_id', 'secret_access_key', 'use_iam_role', 'use_default',
            'profile', 'role', 'role_external_id', 'endpoint', 'token'
        ]
    
    -- PostgreSQL: remove connection fields
    WHEN s.plugin_id = 'postgresql' THEN
        s.config - ARRAY['host', 'port', 'database', 'user', 'password', 'ssl_mode', 'sslmode']
    
    -- MySQL: remove connection fields
    WHEN s.plugin_id = 'mysql' THEN
        s.config - ARRAY['host', 'port', 'database', 'user', 'password', 'tls']
    
    -- BigQuery: remove connection fields
    WHEN s.plugin_id = 'bigquery' THEN
        s.config - ARRAY['project_id', 'dataset_id', 'credentials_path', 'use_default_credentials', 'location']
    
    -- Snowflake: remove connection fields
    WHEN s.plugin_id = 'snowflake' THEN
        s.config - ARRAY['account', 'user', 'password', 'warehouse', 'database', 'schema', 'role']
    
    -- Redshift: remove connection fields
    WHEN s.plugin_id = 'redshift' THEN
        s.config - ARRAY['host', 'port', 'database', 'user', 'password', 'ssl_mode']
    
    -- Databricks: remove connection fields
    WHEN s.plugin_id = 'databricks' THEN
        s.config - ARRAY['host', 'http_path', 'access_token', 'catalog', 'schema']
    
    -- Kafka: remove connection fields
    WHEN s.plugin_id = 'kafka' THEN
        s.config - ARRAY['bootstrap_servers', 'client_id', 'authentication', 'tls', 'schema_registry']
    
    -- MongoDB: remove connection fields
    WHEN s.plugin_id = 'mongodb' THEN
        s.config - ARRAY['host', 'port', 'database', 'user', 'password', 'auth_source', 'replica_set', 'tls']
    
    -- ClickHouse: remove connection fields
    WHEN s.plugin_id = 'clickhouse' THEN
        s.config - ARRAY['host', 'port', 'database', 'user', 'password', 'secure']
    
    -- GCS: remove connection fields
    WHEN s.plugin_id = 'gcs' THEN
        s.config - ARRAY['project_id', 'bucket_name', 'credentials_file', 'credentials_json', 'disable_auth']
    
    -- Azure Blob: remove connection fields
    WHEN s.plugin_id = 'azureblob' THEN
        s.config - ARRAY['connection_string', 'account_name', 'account_key', 'container_name', 'endpoint']

    -- Trino: remove connection fields
    WHEN s.plugin_id = 'trino' THEN
        s.config - ARRAY['host', 'port', 'user', 'password', 'secure', 'ssl_cert_path', 'access_token']

    -- Iceberg: remove connection fields
    WHEN s.plugin_id = 'iceberg' THEN
        s.config - ARRAY['uri', 'warehouse', 'credential', 'token', 'properties', 'prefix']

    -- Redpanda: remove kafka connection fields
    WHEN s.plugin_id = 'redpanda' THEN
        s.config - ARRAY['bootstrap_servers', 'client_id', 'authentication', 'tls']

    -- Confluent: remove kafka connection fields
    WHEN s.plugin_id = 'confluent' THEN
        s.config - ARRAY['bootstrap_servers', 'client_id', 'authentication']

    -- NATS: remove connection fields
    WHEN s.plugin_id = 'nats' THEN
        s.config - ARRAY['host', 'port', 'token', 'username', 'password', 'credentials_file', 'tls', 'tls_insecure']

    -- Airflow: remove connection fields
    WHEN s.plugin_id = 'airflow' THEN
        s.config - ARRAY['host', 'username', 'password', 'api_token']

    -- Redis: remove connection fields
    WHEN s.plugin_id = 'redis' THEN
        s.config - ARRAY['host', 'port', 'password', 'username', 'db', 'tls', 'tls_insecure']

    ELSE s.config
END,
updated_at = now()
WHERE connection_id IS NOT NULL
  AND EXISTS (
      SELECT 1 FROM connections c 
      WHERE c.id = s.connection_id 
      AND c.tags @> ARRAY['migrated', 'schedule']::TEXT[]
  );

---- create above / drop below ----

-- Rollback: Remove connection references and delete migrated connections
UPDATE ingestion_schedules 
SET connection_id = NULL, updated_at = now()
WHERE connection_id IN (
    SELECT id FROM connections WHERE tags @> ARRAY['migrated', 'schedule']::TEXT[]
);

DELETE FROM connections WHERE tags @> ARRAY['migrated', 'schedule']::TEXT[];
