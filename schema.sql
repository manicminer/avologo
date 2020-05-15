--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2
-- Dumped by pg_dump version 12.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'SQL_ASCII';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: log; Type: TABLE; Schema: public; 
--

CREATE TABLE public.log (
    id integer NOT NULL,
    host text,
    message text,
    "timestamp" bigint,
    document tsvector,
    source text,
    severity integer
);


ALTER TABLE public.log OWNER TO brad;

--
-- Name: log_id_seq; Type: SEQUENCE; Schema: public; 
--

CREATE SEQUENCE public.log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.log_id_seq OWNER TO brad;

--
-- Name: log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; 
--

ALTER SEQUENCE public.log_id_seq OWNED BY public.log.id;


--
-- Name: log id; Type: DEFAULT; Schema: public; 
--

ALTER TABLE ONLY public.log ALTER COLUMN id SET DEFAULT nextval('public.log_id_seq'::regclass);


--
-- Name: log log_pkey; Type: CONSTRAINT; Schema: public; 
--

ALTER TABLE ONLY public.log
    ADD CONSTRAINT log_pkey PRIMARY KEY (id);


--
-- Name: idx_fts_search; Type: INDEX; Schema: public; 
--

CREATE INDEX idx_fts_search ON public.log USING gin (document);


--
-- PostgreSQL database dump complete
--
