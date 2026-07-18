DROP TABLE IF EXISTS notification_settings;
DROP INDEX IF EXISTS idx_notification_settings_user_id;
DROP TABLE IF EXISTS notifications;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_notifications_user_id_created_at;
DROP TABLE IF EXISTS notification_channel_verifications;
DROP INDEX IF EXISTS idx_notification_channel_verifications_user_id;
DROP INDEX IF EXISTS idx_notification_channel_verifications_channel_code;