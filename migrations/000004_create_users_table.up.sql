CREATE TABLE public.users (
    id serial4 NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name text NULL,
    surname text NULL,
    username text NULL,
    email text NOT NULL,
    role text NULL,
    password text NOT NULL,
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_role_check CHECK ((role = ANY (ARRAY['Organizer'::text, 'Participant'::text]))),
    CONSTRAINT users_username_key UNIQUE (username)
);
