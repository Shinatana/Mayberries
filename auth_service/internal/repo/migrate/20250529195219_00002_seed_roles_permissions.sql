-- +goose Up

INSERT INTO roles (name, description) VALUES
('user', 'Standard user role'),
('manager', 'Manager with extended privileges'),
('admin', 'Administrator with full access');

INSERT INTO permissions (code, description) VALUES
                                                ('product.create', 'Permission to create products'),
                                                ('product.read', 'Permission to read products'),
                                                ('product.update', 'Permission to update products'),
                                                ('product.delete', 'Permission to delete products');


INSERT INTO roles_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r CROSS JOIN permissions p WHERE r.name = 'admin';

INSERT INTO roles_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'manager' AND p.code IN ('product.read', 'product.update');

INSERT INTO roles_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'user' AND p.code = 'product.read';

-- +goose Down

DELETE FROM roles_permissions;
DELETE FROM permissions;
DELETE FROM roles;