-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	time TIMESTAMP NOT NULL,
	address TEXT NOT NULL,
	description TEXT NOT NULL,
	latitude DOUBLE PRECISION,
	longitude DOUBLE PRECISION,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_time ON events(time);
CREATE TYPE ticket_status AS ENUM ('available', 'booked', 'cancelled');

CREATE TABLE tickets (
	id UUID PRIMARY KEY,
	user_id UUID,
	event_id UUID NOT NULL,
	status ticket_status NOT NULL DEFAULT 'available',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

CREATE INDEX idx_tickets_event_id ON tickets(event_id);
CREATE INDEX idx_tickets_user_id ON tickets(user_id);
CREATE INDEX idx_tickets_event_status ON tickets(event_id, status);
CREATE INDEX idx_tickets_user_status ON tickets(user_id, status);
CREATE INDEX idx_tickets_event_status_updated ON tickets(event_id, status, updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_events_time;
DROP TABLE tickets;
DROP TYPE ticket_status;
DROP TABLE events;
DROP INDEX IF EXISTS idx_tickets_event_status_updated;
DROP INDEX IF EXISTS idx_tickets_user_status;
DROP INDEX IF EXISTS idx_tickets_event_status;
DROP INDEX IF EXISTS idx_tickets_user_id;
DROP INDEX IF EXISTS idx_tickets_event_id;
-- +goose StatementEnd
