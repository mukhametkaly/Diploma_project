CREATE TABLE IF NOT EXISTS merchants (
    merchant_id varchar PRIMARY KEY,
    merchant_name varchar NOT null,
    ie varchar not null,
    address varchar not null,
    status varchar not null,
    bin varchar not null,
    phone varchar not null,
    email varchar not null,
    created_on timestamp not null,
    updated_on timestamp not null

);

ALTER TABLE IF EXISTS merchants OWNER TO postgres;