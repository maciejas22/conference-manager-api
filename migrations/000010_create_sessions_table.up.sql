CREATE TABLE public.sessions (
    session_id varchar(255) DEFAULT lower(encode(gen_random_bytes(16), 'hex'::text)) NOT NULL,
    user_id int4 NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    expires_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    last_accessed_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT sessions_pkey PRIMARY KEY (session_id),
    CONSTRAINT sessions_user_id_key UNIQUE (user_id),
    CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
