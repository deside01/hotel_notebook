-- +goose Up
CREATE TABLE IF NOT EXISTS guest_records (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    phone VARCHAR(15) NOT NULL,
    room TEXT NOT NULL,
    check_in DATE NOT NULL,
    check_out DATE NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_stay CHECK (check_out > check_in)
);

-- +goose Down
DROP TABLE IF EXISTS guest_records;
