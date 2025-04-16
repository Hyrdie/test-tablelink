-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    role_id INTEGER REFERENCES roles(id),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    last_access TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create role_rights table
CREATE TABLE IF NOT EXISTS role_rights (
    id SERIAL PRIMARY KEY,
    role_id INTEGER REFERENCES roles(id),
    section VARCHAR(50) NOT NULL,
    route VARCHAR(255) NOT NULL,
    r_create BOOLEAN DEFAULT false,
    r_read BOOLEAN DEFAULT false,
    r_update BOOLEAN DEFAULT false,
    r_delete BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, section, route)
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Insert default roles
INSERT INTO roles (name) VALUES ('admin') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name) VALUES ('user') ON CONFLICT (name) DO NOTHING;

-- Insert default admin user
INSERT INTO users (role_id, name, email, password) 
VALUES (1, 'Administrator', 'admin@gmail.com', 'adminadmin')
ON CONFLICT (email) DO NOTHING;

-- Insert default role rights for admin
INSERT INTO role_rights (role_id, section, route, r_create, r_read, r_update, r_delete)
VALUES 
    (1, 'be', '/users/user', true, true, true, true),
    (1, 'be', '/auth/login', true, true, true, true)
ON CONFLICT (role_id, section, route) DO NOTHING;

-- Insert default role rights for user
INSERT INTO role_rights (role_id, section, route, r_create, r_read, r_update, r_delete)
VALUES 
    (2, 'be', '/users/user', false, true, false, false),
    (2, 'be', '/auth/login', true, true, true, true)
ON CONFLICT (role_id, section, route) DO NOTHING; 