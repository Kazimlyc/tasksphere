-- name: CreateTask :exec
INSERT INTO tasks (title, user_id, content, status)
VALUES ($1, $2, $3, $4);

-- name: ListTasksByUser :many
SELECT id, title, content, status, created_at, updated_at
FROM tasks
WHERE user_id = $1
ORDER BY id;

-- name: UpdateTask :execrows
UPDATE tasks
SET title = $1,
    content = $2,
    status = $3,
    updated_at = NOW()
WHERE id = $4 AND user_id = $5;

-- name: GetTaskStatusByID :one
SELECT status
FROM tasks
WHERE id = $1 AND user_id = $2;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;
