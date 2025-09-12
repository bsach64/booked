-- name: CreateEvent :exec
INSERT INTO events (
	id, name, time, address, description, latitude, longitude
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
);
