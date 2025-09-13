-- name: CreateEvent :exec
INSERT INTO events (
	id, name, time, address, description, latitude, longitude
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
);

-- name: GetEvents :many
SELECT
	e.*,
	COUNT(t.id) AS total_tickets,
	COUNT(t.id) FILTER (WHERE t.status = 'available') AS available_tickets
FROM
	events e
JOIN
	tickets t
ON
	e.id = t.event_id
GROUP BY
	e.id
ORDER BY
	e.time ASC
LIMIT $1;

-- name: GetNextEvents :many
SELECT 
	e.*,
	COUNT(t.id) AS total_tickets,
	COUNT(t.id) FILTER (WHERE t.status = 'available') AS available_tickets
FROM
	events e
JOIN
	tickets t
ON
	e.id = t.event_id
WHERE
	time > $1
GROUP BY
	e.id
ORDER BY
	e.time ASC
LIMIT $2;

-- name: DeleteEvent :one
DELETE FROM events CASCADE WHERE id = $1 RETURNING id;
