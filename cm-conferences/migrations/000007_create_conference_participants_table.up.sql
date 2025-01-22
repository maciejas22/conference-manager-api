CREATE TABLE public.conference_participants (
    user_id int4 NOT NULL,
    conference_id int4 NOT NULL,
    joined_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ticket_id uuid DEFAULT gen_random_uuid() NOT NULL,
    CONSTRAINT conference_participants_pkey PRIMARY KEY (user_id, conference_id),
    CONSTRAINT conference_participants_ticketid_key UNIQUE (ticket_id),
    CONSTRAINT conference_participants_conference_id_fkey FOREIGN KEY (conference_id) REFERENCES public.conferences(id) ON DELETE CASCADE ON UPDATE CASCADE
);
