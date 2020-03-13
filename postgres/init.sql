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
    'easy',
    'medium',
    'hard'
    );

CREATE TYPE size AS ENUM (
    'small',
    'medium',
    'big',
    'extremelyBig'
    );



create table users
(
    id         VARCHAR PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name  VARCHAR NOT NULL,
    email      VARCHAR NOT NULL,
    CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    created    BIGINT  NOT NULL
);

create table societies
(
    id      VARCHAR PRIMARY KEY,
    name    VARCHAR NOT NULL,
    created BIGINT  NOT NULL
);

create table trash
(
    id            VARCHAR PRIMARY KEY,
    cleaned       BOOLEAN       NOT NULL,
    size          size          NOT NULL,
    accessibility accessibility NOT NULL,
    gps           point         NOT NULL,
    description   VARCHAR,
    finder        VARCHAR REFERENCES users,
    created       BIGINT        NOT NULL
);

create table comments
(
    id           VARCHAR PRIMARY KEY,
    description  VARCHAR,
    trash_exists boolean,
    trash        VARCHAR REFERENCES trash (id),
    "user"       VARCHAR REFERENCES users (id), /*user is a keyword and should be quoted*/
    society      VARCHAR REFERENCES societies (id),
    CONSTRAINT exclusive_writer CHECK ( ("user" is null and society is not null) or
                                        ("user" is not null and society is null)),
    created      BIGINT NOT NULL
);

create table events
(
    id              VARCHAR PRIMARY KEY,
    date            BIGINT  NOT NULL,
    publc           boolean NOT NULL,
    user_creator    VARCHAR REFERENCES users (id),
    society_creator VARCHAR REFERENCES societies (id),
    CONSTRAINT exclusive_creator CHECK ( (user_creator is null and society_creator is not null) or
                                         (user_creator is not null and society_creator is null)),
    created         BIGINT  NOT NULL
);

create table collections
(
    id      VARCHAR PRIMARY KEY,
    trash   VARCHAR REFERENCES trash (id),
    cleaned boolean NOT NULL,
    created BIGINT  NOT NULL
);


create table users_societies_admins
(
    "user"  VARCHAR REFERENCES users (id),
    society VARCHAR REFERENCES societies (id),
    PRIMARY KEY ("user", society)
);

create table users_societies_members
(
    "user"  VARCHAR REFERENCES users (id),
    society VARCHAR REFERENCES societies (id),
    PRIMARY KEY ("user", society)
);

create table users_collections
(
    "user"     VARCHAR REFERENCES users (id),
    collection VARCHAR REFERENCES collections (id),
    PRIMARY KEY ("user", collection)
);

create table societies_events
(
    society VARCHAR REFERENCES societies (id),
    event   VARCHAR REFERENCES events (id),
    PRIMARY KEY (society, event)
);

create table users_events
(
    "user" VARCHAR REFERENCES users (id),
    event  VARCHAR REFERENCES events (id),
    PRIMARY KEY ("user", event)
);