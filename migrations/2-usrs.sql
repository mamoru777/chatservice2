-- +migrate Up
CREATE TABLE usrs
(
    id UUID NOT NULL
        CONSTRAINT usrs_id_pkey
            PRIMARY KEY
);

-- +migrate Down
DROP TABLE usrs;