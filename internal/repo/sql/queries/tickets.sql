-- name: CreateTickets :copyfrom
INSERT INTO tickets (id, event_id)
VALUES ($1, $2);

-- name: GetAvailableTickets :many
SELECT id FROM tickets WHERE event_id = $1 AND (status = 'available' OR status = 'cancelled');

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
UPDATE tickets SET user_id = NULL, status = 'cancelled', updated_at = NOW() WHERE id = ANY($1::uuid[]);

-- name: GetAnalytics :many
SELECT
	event_id,
	COUNT(id) AS total_seats,
	COUNT(id) FILTER (WHERE status = 'booked') AS booked_tickets,
	(COUNT(id)::DOUBLE PRECISION / COUNT(id) FILTER (WHERE status = 'booked')::DOUBLE PRECISION)::DOUBLE PRECISION AS capacity_utilisation
FROM
	tickets
GROUP BY
	event_id
ORDER BY
	(COUNT(id)::DOUBLE PRECISION / COUNT(id) FILTER (WHERE status = 'booked')::DOUBLE PRECISION)::DOUBLE PRECISION
DESC;

-- name: GetDailyBookings :many
SELECT
	event_id,
	COUNT(id) FILTER (WHERE status = 'booked' AND updated_at::date = CURRENT_DATE) AS today_booked_tickets
FROM
	tickets
GROUP BY
	event_id;

-- name: GetCancellationRates :many
SELECT
	event_id,
	(
		COUNT(id) FILTER (WHERE status = 'cancelled')::DOUBLE PRECISION
		/
		NULLIF(
			COUNT(id) FILTER (WHERE status = 'booked' OR status = 'cancelled')::DOUBLE PRECISION,
			0
		)
	)::DOUBLE PRECISION AS cancellation_rate
FROM
	tickets
GROUP BY
	event_id;
