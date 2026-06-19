-- UP
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (
        (receiver_id IS NOT NULL AND group_id IS NULL) OR 
        (receiver_id IS NULL AND group_id IS NOT NULL)
    ),
    CHECK (sender_id != receiver_id)
);

CREATE INDEX idx_messages_receiver ON messages(receiver_id);
CREATE INDEX idx_messages_group ON messages(group_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);
