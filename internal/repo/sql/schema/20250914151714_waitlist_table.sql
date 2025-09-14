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

CREATE INDEX idx_waitlist_user_id ON waitlist(user_id);
CREATE INDEX idx_waitlist_event_id ON waitlist(event_id);
CREATE UNIQUE INDEX idx_waitlist_user_event ON waitlist(user_id, event_id);
CREATE INDEX idx_waitlist_status ON waitlist(status);
-- +goose StatementEnd

-- +goose Down
DROP TABLE waitlist;
DROP TYPE notification_status;
DROP INDEX IF EXISTS idx_waitlist_status;
DROP INDEX IF EXISTS idx_waitlist_user_event;
DROP INDEX IF EXISTS idx_waitlist_event_id;
DROP INDEX IF EXISTS idx_waitlist_user_id;
-- +goose StatementBegin
-- +goose StatementEnd
