-- name: DeleteTokenByToken :exec
DELETE FROM tokens
WHERE token = @token;

-- name: GetNewAccessTokenByRefreshToken :one
SELECT new_access_token
FROM tokens
WHERE refresh_token = @refresh_token AND is_valid = true;
