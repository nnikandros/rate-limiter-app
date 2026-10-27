-- name: GetAllClients :many
SELECT client_id, rps, burst FROM clients;