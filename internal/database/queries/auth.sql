-- name: CreateUser :exec
INSERT INTO users (uuid, auth_type, oauth_id)
VALUES (@uuid, @auth_type, @oauth_id);

-- name: GetUserByOAuthID :one
SELECT uuid
FROM users
WHERE auth_type = @auth_type AND oauth_id = @oauth_id;
