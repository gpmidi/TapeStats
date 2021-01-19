-- Need for UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Auto modified
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.modified = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- User accounts
CREATE TABLE accounts
(
    -- Our info
    id       UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Whens (auto set/updated)
    created  TIMESTAMPTZ NOT NULL             DEFAULT NOW(),
    modified TIMESTAMPTZ NOT NULL             DEFAULT NOW(),

    -- Auth info
    salt     VARCHAR(1024),
    hashed   VARCHAR(1024)
);

-- Tie modified trigger to accounts
CREATE TRIGGER trigger_accounts_set_modified
    BEFORE UPDATE
    ON accounts
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
