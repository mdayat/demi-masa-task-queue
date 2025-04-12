-- name: SelectUsers :many
SELECT * FROM "user";

-- name: InsertUserPrayers :copyfrom
INSERT INTO prayer (id, user_id, name, year, month, day)
VALUES ($1, $2, $3, $4, $5, $6);