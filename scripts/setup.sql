CREATE TABLE events (
     id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    date_time_start timestamp with time zone NOT NULL,
    date_time_end timestamp with time zone,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    CONSTRAINT events_pkey PRIMARY KEY (id)
);