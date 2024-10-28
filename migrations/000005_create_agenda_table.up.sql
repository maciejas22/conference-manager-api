CREATE TABLE public.agenda (
    id serial4 NOT NULL,
    conference_id int4 NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    event text NOT NULL,
    speaker text NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT agenda_pkey PRIMARY KEY (id),
    CONSTRAINT agenda_conference_id_fkey FOREIGN KEY (conference_id) REFERENCES public.conferences(id) ON DELETE CASCADE ON UPDATE CASCADE
);
