-- name: CreateEvent :exec
INSERT INTO events (
	id, name, time, address, description, latitude, longitude
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
);

-- name: GetEvents :many
SELECT * FROM events ORDER BY time ASC LIMIT $1;

-- name: GetNextEvents :many
SELECT * FROM events WHERE time > $1 ORDER BY time ASC LIMIT $2;
