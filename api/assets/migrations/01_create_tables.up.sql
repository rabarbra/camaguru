CREATE TABLE IF NOT EXISTS users
(
    id              INT GENERATED ALWAYS AS IDENTITY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    email           VARCHAR(50) NOT NULL UNIQUE,
    pass            VARCHAR(100) NOT NULL,
    avatar          VARCHAR(500),
    email_verified  BOOLEAN DEFAULT FALSE,
    like_notify     BOOLEAN DEFAULT TRUE,
    comm_notify     BOOLEAN DEFAULT TRUE,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS imgs
(
    id              INT GENERATED ALWAYS AS IDENTITY,
    link            VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id         INT,
    img_type        VARCHAR(50) DEFAULT 'img',
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
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
