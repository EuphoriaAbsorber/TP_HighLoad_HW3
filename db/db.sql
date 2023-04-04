CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users (
	email    CITEXT NOT NULL UNIQUE,
    fullname TEXT NOT NULL,
    nickname CITEXT NOT NULL UNIQUE PRIMARY KEY, 
    about    TEXT
);

CREATE UNLOGGED TABLE forums (
	title   TEXT NOT NULL,
    user1   CITEXT NOT NULL REFERENCES users (nickname),
    slug    CITEXT PRIMARY KEY,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0
);

CREATE UNLOGGED TABLE threads (
    id      SERIAL PRIMARY KEY,
	title   TEXT NOT NULL,
    author  CITEXT NOT NULL REFERENCES users (nickname),
    forum   CITEXT NOT NULL REFERENCES forums (slug),
    message TEXT NOT NULL,
    votes   INT DEFAULT 0,
    slug    CITEXT,
    created TIMESTAMP WITH TIME ZONE
);

CREATE UNLOGGED TABLE posts (
    id  SERIAL PRIMARY KEY,
    parent   INT REFERENCES posts (id),
    author   CITEXT NOT NULL REFERENCES users (nickname),
    message  TEXT NOT NULL,
    forum    CITEXT NOT NULL REFERENCES forums (slug),
    isedited BOOLEAN,
    thread   INT REFERENCES threads (id),
    created  TIMESTAMP WITH TIME ZONE,
    path     INT[]  DEFAULT ARRAY []::INTEGER[]
);

CREATE UNLOGGED TABLE votes (
    nickname  CITEXT NOT NULL REFERENCES users (nickname),
    thread    INT    NOT NULL REFERENCES threads (id),
    voice     INT    NOT NULL,
    CONSTRAINT vote_key UNIQUE (nickname, thread)
);

CREATE UNLOGGED TABLE IF NOT EXISTS forum_users
(
    nickname CITEXT COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
    forum    CITEXT NOT NULL REFERENCES forums (slug),
    CONSTRAINT user_forum_key UNIQUE (nickname, forum)
);

--TRIGGERS

CREATE OR REPLACE FUNCTION post_set_path()
    RETURNS TRIGGER AS
$$
DECLARE
parent_post_id posts.id%type := 0;
BEGIN
    NEW.path = (SELECT path FROM posts WHERE id = NEW.parent) || NEW.id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post
    BEFORE INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE post_set_path();




CREATE OR REPLACE FUNCTION add_user()
    RETURNS TRIGGER AS
$$
BEGIN
INSERT INTO forum_users (nickname, forum)
VALUES (NEW.author, NEW.forum)
    ON CONFLICT do nothing;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_new_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE add_user();

CREATE TRIGGER insert_new_post
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE add_user();


--INDEXES

CREATE UNIQUE INDEX IF NOT EXISTS votes_key ON votes (thread, nickname);

CREATE INDEX IF NOT EXISTS threads_created ON threads (created);
CREATE INDEX IF NOT EXISTS threads_slug ON threads (slug);
CREATE INDEX IF NOT EXISTS threads_forum ON threads (forum);

CREATE INDEX IF NOT EXISTS forums_slug ON forums (slug);
CREATE INDEX IF NOT EXISTS forums_user ON forums (user1);

CREATE INDEX IF NOT EXISTS users_nickname ON users (nickname);

CREATE INDEX IF NOT EXISTS posts_thread ON posts (thread);

CREATE INDEX IF NOT EXISTS forum_users_forum ON forum_users (forum);
CREATE INDEX IF NOT EXISTS forum_users_nickname ON forum_users (nickname);
CREATE INDEX IF NOT EXISTS forum_users_all ON forum_users (forum, nickname);

CREATE INDEX IF NOT EXISTS posts_path_path1_id ON posts (path, (path[1]), id);
CREATE INDEX IF NOT EXISTS posts_path1_id_thread_parent ON posts ((path[1]), thread, id, parent NULLS FIRST);
