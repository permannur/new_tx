CREATE TABLE tbl_department
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(192) NOT NULL,
    state    VARCHAR(10)  NOT NULL,
    priority INTEGER      NOT NULL
);

CREATE UNIQUE INDEX uq_department_name ON tbl_department (name)
    WHERE
        state = 'ENABLED';

-----------------------------------------------------------------------------------------------------

CREATE TABLE tbl_position
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(128) NOT NULL,
    state    VARCHAR(10)  NOT NULL,
    priority INTEGER      NOT NULL
);

CREATE UNIQUE INDEX uq_position_name ON tbl_position (name)
    WHERE
        state = 'ENABLED';

-----------------------------------------------------------------------------------------------------

CREATE TABLE tbl_user
(
    id            serial PRIMARY KEY,
    username      varchar(32)                 NOT NULL,
    password      VARCHAR(256)                NOT NULL,
    firstname     varchar(64)                 NOT NULL,
    lastname      varchar(64)                 NOT NULL,
    department_id int                         NOT NULL,
    position_id   int                         NOT NULL,
    state         varchar(10)                 NOT NULL,
    create_ts     timestamp WITHOUT time zone NOT NULL,
    update_ts     timestamp WITHOUT time zone NOT NULL,
    version       int                         NOT NULL
);

CREATE UNIQUE INDEX uq_user_username ON tbl_user (username)
    WHERE
        state != 'DELETED';

ALTER TABLE tbl_user
    ADD CONSTRAINT fk_user_department_id FOREIGN KEY (department_id)
        REFERENCES tbl_department (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    ADD CONSTRAINT fk_user_position_id FOREIGN KEY (position_id)
        REFERENCES tbl_position (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-----------------------------------------------------------------------------------------------------


CREATE TABLE tbl_user_log
(
    id        uuid PRIMARY KEY,
    user_id   int                         NOT NULL,
    username  varchar(32)                 NOT NULL,
    ip        inet                        NOT NULL,
    action    varchar(32)                 NOT NULL,
    action_ts timestamp WITHOUT time zone NOT NULL,
    sup_info  json                        NOT NULL
);

ALTER TABLE tbl_user_log
    ADD CONSTRAINT fk_user_log_user_id FOREIGN KEY (user_id)
        REFERENCES tbl_user (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-----------------------------------------------------------------------------------------------------

INSERT INTO tbl_department (name, state, priority)
VALUES ('admin', 'ENABLED', 1);

INSERT INTO tbl_position (name, state, priority)
VALUES ('admin', 'ENABLED', 1);

INSERT INTO tbl_user (username, PASSWORD, firstname, lastname, department_id, position_id, state,
                      create_ts, update_ts, version)
VALUES ( 'SYSTEM', '$2a$10$xKAGH3MypqXqOwylyBxPJuFSBRf5j4Ya4yK2Z5pUTMWY5PGs2Zgs2' --123qweASD!
       , 'SYSTEM', 'SYSTEM', 1, 1, 'ACTIVE', now(), now(), 0);

INSERT INTO tbl_user (username, PASSWORD, firstname, lastname, department_id, position_id, state,
                      create_ts, update_ts, version)
VALUES ( 'SUPERADMIN', '$2a$10$xKAGH3MypqXqOwylyBxPJuFSBRf5j4Ya4yK2Z5pUTMWY5PGs2Zgs2' --123qweASD!
       , 'SUPERADMIN', 'SUPERADMIN', 1, 1, 'ACTIVE', now(), now(), 0);
