CREATE TABLE public.news (
    id serial4 NOT NULL,
    title text NOT NULL,
    content text NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT news_pkey PRIMARY KEY (id)
);
