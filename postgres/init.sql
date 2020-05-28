DROP TABLE IF EXISTS users_societies_members;
DROP TABLE IF EXISTS users_societies_admins;
DROP TABLE IF EXISTS users_collections;
DROP TABLE IF EXISTS societies_events;
DROP TABLE IF EXISTS users_events;

DROP TABLE IF EXISTS trash_comments;
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

CREATE TYPE size AS ENUM (
    'unknown',
    'bag',
    'wheelbarrow',
    'car'
    );

CREATE TYPE membership AS ENUM (
    'admin',
    'editor',
    'member'
    );

CREATE TYPE eventRights AS ENUM (
    'creator',
    'editor',
    'viewer'
    );



create table users
(
    id         VARCHAR PRIMARY KEY,
    first_name VARCHAR,
    last_name  VARCHAR,
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
    description varchar,
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
    trash_id   VARCHAR REFERENCES trash (id) on delete cascade not null,
    user_id    VARCHAR REFERENCES users (id) on delete set null not null,
    message    varchar     not null,
    created_at timestamptz NOT NULL
);

create table events
(
    id          VARCHAR PRIMARY KEY,
    date        timestamptz NOT NULL,
    description varchar,
    created_at  timestamptz NOT NULL
);

create table collections
(
    id            VARCHAR PRIMARY KEY,
    trash_id      VARCHAR REFERENCES trash (id) on delete cascade,
    event_id      VARCHAR REFERENCES events (id) on delete cascade, --maybe trigger that it has at least one user-collection or not null event_id. Can have both
    weight        real,
    cleaned_trash boolean     NOT NULL,
    created_at    timestamptz NOT NULL,
    CONSTRAINT correct_weight CHECK (weight > 0.0)
);


create table societies_members
(
    "user_id"  VARCHAR REFERENCES users (id) on delete cascade ,
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
    "user_id"     VARCHAR REFERENCES users (id) on delete cascade,  -- could be trigger that if he is the last then delete the whole collection
    collection_id VARCHAR REFERENCES collections (id) on delete cascade,
    PRIMARY KEY ("user_id", collection_id)
);

create table events_societies
(
    society_id VARCHAR REFERENCES societies (id) on delete cascade , ----> trigger on delete society if society is organizer
    event_id   VARCHAR REFERENCES events (id) on delete cascade,
    permission eventRights not null,
    PRIMARY KEY (society_id, event_id)
);

create table events_users
(
    user_id    VARCHAR REFERENCES users (id) on delete cascade , ----> trigger on delete user if user is organizer
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

CREATE OR REPLACE FUNCTION del_societies_events()
    RETURNS trigger AS
    $$
    DECLARE
        org_events  TEXT[];
    BEGIN
        IF EXISTS (SELECT FROM events_societies WHERE (society_id = OLD.id and permission = 'creator')) THEN
            org_events = Array(SELECT event_id FROM events_societies WHERE society_id = OLD.id and permission = 'creator');
            delete from events where id = ANY(org_events);
        END IF;
        return OLD;
    END;
    $$
    LANGUAGE plpgsql;

CREATE TRIGGER del_soc_trigger
    BEFORE DELETE ON societies
    FOR EACH ROW
    EXECUTE PROCEDURE del_societies_events();

--
--

CREATE OR REPLACE FUNCTION del_user_events()
    RETURNS trigger AS
$$
DECLARE
    org_events  TEXT[];
BEGIN
    IF EXISTS (SELECT FROM events_users WHERE (user_id = OLD.id and permission = 'creator')) THEN
        org_events = Array(SELECT event_id FROM events_users WHERE user_id = OLD.id and permission = 'creator');
        delete from events where id = ANY(org_events);
    END IF;
    return OLD;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER del_user_trigger
    BEFORE DELETE ON users
    FOR EACH ROW
EXECUTE PROCEDURE del_user_events();




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

insert into trash (id, cleaned, size, accessibility, trash_type, location, description, finder_id, created_at)
values ('1', false, 'bag', 'easy', 0, ST_GeomFromText('POINT(48 19)', 4326), '', null, '2003-2-11');
--
insert into collections (id, trash_id, event_id, weight, cleaned_trash, created_at)
values ('1','1', '', 655, false, '2005-2-11')
