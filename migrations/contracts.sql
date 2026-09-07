-- Drop table

-- DROP TABLE contracts;

CREATE TABLE contracts (
	id int8 NOT NULL,
	link varchar NULL,
	created_at timestamp NULL,
	description varchar NULL,
	name_product varchar NULL,
	CONSTRAINT contracts_pkey PRIMARY KEY (id)
);
