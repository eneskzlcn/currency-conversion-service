CREATE TABLE IF NOT EXISTS currencies (
    code VARCHAR NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS currency_exchange_values(
    currency_from VARCHAR NOT NULL REFERENCES currencies ON DELETE CASCADE ON UPDATE CASCADE,
    currency_to VARCHAR NOT NULL REFERENCES currencies ON DELETE CASCADE ON UPDATE CASCADE,
    exchange_rate FLOAT NOT NULL,
    markup_rate FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(currency_from, currency_to)
);

CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_wallets(
    user_id INTEGER REFERENCES users ON DELETE CASCADE ON UPDATE CASCADE,
    currency_code VARCHAR REFERENCES currencies ON DELETE CASCADE ON UPDATE CASCADE,
    balance FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, currency_code)
);

CREATE TABLE IF NOT EXISTS exchange_rate_offers(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER REFERENCES users ON DELETE CASCADE ON UPDATE CASCADE,
    currency_from VARCHAR NOT NULL REFERENCES currencies ON DELETE CASCADE ON UPDATE CASCADE,
    currency_to VARCHAR NOT NULL REFERENCES currencies ON DELETE CASCADE ON UPDATE CASCADE,
    exchange_rate FLOAT NOT NULL,
    offer_created_at TIMESTAMP NOT NULL,
    offer_expires_at INT NOT NULL
);

INSERT INTO currencies (code) VALUES ('TRY');
INSERT INTO currencies (code) VALUES ('USD');

INSERT INTO currency_exchange_values (currency_from, currency_to, exchange_rate, markup_rate)
VALUES ('TRY', 'USD', 1.40, 0.2);

INSERT INTO currency_exchange_values (currency_from, currency_to, exchange_rate,markup_rate)
VALUES ('USD', 'TRY', 16.00, 0.3);

INSERT INTO users (username, password) VALUES ('eneskzlcn', 'eneskzlcn');
INSERT INTO users (username, password) VALUES ('sedatcan', 'sedatcan');

INSERT INTO user_wallets (user_id, currency_code, balance) VALUES (1, 'TRY', 2250.20);
INSERT INTO user_wallets (user_id, currency_code, balance) VALUES (1, 'USD', 340.20);

INSERT INTO user_wallets (user_id, currency_code, balance) VALUES (2, 'USD', 2441.20);
INSERT INTO user_wallets (user_id, currency_code, balance) VALUES (2, 'TRY', 8200.20);

