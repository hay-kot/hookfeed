-- name: {{ .Computed.domain_var }}ByID :one
SELECT
    *
FROM
    {{ .Scaffold.sql_table }}
WHERE
    id = $1;

-- name: {{ .Computed.domain_var }}GetAll :many
SELECT
    *
FROM
    {{ .Scaffold.sql_table }}
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
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: {{ .Computed.domain_var }}GetAllCount :one
SELECT
    COUNT(*)
FROM
    {{ .Scaffold.sql_table }};

-- name: {{ .Computed.domain_var }}DeleteByID :exec
DELETE FROM
    {{ .Scaffold.sql_table }}
WHERE
    id = $1;
{{ if .Scaffold.user_relation }}
-- name: {{ .Computed.domain_var }}GetAllByUser :many
SELECT
    *
FROM
    {{ .Scaffold.sql_table }}
WHERE
    user_id = $1
ORDER BY
    created_at
LIMIT
    $2 OFFSET $3;

-- name: {{ .Computed.domain_var }}GetAllByUserCount :one
SELECT
    COUNT(*)
FROM
    {{ .Scaffold.sql_table }}
WHERE
    user_id = $1;
{{ end }}
