CREATE TABLE notification_receivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error TEXT NULL,
    sent_at TIMESTAMP WITH TIME ZONE NULL,
    delivered_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_notification
        FOREIGN KEY (notification_id) 
        REFERENCES notifications(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_notification_receivers_notification_id ON notification_receivers(notification_id);
CREATE INDEX idx_notification_receivers_status ON notification_receivers(status);
CREATE INDEX idx_notification_receivers_email ON notification_receivers(email);
