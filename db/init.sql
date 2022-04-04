CREATE DATABASE gotp ENCODING = 'UTF8' LOCALE = 'en_US.utf8';


ALTER DATABASE gotp OWNER TO postgres;

\connect gotp

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;


CREATE TABLE public.gotp (
    uuid character varying NOT NULL COLLATE pg_catalog."C.UTF-8",
    otptype integer NOT NULL,
    name character varying NOT NULL COLLATE pg_catalog."C.UTF-8",
    secret_token character varying NOT NULL COLLATE pg_catalog."C.UTF-8",
    counter bigint DEFAULT 0 NOT NULL,
    update_time timestamp with time zone NOT NULL
);


ALTER TABLE public.gotp OWNER TO postgres;


COPY public.gotp (uuid, otptype, name, secret_token, counter, update_time) FROM stdin;
1234567dfcg	1	CustomSite	4S62BZNFXXSZLCRO	1	1970-04-04 15:07:21.66354+00
\.


ALTER TABLE ONLY public.gotp
    ADD CONSTRAINT gotp_pk PRIMARY KEY (uuid);
