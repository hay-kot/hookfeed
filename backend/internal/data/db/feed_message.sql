-- name: FeedMessageByID :one
SELECT
    *
FROM
    feed_messages
WHERE
    id = $1;

-- name: FeedMessageGetAll :many
SELECT
    *
FROM
    feed_messages
ORDER BY
    -- For created_at
    CASE
        WHEN sqlc.arg('order_by') :: text = 'created_at:asc' THEN created_at
    END ASC NULLS LAST,
    CASE
        WHEN sqlc.arg('order_by') = 'created_at:desc' THEN created_at
    END DESC NULLS LAST,
    -- Add this to the end of the order by clause so that items are ordered
    -- consistently within an ordered set for consistent results between pages
    id DESC
LIMIT
    sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: FeedMessageGetAllCount :one
SELECT
    COUNT(*)
FROM
    feed_messages;

-- name: FeedMessageDeleteByID :exec
DELETE FROM
    feed_messages
WHERE
    id = $1;
