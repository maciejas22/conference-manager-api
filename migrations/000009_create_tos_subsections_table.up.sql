CREATE TABLE public.subsections (
    id serial4 NOT NULL,
    section_id int4 NOT NULL,
    title text NOT NULL,
    content text NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT subsections_pkey PRIMARY KEY (id),
    CONSTRAINT subsections_section_id_fkey FOREIGN KEY (section_id) REFERENCES public.sections(id) ON DELETE CASCADE ON UPDATE CASCADE
);
