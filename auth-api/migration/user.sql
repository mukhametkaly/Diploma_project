CREATE TABLE IF NOT EXISTS users (
                                        username varchar PRIMARY KEY,
                                        password varchar NOT NULL,
                                        full_name varchar NOT NULL,
                                        salt varchar NOT NULL,
                                        user_role varchar NOT NULL,
                                        IIN varchar ,
                                        mail varchar ,
                                        mobile varchar ,
                                        merchant_id varchar

);

ALTER TABLE IF EXISTS users OWNER TO postgres;
