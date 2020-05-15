package main

var dbSchema = `
CREATE TABLE public.log (
    id integer NOT NULL,
    host text,
    message text,
    "timestamp" bigint,
    document tsvector,
    source text,
    severity integer
);
CREATE SEQUENCE public.log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.log_id_seq OWNED BY public.log.id;
ALTER TABLE ONLY public.log ALTER COLUMN id SET DEFAULT nextval('public.log_id_seq'::regclass);
ALTER TABLE ONLY public.log
    ADD CONSTRAINT log_pkey PRIMARY KEY (id);
CREATE INDEX idx_fts_search ON public.log USING gin (document);
`