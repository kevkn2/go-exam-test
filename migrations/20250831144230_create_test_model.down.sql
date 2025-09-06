-- Drop questions first (because it depends on online_tests)
DROP INDEX IF EXISTS idx_questions_online_test_id;
DROP INDEX IF EXISTS idx_questions_deleted_at;
DROP TABLE IF EXISTS questions;

-- Then drop online_tests
DROP INDEX IF EXISTS idx_online_tests_deleted_at;
DROP TABLE IF EXISTS online_tests;
