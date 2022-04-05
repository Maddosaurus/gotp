CREATE DATABASE pallas ENCODING = 'UTF8' LOCALE = 'en_US.utf8';


ALTER DATABASE pallas OWNER TO postgres;

\connect pallas

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


CREATE TABLE public.pallas (
    uuid character varying NOT NULL COLLATE pg_catalog."C.UTF-8",
    otptype integer NOT NULL,
    name character varying NOT NULL COLLATE pg_catalog."C.UTF-8",
    secret_token character varying NOT NULL COLLATE pg_catalog."C.UTF-8",
    counter bigint DEFAULT 0 NOT NULL,
    update_time timestamp with time zone NOT NULL
);


ALTER TABLE public.pallas OWNER TO postgres;


COPY public.pallas (uuid, otptype, name, secret_token, counter, update_time) FROM stdin;
38518e4a-0b71-4d85-b925-1abdf3b56b03	0	Site1	JBSWY3DPEHPK3PX3	0	1970-04-01 15:07:21.66354+00
d15e9e34-6e9e-4fcd-b795-ec848430b9c4	0	Twertch	JBSWY3DPEHPK3PX4	0	1970-04-02 15:07:21.66354+00
dde17ab2-df84-432d-9a6d-dab3a3476d9b	1	CusSite	4S62BZNFXXSZLCRO	1	1970-04-03 15:07:21.66354+00
\.


ALTER TABLE ONLY public.pallas
    ADD CONSTRAINT pallas_pk PRIMARY KEY (uuid);
