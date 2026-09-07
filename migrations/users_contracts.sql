-- Drop table

-- DROP TABLE users_contracts;

CREATE TABLE users_contracts (
	user_id int8 NOT NULL,
	contract_id int8 NOT NULL,
	active bool NULL DEFAULT true,
	CONSTRAINT users_contracts_pkey PRIMARY KEY (user_id, contract_id),
	CONSTRAINT users_contracts_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.contracts(id),
	CONSTRAINT users_contracts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(chat_id)
);