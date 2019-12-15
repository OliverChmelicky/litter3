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
    id        VARCHAR(50) PRIMARY KEY,
    firstName VARCHAR(50) NOT NULL,
    lastName  VARCHAR(50) NOT NULL,
    email     VARCHAR(50) NOT NULL,
    CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

create table societies
(
    id     VARCHAR(50) PRIMARY KEY,
    name   VARCHAR(50) NOT NULL,
    region VARCHAR(50) NOT NULL,
    admin  VARCHAR(50) REFERENCES users (id)
);

create table comments
(
    id        VARCHAR(50) PRIMARY KEY,
    description  VARCHAR(200),
    date  BIGINT  NOT NULL,
    trash_exists boolean,
    trash VARCHAR(50) REFERENCES users (id),
    "user"  VARCHAR(50) REFERENCES users, /*user is a keyword and should be quoted*/
    society VARCHAR(50) REFERENCES societies,
    CONSTRAINT exclusive_writer CHECK ( ("user" is null and society is not null) or ("user" is not null and society is null  ))
);

create table events
(
    id    VARCHAR(50) PRIMARY KEY,
    date  BIGINT  NOT NULL,
    publc boolean NOT NULL,
    userCreator VARCHAR(50) REFERENCES users,
    societyCreator VARCHAR(50) REFERENCES societies
        CONSTRAINT  exclusive_creator CHECK ( (userCreator is null and societyCreator is not null) or (userCreator is not null and societyCreator is null  ))
);

create table collections
(
    id      VARCHAR(50) PRIMARY KEY,
    date    BIGINT  NOT NULL,
    cleaned boolean NOT NULL
);

create table trash
(
    id            varchar(50) PRIMARY KEY,
    size          size   NOT NULL,
    place         VARCHAR(50),
    found         BIGINT NOT NULL,
    cleaned       BOOLEAN NOT NULL,
    accessibility accessibility NOT NULL,
    gps           point NOT NULL,
    finder        VARCHAR(50) REFERENCES users
);



create table users_societies
(
    "user"  VARCHAR(50) REFERENCES users,
    society VARCHAR(50) REFERENCES societies,
    PRIMARY KEY ("user", society)
);

create table users_collected
(
    "user" varchar(50) REFERENCES users,
    collection  varchar(50) REFERENCES collections,
    PRIMARY KEY ("user", collection)
);

create table societies_events
(
    society varchar(50) REFERENCES societies,
    event  varchar(50) REFERENCES events,
    PRIMARY KEY (society, event)
);

create table users_events
(
    "user" varchar(50) REFERENCES users,
    event  varchar(50) REFERENCES events,
    PRIMARY KEY ("user", event)
);