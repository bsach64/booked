-- +goose Up
CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE users (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	hashed_password TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	role user_role NOT NULL DEFAULT 'user',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
DROP TYPE user_role;
