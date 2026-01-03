-- name: CreateTask :exec
INSERT INTO tasks (title, user_id)
VALUES ($1, $2);

-- name: ListTasksByUser :many
SELECT id, title, created_at, updated_at
FROM tasks
WHERE user_id = $1
ORDER BY id;

-- name: UpdateTask :execrows
UPDATE tasks
SET title = $1,
    updated_at = NOW()
WHERE id = $2 AND user_id = $3;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;
