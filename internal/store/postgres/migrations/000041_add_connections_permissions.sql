-- Add permissions for connections resource
INSERT INTO permissions (name, description, resource_type, action) VALUES
('view_connections', 'View connections', 'connections', 'view'),
('manage_connections', 'Create/update/delete connections', 'connections', 'manage');

-- Grant all permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'admin'),
    id
FROM permissions
WHERE name IN ('view_connections', 'manage_connections');

-- Grant view permission to user role
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'user'),
    id
FROM permissions
WHERE name = 'view_connections';

---- create above / drop below ----

-- Remove permissions
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE name IN ('view_connections', 'manage_connections')
);

DELETE FROM permissions WHERE name IN ('view_connections', 'manage_connections');
