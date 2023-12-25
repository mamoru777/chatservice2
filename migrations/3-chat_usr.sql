-- +migrate Up
CREATE TABLE chat_usr
(
    chat_id UUID
        CONSTRAINT chat_usr_chat_id_fkey
            REFERENCES chats(id),
    usr_id UUID
        CONSTRAINT chat_usr_usr_id_fkey
            REFERENCES usrs(id),
    CONSTRAINT chat_usr_id_pkey PRIMARY KEY (chat_id, usr_id)
);

-- +migrate Down
DROP TABLE chat_usr;
