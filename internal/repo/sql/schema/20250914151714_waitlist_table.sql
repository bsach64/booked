-- +goose Up
-- +goose StatementBegin
CREATE TYPE notification_status AS ENUM ('to_notify', 'notified');
CREATE TABLE waitlist (
	id UUID PRIMARY KEY,
	count INT NOT NULL,
	user_id UUID NOT NULL,
	event_id UUID NOT NULL,
	status notification_status NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE waitlist;
-- +goose StatementBegin
-- +goose StatementEnd
