-- name: FeedMessageByID :one
SELECT
    sqlc.embed(feed_messages_view)
FROM
    feed_messages_view
WHERE
    id = $1;

-- name: FeedMessageCreate :one
INSERT INTO feed_messages (
    feed_slug,
    raw_request,
    raw_headers,
    raw_query_params,
    title,
    message,
    priority,
    logs,
    metadata,
    state,
    received_at,
    processed_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING id, feed_slug, raw_request, raw_headers, raw_query_params, title, message, priority, logs, metadata, state, state_changed_at, received_at, processed_at, created_at, updated_at;

-- name: FeedMessageGetAll :many
SELECT
    sqlc.embed(feed_messages_view)
FROM
    feed_messages_view
ORDER BY
    -- For received_at
    CASE
        WHEN sqlc.arg('order_by') :: text = 'received_at:asc' THEN received_at
    END ASC NULLS LAST,
    CASE
        WHEN sqlc.arg('order_by') = 'received_at:desc' THEN received_at
    END DESC NULLS LAST,
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
    feed_messages_view;

-- name: FeedMessagesByFeedSlug :many
SELECT
    sqlc.embed(feed_messages_view)
FROM
    feed_messages_view
WHERE
    feed_slug = $1
    AND (sqlc.narg('priority')::integer IS NULL OR priority = sqlc.narg('priority'))
    AND (sqlc.narg('state')::text IS NULL OR state = sqlc.narg('state'))
    AND (sqlc.narg('since')::timestamp IS NULL OR received_at >= sqlc.narg('since'))
    AND (sqlc.narg('until')::timestamp IS NULL OR received_at <= sqlc.narg('until'))
ORDER BY
    received_at DESC,
    id DESC
LIMIT
    sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: FeedMessagesByFeedSlugCount :one
SELECT
    COUNT(*)
FROM
    feed_messages_view
WHERE
    feed_slug = $1
    AND (sqlc.narg('priority')::integer IS NULL OR priority = sqlc.narg('priority'))
    AND (sqlc.narg('state')::text IS NULL OR state = sqlc.narg('state'))
    AND (sqlc.narg('since')::timestamp IS NULL OR received_at >= sqlc.narg('since'))
    AND (sqlc.narg('until')::timestamp IS NULL OR received_at <= sqlc.narg('until'));

-- name: FeedMessageUpdateState :one
UPDATE feed_messages
SET
    state = $2,
    state_changed_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING id, feed_slug, raw_request, raw_headers, raw_query_params, title, message, priority, logs, metadata, state, state_changed_at, received_at, processed_at, created_at, updated_at;

-- name: FeedMessageDeleteByID :exec
DELETE FROM
    feed_messages
WHERE
    id = $1;

-- name: FeedMessageBulkUpdateState :exec
UPDATE feed_messages
SET
    state = $2,
    state_changed_at = CURRENT_TIMESTAMP
WHERE
    id = ANY($1::uuid[]);

-- name: FeedMessageBulkDelete :one
DELETE FROM
    feed_messages
WHERE
    id = ANY(sqlc.arg('message_ids')::uuid[])
RETURNING COUNT(*);

-- name: FeedMessageBulkDeleteByFilter :one
DELETE FROM
    feed_messages
WHERE
    feed_slug = $1
    AND (sqlc.narg('priority')::integer IS NULL OR priority = sqlc.narg('priority'))
    AND (sqlc.narg('older_than')::timestamp IS NULL OR received_at < sqlc.narg('older_than'))
RETURNING COUNT(*);

-- name: FeedMessageSearch :many
SELECT
    sqlc.embed(v)
FROM
    feed_messages_view v
    INNER JOIN feed_messages fm ON v.id = fm.id
WHERE
    (sqlc.narg('feed_slug')::text IS NULL OR v.feed_slug = sqlc.narg('feed_slug'))
    AND (sqlc.narg('priority')::integer IS NULL OR v.priority = sqlc.narg('priority'))
    AND (sqlc.narg('state')::text IS NULL OR v.state = sqlc.narg('state'))
    AND (sqlc.narg('since')::timestamp IS NULL OR v.received_at >= sqlc.narg('since'))
    AND (sqlc.narg('until')::timestamp IS NULL OR v.received_at <= sqlc.narg('until'))
    AND (sqlc.narg('query')::text IS NULL OR fm.search_vector @@ plainto_tsquery('english', sqlc.narg('query')))
ORDER BY
    v.received_at DESC,
    v.id DESC
LIMIT
    sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: FeedMessageSearchCount :one
SELECT
    COUNT(*)
FROM
    feed_messages_view v
    INNER JOIN feed_messages fm ON v.id = fm.id
WHERE
    (sqlc.narg('feed_slug')::text IS NULL OR v.feed_slug = sqlc.narg('feed_slug'))
    AND (sqlc.narg('priority')::integer IS NULL OR v.priority = sqlc.narg('priority'))
    AND (sqlc.narg('state')::text IS NULL OR v.state = sqlc.narg('state'))
    AND (sqlc.narg('since')::timestamp IS NULL OR v.received_at >= sqlc.narg('since'))
    AND (sqlc.narg('until')::timestamp IS NULL OR v.received_at <= sqlc.narg('until'))
    AND (sqlc.narg('query')::text IS NULL OR fm.search_vector @@ plainto_tsquery('english', sqlc.narg('query')));

-- name: FeedMessageDeleteOldByCount :exec
DELETE FROM feed_messages fm
WHERE fm.feed_slug = $1
AND fm.id IN (
    SELECT id
    FROM feed_messages
    WHERE feed_slug = $1
    ORDER BY received_at DESC
    OFFSET $2
);

-- name: FeedMessageDeleteOldByAge :exec
DELETE FROM feed_messages
WHERE
    feed_slug = $1
    AND received_at < NOW() - INTERVAL '1 day' * $2;
