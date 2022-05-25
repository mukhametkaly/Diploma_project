CREATE TABLE IF NOT EXISTS shopping_carts (
    id serial PRIMARY KEY,
    merchant_id varchar NOT NULL,
    total_sum   double precision NOT NULL,
    created_on timestamp NOT NULL,
    provided_time timestamp NOT NULL,
    status      varchar
);
ALTER TABLE IF EXISTS shopping_carts OWNER TO postgres;

CREATE TABLE IF NOT EXISTS shopping_cart_products (
    id serial PRIMARY KEY,
    barcode varchar NOT NULL,
    name varchar NOT NULL,
    amount double precision default 0,
    shopping_cart_id int,
    purchase_price double precision,
    selling_price double precision,
    total double precision,
    created_on timestamp not null
);
ALTER TABLE IF EXISTS shopping_cart_products OWNER TO postgres;

