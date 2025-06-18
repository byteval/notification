CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    layout_id UUID NOT NULL REFERENCES layouts(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    data JSONB,
    channels VARCHAR(20)[] NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    receiver VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL
);

CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_receiver ON notifications(receiver);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);