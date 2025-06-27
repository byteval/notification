CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    layout_id UUID,
    title VARCHAR(255) NOT NULL,
    data JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
);

CREATE INDEX idx_notifications_created_at ON notifications(created_at);