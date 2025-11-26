-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    question_id Integer NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS answers;
-- +goose StatementEnd
