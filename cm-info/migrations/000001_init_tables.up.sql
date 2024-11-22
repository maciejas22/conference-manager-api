CREATE TABLE IF NOT EXISTS public.terms_of_service (
    id serial4 NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    introduction text NOT NULL,
    acknowledgement text NOT NULL,
    CONSTRAINT terms_of_service_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.sections (
    id serial4 NOT NULL,
    terms_of_service_id int4 NOT NULL,
    title text NOT NULL,
    content text NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT sections_pkey PRIMARY KEY (id),
    CONSTRAINT sections_terms_of_service_id_fkey FOREIGN KEY (terms_of_service_id) REFERENCES public.terms_of_service(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.subsections (
    id serial4 NOT NULL,
    section_id int4 NOT NULL,
    title text NOT NULL,
    content text NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT subsections_pkey PRIMARY KEY (id),
    CONSTRAINT subsections_section_id_fkey FOREIGN KEY (section_id) REFERENCES public.sections(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.news (
    id serial4 NOT NULL,
    title text NOT NULL,
    content text NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT news_pkey PRIMARY KEY (id)
);
