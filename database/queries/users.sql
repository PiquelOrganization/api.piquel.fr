-- name: AddUser :one
INSERT INTO "users" ("username", "name", "image", "email", "role")
VALUES ($1, $2, $3, $4, $5) RETURNING "id";

-- name: GetUserByUsername :one
SELECT * FROM "users" WHERE "username" = $1;

-- name: GetUserById :one
SELECT * FROM "users" WHERE "id" = $1;

-- name: GetUserByEmail :one
SELECT * FROM "users" WHERE "email" = $1;

-- name: UpdateUser :exec
UPDATE "users" SET "username" = $2, "name" = $3, "image" = $4 WHERE "id" = $1;
