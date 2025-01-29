-- +goose up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL
);

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

CREATE TABLE feed_follows (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL,
	feed_id UUID NOT NULL,
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) 
	REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) 
	REFERENCES feeds(id) ON DELETE CASCADE
);

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
