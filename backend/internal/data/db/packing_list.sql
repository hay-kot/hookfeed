-- name: PackingListByID :one
SELECT
    *
FROM
    packing_lists
WHERE
    id = $1;

-- name: PackingListCreate :one
INSERT INTO
    packing_lists (user_id, name, description, due_date, days, tags, status)
VALUES
    ($1, $2, $3, $4, $5, $6, 'not-started') RETURNING *;

-- name: PackingListUpdate :one
UPDATE
    packing_lists
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    due_date = COALESCE(sqlc.narg('due_date'), due_date),
    status = COALESCE(sqlc.narg('status'), status),
    days = COALESCE(sqlc.narg('days'), days),
    tags = CASE
        WHEN @update_tags::bool THEN @tags
        ELSE tags
    END
WHERE
    user_id = sqlc.arg('user_id')
    AND id = sqlc.arg('id') RETURNING *;

-- name: PackingListGetAllByUser :many
SELECT
    sqlc.embed(pl),
    COALESCE(item_counts.item_count, 0) AS item_count,
    COALESCE(item_counts.packed_count, 0) AS packed_count
FROM
    packing_lists pl
    LEFT JOIN (
        SELECT
            packing_list_id,
            COUNT(*) AS item_count,
            COUNT(*) FILTER (WHERE is_packed = TRUE) AS packed_count
        FROM
            packing_list_items
        GROUP BY
            packing_list_id
    ) item_counts ON pl.id = item_counts.packing_list_id
WHERE
    user_id = $1
ORDER BY
    -- Add this to the end of the order by clause so that items are ordered
    -- consistently within an ordered set for consistent results between pages
    pl.due_date ASC NULLS FIRST,

    -- For created_at
    CASE
        WHEN sqlc.arg('order_by') :: text = 'created_at:asc' THEN pl.created_at
    END ASC NULLS LAST,
    CASE
        WHEN sqlc.arg('order_by') = 'created_at:desc' THEN pl.created_at
    END DESC NULLS LAST
LIMIT
    sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: PackingListGetAllByUserCount :one
SELECT
    COUNT(*)
FROM
    packing_lists
WHERE
    user_id = $1;

-- name: PackingListDeleteByID :exec
DELETE FROM
    packing_lists
WHERE
    id = $1;
