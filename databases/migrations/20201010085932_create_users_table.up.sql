CREATE TABLE IF NOT EXISTS users
(
    id       BIGSERIAL    NOT NULL,
    email    VARCHAR(120) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
