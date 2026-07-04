CREATE SCHEMA IF NOT EXISTS webchat;

CREATE TABLE webchat.users(
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL CHECK (char_length(full_name) >= 1),
    email VARCHAR(100) NOT NULL UNIQUE
    CHECK (char_length(email) BETWEEN 7 AND 100
    AND
    (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')),
    password VARCHAR(100) NOT NULL CHECK (char_length(password) BETWEEN 7 AND 100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE webchat.rooms(
    id SERIAL PRIMARY KEY,
    name VARCHAR(25) NOT NULL CHECK (char_length(name) >= 1),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by INT REFERENCES webchat.users(id) NOT NULL
);

CREATE TABLE webchat.messages(
    id SERIAL PRIMARY KEY,
    text VARCHAR(5000) NOT NULL CHECK (char_length(text) >= 1),
    room_id INT REFERENCES webchat.rooms(id) NOT NULL,
    user_id INT REFERENCES webchat.users(id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rooms_created_by ON webchat.rooms(created_by);
CREATE INDEX idx_messages_room_id ON webchat.messages(room_id);
CREATE INDEX idx_messages_user_id ON webchat.messages(user_id);
CREATE INDEX idx_messages_created_at ON webchat.messages(created_at);
