INSERT INTO users_contracts (user_id, contract_id, active)
VALUES ($1, $2, true)
ON CONFLICT (user_id, contract_id)
DO UPDATE SET active = true
WHERE users_contracts.active = false;