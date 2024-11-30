CREATE TABLE IF NOT EXISTS public.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS public.sessions (
    session_id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    expires_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    last_accessed_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT sessions_pkey PRIMARY KEY (session_id),
    CONSTRAINT sessions_user_id_key UNIQUE (user_id),
    CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

