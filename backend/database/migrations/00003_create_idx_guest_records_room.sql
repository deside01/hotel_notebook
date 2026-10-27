-- +goose Up
CREATE INDEX IF NOT EXISTS idx_guest_records_room 
ON guest_records (lower(room) text_pattern_ops);

-- +goose Down
DROP INDEX IF EXISTS idx_guest_records_room;
