-- name: CreateTickets :copyfrom
INSERT INTO tickets (id, event_id)
VALUES ($1, $2);

-- name: GetAvailableTickets :many
SELECT id FROM tickets WHERE event_id = $1 AND status = 'available';

-- name: BookTickets :exec
UPDATE tickets SET user_id = $1, status = 'booked', updated_at = NOW() WHERE id = ANY($2::uuid[]);

-- name: GetBookingHistory :many
SELECT
	events.id,
	name,
	time,
	address,
	description,
	latitude,
	longitude,
	COUNT(tickets.id),
	MAX(tickets.updated_at)::TIMESTAMP
FROM
	tickets
JOIN
	events
ON
	tickets.event_id = events.id
WHERE
	tickets.user_id = $1 AND
	tickets.status = 'booked'
GROUP BY events.id;

-- name: GetBookedTickets :many
SELECT id FROM tickets WHERE event_id = $1 AND user_id = $2 AND status = 'booked';

-- name: CancelTickets :exec
UPDATE tickets SET user_id = NULL, status = 'available', updated_at = NOW() WHERE id = ANY($1::uuid[]);
