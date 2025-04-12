-- name: SelectUsers :many
SELECT * FROM "user";

-- name: InsertUserPrayers :copyfrom
INSERT INTO prayer (id, user_id, name, year, month, day)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUncheckedPrayer :exec
UPDATE prayer SET status = 'missed'
WHERE status = 'pending' AND day <= $1 AND month = $2 AND year = $3;