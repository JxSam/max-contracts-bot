UPDATE users
SET name = $2
WHERE chat_id = $1;
