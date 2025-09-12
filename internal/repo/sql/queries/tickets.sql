-- name: CreateTickets :copyfrom
INSERT INTO tickets (id, event_id)
VALUES ($1, $2);

-- name: GetAvailableTickets :many
SELECT id FROM tickets WHERE event_id = $1 AND status = 'available';
