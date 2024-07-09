CREATE TABLE IF NOT EXISTS users
(
    id              INT GENERATED ALWAYS AS IDENTITY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    email           VARCHAR(50) NOT NULL UNIQUE,
    pass            VARCHAR(100) NOT NULL,
    email_verified  BOOLEAN DEFAULT f,
    like_notify     BOOLEAN DEFAULT t,
    comm_notify     BOOLEAN DEFAULT t,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS imgs
(
    id              INT GENERATED ALWAYS AS IDENTITY,
    link            VARCHAR(50) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id         INT,
    PRIMARY KEY(id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
	    REFERENCES users(id)
	    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments
(
    id              INT GENERATED ALWAYS AS IDENTITY,
    content         TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id         INT NOT NULL,
    img_id          INT NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
	    REFERENCES users(id)
	    ON DELETE CASCADE,
    CONSTRAINT fk_img
        FOREIGN KEY(img_id)
	    REFERENCES imgs(id)
	    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes
(
    id              INT GENERATED ALWAYS AS IDENTITY,
    user_id         INT NOT NULL,
    img_id          INT NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
	    REFERENCES users(id)
	    ON DELETE CASCADE,
    CONSTRAINT fk_img
        FOREIGN KEY(img_id) 
	    REFERENCES imgs(id)
	    ON DELETE CASCADE
);
