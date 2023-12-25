-- +migrate Up
CREATE TABLE chats
(
    id UUID DEFAULT UUID_GENERATE_V4()
        CONSTRAINT chats_id_pkey
            PRIMARY KEY
);

-- +migrate Down
DROP TABLE chats;