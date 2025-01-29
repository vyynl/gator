-- +goose up
CREATE TABLE feeds (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	last_fetched_at TIMESTAMP,
	name TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	user_id UUID NOT NULL,
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) 
	REFERENCES users(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE feeds;
