CREATE TABLE IF NOT EXISTS score (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    score INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE MATERIALIZED VIEW score_leaderboard AS
    SELECT
        user_id,
        SUM(score) AS user_score,
        ROW_NUMBER() OVER (ORDER BY SUM(score) DESC, MIN(created_at)) AS position
    FROM
        score
    GROUP by
        user_id
    ORDER by
        position;

CREATE INDEX score_leaderboard_user_id_idx ON score_leaderboard(user_id);