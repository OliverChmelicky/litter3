DROP TABLE IF EXISTS users_societies_members;
DROP TABLE IF EXISTS users_societies_admins;
DROP TABLE IF EXISTS users_collections;
DROP TABLE IF EXISTS societies_events;
DROP TABLE IF EXISTS users_events;

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS trash; /*trash is countable and uncountable*/
DROP TABLE IF EXISTS societies;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS accessibility;
DROP TYPE IF EXISTS "size";

CREATE TYPE accessibility AS ENUM (
    'unknown',
    'easy',
    'car',
    'cave',
    'underWater'
    );

CREATE TYPE trashType AS ENUM (
    'unknown',
    'household',
    'automotive',
    'construction',
    'plastics',
    'electronic',
    'glass',
    'metal',
    'liquid',
    'dangerous',
    'carcass',
    'organic'
    );

CREATE TYPE membership AS ENUM (
    'admin',
    'member'
    );

CREATE TYPE size AS ENUM (
    'unknown',
    'bag',
    'wheelbarrow',
    'car'
    );



create table users
(
    id         VARCHAR PRIMARY KEY,
    first_name VARCHAR     NOT NULL,
    last_name  VARCHAR     NOT NULL,
    email      VARCHAR     NOT NULL,
    CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    created_at timestamptz NOT NULL
);

create table societies
(
    id         VARCHAR PRIMARY KEY,
    name       VARCHAR     NOT NULL,
    created_at timestamptz NOT NULL
);

create table trash
(
    id            VARCHAR PRIMARY KEY,
    cleaned       BOOLEAN                NOT NULL default false,
    size          size                   NOT NULL default 'unknown',
    accessibility accessibility          NOT NULL default 'unknown',
    trash_type    trashType              NOT NULL default 'unknown',
    location      GEOGRAPHY(POINT, 4326) NOT NULL,
    description   VARCHAR,
    finder_id     VARCHAR REFERENCES users,
    created_at    timestamptz            NOT NULL
);

create table comments
(
    id           VARCHAR PRIMARY KEY,
    description  VARCHAR,
    trash_exists boolean,
    trash_id     VARCHAR REFERENCES trash (id),
    "user_id"    VARCHAR REFERENCES users (id), /*user is a keyword and should be quoted*/
    society_id   VARCHAR REFERENCES societies (id),
    CONSTRAINT exclusive_writer CHECK ( ("user_id" is null and society_id is not null) or
                                        ("user_id" is not null and society_id is null)),
    created_at   timestamptz NOT NULL
);

create table events
(
    id         VARCHAR PRIMARY KEY,
    date       BIGINT      NOT NULL,
    publc      boolean     NOT NULL,
    user_id    VARCHAR REFERENCES users (id),
    society_id VARCHAR REFERENCES societies (id),
    CONSTRAINT exclusive_creator CHECK ( (user_id is null and society_id is not null) or
                                         (user_id is not null and society_id is null)),
    created_at timestamptz NOT NULL
);

create table collections
(
    id            VARCHAR PRIMARY KEY,
    trash_id      VARCHAR REFERENCES trash (id),
    cleaned_trash boolean     NOT NULL,
    created_at    timestamptz NOT NULL
);


create table societies_members
(
    "user_id"  VARCHAR REFERENCES users (id),
    society_id VARCHAR REFERENCES societies (id),
    permission membership  not null,
    created_at timestamptz NOT NULL,
    PRIMARY KEY ("user_id", society_id)
);

create table societies_applicants
(
    "user_id"  VARCHAR REFERENCES users (id),
    society_id VARCHAR REFERENCES societies (id),
    created_at timestamptz NOT NULL,
    PRIMARY KEY ("user_id", society_id)
);

create table users_collections
(
    "user_id"     VARCHAR REFERENCES users (id),
    collection_id VARCHAR REFERENCES collections (id),
    created_at    timestamptz NOT NULL,
    PRIMARY KEY ("user_id", collection_id)
);


CREATE TYPE eventRights AS ENUM (
    'admin',
    'viewer'
    );

/*mozno zjednot societies_events a users_events*/
create table societies_events
(
    society_id VARCHAR REFERENCES societies (id),
    event_id   VARCHAR REFERENCES events (id),
    permission eventRights not null,
    created_at timestamptz NOT NULL,
    PRIMARY KEY (society_id, event_id)
);

create table users_events
(
    "user_id"  VARCHAR REFERENCES users (id),
    event_id   VARCHAR REFERENCES events (id),
    permission eventRights not null,
    created_at timestamptz NOT NULL,
    PRIMARY KEY ("user_id", event_id)
);


create table friends
(
    user1_id   VARCHAR REFERENCES users (id),
    user2_id   VARCHAR REFERENCES users (id),
    created_at timestamptz NOT NULL,
    PRIMARY KEY (user1_id, user2_id)
);

create table friend_requests
(
    user1_id   VARCHAR REFERENCES users (id),
    user2_id   VARCHAR REFERENCES users (id),
    created_at timestamptz NOT NULL,
    PRIMARY KEY (user1_id, user2_id)
);

create table trash_comments
(
    Id         varchar PRIMARY KEY,
    UserId     varchar     not null,
    TrashId    varchar     not null,
    message    varchar     not null,
    user1_id   VARCHAR REFERENCES users (id),
    user2_id   VARCHAR REFERENCES users (id),
    created_at timestamptz NOT NULL
);

--not possible before create
-- CREATE OR REPLACE FUNCTION correct_create_friends()
--     RETURNS trigger AS
-- $BODY$
-- BEGIN
--     IF  1 = STRCMP(NEW.user1_id, NEW.user2_id) THEN
--         tmp varchar;
--         tmp = new.user1_id;
--     END IF;
--
--     RETURN NEW;
-- END;
-- $BODY$
--
-- CREATE TRIGGER correctCreateFriendRequest
--     BEFORE UPDATE ON friend_requests
--     FOR EACH ROW
-- EXECUTE PROCEDURE correct_create_friends();