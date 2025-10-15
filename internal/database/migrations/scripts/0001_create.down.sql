DROP INDEX IF EXISTS idx_leaderboards_created_at;
DROP INDEX IF EXISTS idx_leaderboards_score;
DROP INDEX IF EXISTS idx_leaderboards_user_id;
DROP INDEX IF EXISTS idx_questions_difficulty;
DROP INDEX IF EXISTS idx_questions_topic;
DROP INDEX IF EXISTS idx_users_last_activity;
DROP INDEX IF EXISTS idx_users_telegram_id;

DROP TABLE IF EXISTS leaderboards;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS users;