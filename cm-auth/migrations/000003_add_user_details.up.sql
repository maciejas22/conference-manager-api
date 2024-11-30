ALTER TABLE public.users
ADD COLUMN name VARCHAR(255),
ADD COLUMN surname VARCHAR(255),
ADD COLUMN username VARCHAR(255),
ADD COLUMN role VARCHAR(255),
ADD COLUMN stripe_account_id VARCHAR(255),
ADD CONSTRAINT users_role_check CHECK ((role = ANY (ARRAY['Organizer'::text, 'Participant'::text])));

