-- +goose Up
CREATE TABLE IF NOT EXISTS template_meta (
    id BIGSERIAL PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO template_meta (key, value)
VALUES ('schema', 'ok')
ON CONFLICT (key) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS template_meta;
