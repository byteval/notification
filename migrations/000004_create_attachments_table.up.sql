CREATE TABLE notification_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_notification
      FOREIGN KEY (notification_id) 
      REFERENCES notifications(id)
      ON DELETE CASCADE
);

CREATE INDEX idx_notification_attachments_notification_id 
ON notification_attachments(notification_id);