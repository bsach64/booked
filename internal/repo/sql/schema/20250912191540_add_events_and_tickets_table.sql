-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	time TIMESTAMP NOT NULL,
	address TEXT NOT NULL,
	description TEXT NOT NULL,
	seat_count BIGINT NOT NULL,
	latitude DOUBLE PRECISION,
	longitude DOUBLE PRECISION,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE ticket_status AS ENUM ('available', 'booked');

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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tickets;
DROP TYPE ticket_status;
DROP TABLE events;
-- +goose StatementEnd
