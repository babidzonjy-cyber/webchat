CREATE TABLE webchat.room_members(
    room_id INT NOT NULL REFERENCES webchat.rooms(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES webchat.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (room_id, user_id)
);

CREATE INDEX idx_room_members_user_id ON webchat.room_members(user_id);
