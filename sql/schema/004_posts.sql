-- +goose up
CREATE TABLE posts (
	id UUID NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL UNIQUE,
	description TEXT,
	published_at TIMESTAMP,
	feed_id UUID NOT NULL
);

-- +goose down
DROP TABLE posts;
