CREATE TABLE public.sections (
    id serial4 NOT NULL,
    terms_of_service_id int4 NOT NULL,
    title text NOT NULL,
    content text NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT sections_pkey PRIMARY KEY (id),
    CONSTRAINT sections_terms_of_service_id_fkey FOREIGN KEY (terms_of_service_id) REFERENCES public.terms_of_service(id) ON DELETE CASCADE ON UPDATE CASCADE
);
