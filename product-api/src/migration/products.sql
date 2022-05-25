CREATE TABLE IF NOT EXISTS products (
        id serial PRIMARY KEY,
        barcode varchar NOT NULL,
        name varchar NOT NULL,
        category_name varchar NOT NULL,
        category_id int,
        merchant_id varchar NOT NULL,
        stock_id int,

        purchase_price double precision,
        selling_price double precision,
        amount double precision,
        margin double precision,
        unit_type varchar

);

ALTER TABLE IF EXISTS products OWNER TO postgres;
