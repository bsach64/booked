-- name: CreateTickets :copyfrom
INSERT INTO tickets (id, event_id)
VALUES ($1, $2);
