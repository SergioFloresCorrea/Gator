-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL,
	feed_id INT NOT NULL,
	UNIQUE(user_id, feed_id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
