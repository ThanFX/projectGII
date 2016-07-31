CREATE TABLE IF NOT EXISTS public.persons
(
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    job_id          INT,
    state           INT NOT NULL,
    characteristic  JSONB
);

CREATE UNIQUE INDEX IF NOT EXISTS persons_id_uindex ON public.persons (id);
COMMENT ON TABLE public.persons IS 'Таблица персонажей';

CREATE TABLE IF NOT EXISTS public.time
(
    id INTEGER DEFAULT 1 NOT NULL,
    real_time INTEGER,
    world_time INTEGER
);

CREATE TABLE public.world_map
(
    chunk_id SERIAL PRIMARY KEY NOT NULL,
    x INT DEFAULT 0 NOT NULL,
    y INT DEFAULT 0 NOT NULL,
    is_explored BOOLEAN DEFAULT FALSE  NOT NULL,
    terrains JSONB DEFAULT '{"terrains":[]}' NOT NULL
);
CREATE UNIQUE INDEX world_map_chunk_id_uindex ON public.world_map (chunk_id);