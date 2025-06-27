-- Drop indexes
DROP INDEX IF EXISTS idx_notification_receivers_notification_id;
DROP INDEX IF EXISTS idx_notification_receivers_status;
DROP INDEX IF EXISTS idx_notification_receivers_email;

-- Finally drop the table
DROP TABLE IF EXISTS notification_receivers;
