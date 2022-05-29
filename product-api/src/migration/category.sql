CREATE TABLE IF NOT EXISTS categories (
      id serial PRIMARY KEY,
      category_name varchar NOT NULL,
      merchant_id varchar NOT NULL,
      description varchar,
      products_count int not null default 0,
      created_on timestamp,
      updated_on timestamp,
      unique (category_name, merchant_id)
);

ALTER TABLE IF EXISTS products OWNER TO postgres;
