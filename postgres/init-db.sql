DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS societies;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS trash; /*trash is countable and uncountable*/

DROP TABLE IF EXISTS users_societies;
DROP TABLE IF EXISTS approved_trash;
DROP TABLE IF EXISTS events_trash;
DROP TABLE IF EXISTS users_collected_trash;
DROP TABLE IF EXISTS users_events;


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

create table events
(
    id    VARCHAR(50) PRIMARY KEY,
    date  BIGINT  NOT NULL,
    publc boolean NOT NULL
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
    "user"  VARCHAR(50) REFERENCES users, /*user is a keyword and should be quoted*/
    society VARCHAR(50) REFERENCES societies,
    PRIMARY KEY ("user", society)
);


create table approved_trash
(
    "user" VARCHAR(50),
    trash  VARCHAR(50),
    PRIMARY KEY ("user", trash)
);

create table events_trash
(
    event varchar(50) REFERENCES events,
    trash varchar(50) REFERENCES trash,
    PRIMARY KEY (event, trash)
);

create table users_collected_trash
(
    "user" varchar(50) REFERENCES users,
    trash  varchar(50) REFERENCES trash,
    PRIMARY KEY ("user", trash)
);

create table users_events
(
    "user" varchar(50) REFERENCES users,
    event  varchar(50) REFERENCES events,
    PRIMARY KEY ("user", event)
);