CREATE TABLE public.conference_organizers (
    user_id int4 NOT NULL,
    conference_id int4 NOT NULL,
    joined_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT conference_organizers_pkey PRIMARY KEY (user_id, conference_id),
    CONSTRAINT conference_organizers_conference_id_fkey FOREIGN KEY (conference_id) REFERENCES public.conferences(id) ON DELETE CASCADE ON UPDATE CASCADE
);
