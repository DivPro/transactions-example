-- +goose Up
CREATE TABLE users
(
    id          UUID    NOT NULL PRIMARY KEY,
    first_name  VARCHAR NOT NULL,
    last_name   VARCHAR NOT NULL,
    second_name VARCHAR NOT NULL         DEFAULT '',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE deposits
(
    id         UUID                                   NOT NULL PRIMARY KEY,
    user_id    UUID                                   NOT NULL REFERENCES users,
    amount     DECIMAL                                NOT NULL CHECK ( amount > 0 ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
CREATE INDEX ON deposits (user_id);

CREATE TABLE transactions
(
    id         UUID                                   NOT NULL,
    user_id    UUID                                   NOT NULL REFERENCES users,
    target_id  UUID                                   NOT NULL REFERENCES users,
    amount     DECIMAL                                NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
) PARTITION BY HASH(user_id);
CREATE INDEX ON transactions (id);
CREATE INDEX ON transactions (user_id);
CREATE INDEX ON transactions (target_id);

CREATE TABLE transactions_1 partition of transactions FOR VALUES WITH (MODULUS 3, REMAINDER 0);
CREATE TABLE transactions_2 partition of transactions FOR VALUES WITH (MODULUS 3, REMAINDER 1);
CREATE TABLE transactions_3 partition of transactions FOR VALUES WITH (MODULUS 3, REMAINDER 2);


CREATE TABLE balances
(
    user_id UUID    NOT NULL REFERENCES users PRIMARY KEY,
    amount  DECIMAL NOT NULL CHECK ( amount > 0 )
);

-- +goose Down
DROP TABLE balances;
DROP TABLE transactions;
DROP TABLE deposits;
DROP TABLE users;