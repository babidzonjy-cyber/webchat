DROP INDEX IF EXISTS webchat.idx_messages_created_at;
DROP INDEX IF EXISTS webchat.idx_messages_user_id;
DROP INDEX IF EXISTS webchat.idx_messages_room_id;
DROP INDEX IF EXISTS webchat.idx_rooms_created_by;

DROP TABLE IF EXISTS webchat.messages;
DROP TABLE IF EXISTS webchat.rooms;
DROP TABLE IF EXISTS webchat.users;
