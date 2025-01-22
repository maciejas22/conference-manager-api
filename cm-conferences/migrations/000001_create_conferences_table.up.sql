CREATE TABLE public.conferences (
    id serial4 NOT NULL,
    title text NOT NULL,
    start_date timestamp NOT NULL,
    "location" text NOT NULL,
    additional_info text NULL,
    registration_deadline timestamp NULL,
    participants_limit int4 NULL,
    end_date timestamp NOT NULL,
    website text NULL,
    acronym text NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ticket_price int4 DEFAULT 0 NOT NULL,
    CONSTRAINT conferences_pkey PRIMARY KEY (id)
);
