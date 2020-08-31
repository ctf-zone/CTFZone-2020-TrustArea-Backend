CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE logs (
        id SERIAL PRIMARY KEY,
        log_type VARCHAR NOT NULL,
        log_message VARCHAR NOT NULL,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
        id uuid DEFAULT uuid_generate_v4 () UNIQUE PRIMARY KEY,
        username VARCHAR NOT NULL UNIQUE,
        first_name VARCHAR NOT NULL,
        last_name VARCHAR NOT NULL,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
		id SERIAL PRIMARY KEY,
		owner uuid NOT NULL,
		description VARCHAR NOT NULL,
		challenge VARCHAR NOT NULL,
		reward VARCHAR NOT NULL,
		solved BOOLEAN NOT NULL DEFAULT FALSE,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (owner) REFERENCES users (id)
);

CREATE TABLE solutions (
		id SERIAL PRIMARY KEY,
		task_id INTEGER NOT NULL,
		solution VARCHAR NOT NULL,
		owner uuid NOT NULL,
		FOREIGN KEY (task_id) REFERENCES tasks (id),
		FOREIGN KEY (owner) REFERENCES users (id)
);
