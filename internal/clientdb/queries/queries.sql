-- name: GetAllClients :many
SELECT * FROM clients;

-- name: GetRateLimitByClient :one
SELECT rps,burst FROM clients
WHERE client_id=?;