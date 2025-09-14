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
	COUNT(t.id) FILTER (WHERE t.status = 'available' OR t.status = 'cancelled') AS available_tickets
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

-- name: GetEventByID :one
SELECT
	e.*,
	COUNT(t.id) AS total_tickets
FROM
	events e
JOIN
	tickets t
ON
	e.id = t.event_id
WHERE
	e.id = $1
GROUP BY
	e.id;


-- name: UpdateEvent :exec
UPDATE 
	events
SET
	name = $2, time = $3, address = $4, description = $5, latitude = $6, longitude = $7, updated_at = NOW()
WHERE
	id = $1;

