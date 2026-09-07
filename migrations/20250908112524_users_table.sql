-- Drop table

-- DROP TABLE users;

CREATE TABLE users (
	chat_id int8 NOT NULL,
	email varchar(128) NULL DEFAULT NULL::character varying,
	"name" varchar(128) NULL DEFAULT NULL::character varying,
	CONSTRAINT users_pkey PRIMARY KEY (chat_id)
);
