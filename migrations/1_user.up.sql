CREATE TABLE IF NOT EXISTS public.user (
    id UUID PRIMARY KEY UNIQUE DEFAULT (gen_random_uuid()),
    email TEXT NOT NULL UNIQUE,
    refresh_token TEXT UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_user_id ON public.user(id);