CREATE TABLE products
(
    id       SERIAL,
    name     TEXT  NOT NULL,
    price    FLOAT NOT NULL,
    articles JSONB NOT NULL DEFAULT '{}',

    PRIMARY KEY (id)
)
