CREATE TABLE public.terms_of_service (
    id serial4 NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    introduction text NOT NULL,
    acknowledgement text NOT NULL,
    CONSTRAINT terms_of_service_pkey PRIMARY KEY (id)
);
