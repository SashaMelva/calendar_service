CREATE TABLE events (
     id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    "time" time without time zone,
    date date NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    CONSTRAINT events_pkey PRIMARY KEY (id)
);