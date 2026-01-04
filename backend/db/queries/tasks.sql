-- name: CreateTask :exec
INSERT INTO tasks (title, user_id, content)
VALUES ($1, $2, $3);

-- name: ListTasksByUser :many
SELECT id, title, content, created_at, updated_at
FROM tasks
WHERE user_id = $1
ORDER BY id;

-- name: UpdateTask :execrows
UPDATE tasks
SET title = $1,
    content = $2,
    updated_at = NOW()
WHERE id = $3 AND user_id = $4;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;
