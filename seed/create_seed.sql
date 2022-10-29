CREATE TABLE IF NOT EXISTS currency (
    code VARCHAR NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS exchange(
    currency_from VARCHAR NOT NULL REFERENCES currency ON DELETE CASCADE ON UPDATE CASCADE,
    currency_to VARCHAR NOT NULL REFERENCES currency ON DELETE CASCADE ON UPDATE CASCADE,
    currency_rate FLOAT NOT NULL,
    markup_rate FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(currency_from, currency_to)
);

CREATE TABLE IF NOT EXISTS user(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS wallet(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    balance FLOAT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS user_wallets(
    user_id INTEGER REFERENCES user ON DELETE CASCADE ON UPDATE CASCADE,
    currency_code VARCHAR REFERENCES currency ON DELETE CASCADE ON UPDATE CASCADE,
    wallet_id INTEGER REFERENCES wallet ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, currency_code, wallet_id)
);

INSERT INTO currency (code) VALUES ('TRY');
INSERT INTO currency (code) VALUES ('USD');

INSERT INTO exchange (currency_from, currency_to, currency_rate)
VALUES ('TRY', 'USD', 1.40, 0.2);

INSERT INTO exchange (currency_from, currency_to, currency_rate)
VALUES ('USD', 'TRY', 16.00, 0.3);