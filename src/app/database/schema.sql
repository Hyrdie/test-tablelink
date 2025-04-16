-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create role_rights table
CREATE TABLE IF NOT EXISTS role_rights (
    id SERIAL PRIMARY KEY,
    role_id INTEGER REFERENCES roles(id),
    section VARCHAR(50) NOT NULL,
    route VARCHAR(255) NOT NULL,
    r_create BOOLEAN DEFAULT FALSE,
    r_read BOOLEAN DEFAULT FALSE,
    r_update BOOLEAN DEFAULT FALSE,
    r_delete BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    role_id INTEGER REFERENCES roles(id),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    last_access TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default admin role
INSERT INTO roles (name) VALUES ('Admin') ON CONFLICT DO NOTHING;

-- Insert default role rights for admin
INSERT INTO role_rights (role_id, section, route, r_create, r_read, r_update, r_delete)
SELECT 
    r.id,
    'be',
    '/users/user',
    TRUE,
    TRUE,
    TRUE,
    TRUE
FROM roles r
WHERE r.name = 'Admin'
ON CONFLICT DO NOTHING;

-- Insert default admin user
INSERT INTO users (role_id, name, email, password)
SELECT 
    r.id,
    'Administrator',
    'admin@gmail.com',
    '$2a$10$YourHashedPasswordHere' -- Replace with actual hashed password
FROM roles r
WHERE r.name = 'Admin'
ON CONFLICT DO NOTHING; 