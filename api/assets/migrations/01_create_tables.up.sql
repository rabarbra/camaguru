CREATE TABLE IF NOT EXISTS users
(
    id              INTEGER PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    email           TEXT NOT NULL UNIQUE,
    pass            TEXT NOT NULL,
    email_verified  INTEGER,
    like_notify     INTEGER,
    comm_notify     INTEGER
);

CREATE TABLE IF NOT EXISTS imgs
(
    id              INTEGER PRIMARY KEY,
    link            TEXT NOT NULL,
    created_at      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS comments
(
    id              INTEGER PRIMARY KEY,
    content         TEXT,
    created_at      TEXT NOT NULL,
    user_id         INTEGER,
    img_id          INTEGER,
    FOREIGN KEY(user_id)    REFERENCES users(id),
    FOREIGN KEY(img_id)     REFERENCES imgs(id)
);

CREATE TABLE IF NOT EXISTS likes
(
    id              INTEGER PRIMARY KEY,
    user_id         INTEGER,
    img_id          INTEGER,
    FOREIGN KEY(user_id)    REFERENCES users(id),
    FOREIGN KEY(img_id)     REFERENCES imgs(id)
);

PRAGMA foreign_keys = ON;
