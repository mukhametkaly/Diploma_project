drop table short_inventory;
CREATE TABLE IF NOT EXISTS short_inventory (
                                               id serial PRIMARY KEY,
                                               document_number varchar,
                                               merchant_id varchar NOT NULL,
                                               total_sum   double precision NOT NULL default 0,
                                               created_on timestamp NOT NULL,
                                               updated_on timestamp NOT NULL,
                                               provided_time timestamp ,
                                               status      varchar
);
ALTER TABLE IF EXISTS short_inventory OWNER TO postgres;

drop table short_waybill;
CREATE TABLE IF NOT EXISTS short_waybill (
                                             id serial PRIMARY KEY,
                                             document_number varchar,
                                             merchant_id varchar NOT NULL,
                                             total_sum   double precision NOT NULL default 0,
                                             created_on timestamp NOT NULL,
                                             updated_on timestamp NOT NULL,
                                             provided_time timestamp ,
                                             status      varchar
);
ALTER TABLE IF EXISTS short_waybill OWNER TO postgres;


CREATE TABLE IF NOT EXISTS inventory_product (
                                                 id serial PRIMARY KEY,
                                                 barcode varchar NOT NULL,
                                                 name varchar NOT NULL,
                                                 actual_amount double precision default 0,
                                                 amount double precision default 0,
                                                 inventory_id int,
                                                 purchase_price double precision,
                                                 selling_price double precision,
                                                 total double precision,
                                                 created_on timestamp not null


);
ALTER TABLE IF EXISTS inventory_product OWNER TO postgres;

CREATE TABLE IF NOT EXISTS waybill_product (
                                               id serial PRIMARY KEY,
                                               barcode varchar NOT NULL,
                                               name varchar NOT NULL,
                                               amount double precision default 0,
                                               received_amount double precision default 0,
                                               waybill_id int NOT NULL,
                                               purchase_price double precision,
                                               selling_price double precision,
                                               total double precision,
                                               created_on timestamp not null

);
ALTER TABLE IF EXISTS waybill_product OWNER TO postgres;