ALTER TABLE users
    DROP COLUMN IF EXISTS ban_reason,
    DROP COLUMN IF EXISTS banned_until,
    DROP COLUMN IF EXISTS banned_at,
    DROP COLUMN IF EXISTS is_banned;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_attribute a
        JOIN pg_type t ON a.atttypid = t.oid
        WHERE t.typname = 'isbanned'
    ) THEN
        DROP TYPE isbanned;
    END IF;
END$$;