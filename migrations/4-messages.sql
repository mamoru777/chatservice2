-- +migrate Up
CREATE TABLE messages
(
    id UUID DEFAULT UUID_GENERATE_V4()
        CONSTRAINT messages_id_pkey
            PRIMARY KEY,
    chat_id UUID NOT NULL
        CONSTRAINT messages_chat_id_fkey
            REFERENCES chats(id),
    usr_id UUID NOT NULL
        CONSTRAINT messages_usr_id_fkey
            REFERENCES usrs(id),
    text TEXT NOT NULL,
    data TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE messages;