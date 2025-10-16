-- name: PackingListItemByID :one
SELECT
    *
FROM
    packing_list_items
WHERE
    id = $1;

-- name: PackingListItemCreate :one
INSERT INTO
    packing_list_items (
        packing_list_id,
        name,
        category,
        notes,
        quantity,
        is_packed
    )
VALUES
    ($1, $2, $3, $4, $5, false) RETURNING *;

-- name: PackingListItemUpdate :one
UPDATE
    packing_list_items
SET
    name = COALESCE(sqlc.narg('name'), name),
    category = COALESCE(sqlc.narg('category'), category),
    notes = COALESCE(sqlc.narg('notes'), notes),
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    is_packed = COALESCE(sqlc.narg('is_packed'), is_packed)
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: PackingListItemGetAll :many
SELECT
    *
FROM
    packing_list_items
WHERE
    packing_list_id = sqlc.arg('packing_list_id')
ORDER BY
    created_at DESC,
    id DESC;

-- name: PackingListItemDeleteByID :exec
DELETE FROM
    packing_list_items
WHERE
    id = $1;
