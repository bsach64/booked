-- name: CreateUser :exec
INSERT INTO users (
	id, name, hashed_password, email, role
) VALUES (
	$1, $2, $3, $4, $5
);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
