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

-- CREATE TYPE trashType AS ENUM (
--     'unknown',
--     'household',
--     'automotive',
--     'construction',
--     'plastics',
--     'electronic',
--     'glass',
--     'metal',
--     'liquid',
--     'dangerous',
--     'carcass',
--     'organic'
--     );

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
    email      VARCHAR     NOT NULL unique,
    uid        VARCHAR     NOT NULL unique,
    avatar     varchar,
    CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    created_at timestamptz NOT NULL
);

create table societies
(
    id         VARCHAR PRIMARY KEY,
    name       VARCHAR     NOT NULL,
    avatar     varchar,
    created_at timestamptz NOT NULL
);

create table trash
(
    id            VARCHAR PRIMARY KEY,
    cleaned       BOOLEAN       default false,
    size          size          default 'unknown',
    accessibility accessibility default 'unknown',
    trash_type    numeric,
    location      GEOGRAPHY(POINT, 4326) NOT NULL,
    description   VARCHAR,
    finder_id     VARCHAR                REFERENCES users (id) on delete set null,
    created_at    timestamptz            NOT NULL
);

create table trash_comments
(
    id         varchar PRIMARY KEY,
    trash_id   VARCHAR REFERENCES trash (id) on delete cascade,
    user_id    VARCHAR REFERENCES users (id) on delete cascade,
    message    varchar     not null,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

create table events
(
    id          VARCHAR PRIMARY KEY,
    date        timestamptz NOT NULL,
    publc       boolean     NOT NULL,
    description varchar,
--     user_id    VARCHAR REFERENCES users (id),
--     society_id VARCHAR REFERENCES societies (id),
--     CONSTRAINT exclusive_creator CHECK ( (user_id is null and society_id is not null) or
--                                          (user_id is not null and society_id is null)),
    created_at  timestamptz NOT NULL
);

create table collections
(
    id            VARCHAR PRIMARY KEY,
    trash_id      VARCHAR REFERENCES trash (id) on delete cascade,
    event_id      VARCHAR REFERENCES events (id) on delete cascade,
    weight        real,
    cleaned_trash boolean     NOT NULL,
    created_at    timestamptz NOT NULL,
    CONSTRAINT correct_weight CHECK (weight > 0.0)
);


create table societies_members
(
    "user_id"  VARCHAR REFERENCES users (id),
    society_id VARCHAR REFERENCES societies (id) on delete cascade,
    permission membership  not null,
    created_at timestamptz NOT NULL,
    PRIMARY KEY ("user_id", society_id)
);

create table societies_applicants
(
    "user_id"  VARCHAR REFERENCES users (id) on delete cascade,
    society_id VARCHAR REFERENCES societies (id) on delete cascade,
    created_at timestamptz NOT NULL,
    PRIMARY KEY ("user_id", society_id)
);

create table users_collections
(
    "user_id"     VARCHAR REFERENCES users (id) on delete cascade,
    collection_id VARCHAR REFERENCES collections (id) on delete cascade,
    PRIMARY KEY ("user_id", collection_id)
);


CREATE TYPE eventRights AS ENUM (
    'creator',
    'editor',
    'viewer'
    );

create table events_societies
(
    society_id VARCHAR REFERENCES societies (id), --cannot on delete cascade because it can be an organizer --> trigger or in app solution
    event_id   VARCHAR REFERENCES events (id) on delete cascade,
    permission eventRights not null,
    PRIMARY KEY (society_id, event_id)
);

create table events_users
(
    user_id    VARCHAR REFERENCES users (id), --cannot on delete cascade because he can be an organizer --> trigger or in app solution
    event_id   VARCHAR REFERENCES events (id) on delete cascade,
    permission eventRights not null,
    PRIMARY KEY ("user_id", event_id)
);

create table events_trash
(
    trash_id VARCHAR REFERENCES trash (id) on delete cascade,
    event_id VARCHAR REFERENCES events (id) on delete cascade,
    PRIMARY KEY (trash_id, event_id)
);


create table friends
(
    user1_id   VARCHAR REFERENCES users (id) on delete cascade,
    user2_id   VARCHAR REFERENCES users (id) on delete cascade,
    created_at timestamptz NOT NULL,
    PRIMARY KEY (user1_id, user2_id)
);

create table friend_requests
(
    user1_id   VARCHAR REFERENCES users (id) on delete cascade,
    user2_id   VARCHAR REFERENCES users (id) on delete cascade,
    created_at timestamptz NOT NULL,
    PRIMARY KEY (user1_id, user2_id)
);


create table trash_images
(
    trash_id VARCHAR REFERENCES trash (id) on delete cascade,
    url      varchar,
    PRIMARY KEY (trash_id, url)
);

create table collection_images
(
    collection_id VARCHAR REFERENCES collections (id) on delete cascade,
    url           varchar,
    PRIMARY KEY (collection_id, url)
);

INSERT INTO users (id, first_name, last_name, email, uid, created_at)
VALUES ('ad10c133-f18f-417d-be04-795e290683c1', 'Heinrich', 'Herrer', 'a@a.cz', 'vBP3XYCYQPZGMYAJVeVgd9pttkr2', '2003-2-1');

INSERT INTO users (id, first_name, last_name, email, uid, created_at)
VALUES ('92a99679-9e9f-401c-8ae4-eb9c869c2120', 'Peter', 'Aufschneider', 'b@b.cz', 'fCtIJKvNRbXAulI9IFScKoJcyuk2', '2003-2-1');


INSERT INTO societies (id, name, avatar, created_at)
VALUES ('1', 'Prve', 'Aufschneider', '2003-2-1');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('2', 'Druhe', 'Aufschneider', '2003-2-2');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('3', 'Tetie', 'Aufschneider', '2003-2-3');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('4', 'Stvrte', 'Aufschneider', '2003-2-4');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('5', 'Piate', 'Aufschneider', '2003-2-5');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('6', 'Sieste', 'Aufschneider', '2003-2-6');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('7', 'Siedme', 'Aufschneider', '2003-2-7');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('8', 'Osme', 'Aufschneider', '2003-2-8');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('9', 'Deviate', 'Aufschneider', '2003-2-9');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('10', 'Desiate', 'Aufschneider', '2003-2-10');
INSERT INTO societies (id, name, avatar, created_at)
VALUES ('11', 'Jedenast', 'Aufschneider', '2003-2-11');