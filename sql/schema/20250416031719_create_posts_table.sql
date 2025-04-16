-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts(
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title TEXT,
	url TEXT UNIQUE NOT NULL,
	description TEXT,
	published_at TIMESTAMP NOT NULL,
	feed_id INT NOT NULL,
	FOREIGN KEY (feed_id)
	REFERENCES feeds (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
