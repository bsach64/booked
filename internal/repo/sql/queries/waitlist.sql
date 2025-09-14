-- name: GetWaitlistEntry :one
SELECT * FROM waitlist WHERE user_id = $1 AND event_id = $2;

-- name: UpdateWaitlistStatus :exec
UPDATE waitlist SET status = $1 WHERE id = ANY($2::uuid[]);

-- name: AddToWaitlist :exec
INSERT INTO waitlist (id, user_id, event_id, count, status)
VALUES ($1, $2, $3, $4, 'to_notify');

-- name: GetWaitlistNotificationDetails :many
SELECT
    u.name AS user_name,
    u.email AS user_email,
    e.name AS event_name,
    w."count",
    e.time AS event_time,
    e.id AS event_id,
    w.id AS waitlist_id
FROM users u
JOIN waitlist w ON u.id = w.user_id
JOIN events e ON e.id = w.event_id
WHERE w.status = 'to_notify';
