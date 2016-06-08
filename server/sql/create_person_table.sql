CREATE TABLE IF NOT EXISTS persons
(
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    job_id          INT,
    state           INT NOT NULL,
    characteristic  JSONB
);

CREATE UNIQUE INDEX IF NOT EXISTS persons_id_uindex ON public.persons (id);
COMMENT ON TABLE public.persons IS 'Таблица персонажей'