-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    packing_lists
ADD
    COLUMN STATUS TEXT NOT NULL DEFAULT 'in-progress',
ADD
    COLUMN due_date DATE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    packing_lists DROP COLUMN STATUS,
    DROP COLUMN due_date;

-- +goose StatementEnd
