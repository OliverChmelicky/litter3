DROP TABLE IF EXISTS users_societies;
DROP TABLE IF EXISTS users_collected;
DROP TABLE IF EXISTS societies_events;
DROP TABLE IF EXISTS users_events;

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS trash; /*trash is countable and uncountable*/
DROP TABLE IF EXISTS societies;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS accessibility
DROP TYPE IF EXISTS "size"

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
    id        VARCHAR PRIMARY KEY,
    firstName VARCHAR NOT NULL,
    lastName  VARCHAR NOT NULL,
    email     VARCHAR NOT NULL,
    CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

create table societies
(
    id     VARCHAR PRIMARY KEY,
    name   VARCHAR NOT NULL,
    region VARCHAR NOT NULL,
    admin  VARCHAR REFERENCES users (id)
);

create table comments
(
    id        VARCHAR PRIMARY KEY,
    description  VARCHAR),
    date  BIGINT  NOT NULL,
    trash_exists boolean,
    trash VARCHAR REFERENCES users (id),
    "user"  VARCHAR REFERENCES users, /*user is a keyword and should be quoted*/
    society VARCHAR REFERENCES societies,
    CONSTRAINT exclusive_writer CHECK ( ("user" is null and society is not null) or ("user" is not null and society is null  ))
);

create table events
(
    id    VARCHAR PRIMARY KEY,
    date  BIGINT  NOT NULL,
    publc boolean NOT NULL,
    userCreator VARCHAR REFERENCES users,
    societyCreator VARCHAR REFERENCES societies
        CONSTRAINT  exclusive_creator CHECK ( (userCreator is null and societyCreator is not null) or (userCreator is not null and societyCreator is null  ))
);

create table collections
(
    id      VARCHAR PRIMARY KEY,
    date    BIGINT  NOT NULL,
    cleaned boolean NOT NULL
);

create table trash
(
    id            VARCHAR PRIMARY KEY,
    size          size   NOT NULL,
    place         VARCHAR,
    found         BIGINT NOT NULL,
    cleaned       BOOLEAN NOT NULL,
    accessibility accessibility NOT NULL,
    gps           point NOT NULL,
    finder        VARCHAR REFERENCES users
);



create table users_societies
(
    "user"  VARCHAR REFERENCES users,
    society VARCHAR REFERENCES societies,
    PRIMARY KEY ("user", society)
);

create table users_collected
(
    "user" VARCHAR REFERENCES users,
    collection  VARCHAR REFERENCES collections,
    PRIMARY KEY ("user", collection)
);

create table societies_events
(
    society VARCHAR REFERENCES societies,
    event  VARCHAR REFERENCES events,
    PRIMARY KEY (society, event)
);

create table users_events
(
    "user" VARCHAR REFERENCES users,
    event  VARCHAR REFERENCES events,
    PRIMARY KEY ("user", event)
);