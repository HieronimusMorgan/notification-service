CREATE TABLE notifications
(
    id             SERIAL PRIMARY KEY,
    target_token   TEXT NOT NULL,
    title          TEXT NOT NULL,
    body           TEXT NOT NULL,
    platform       TEXT NOT NULL,               -- 'android' or 'web'
    priority       TEXT      DEFAULT 'high',    -- 'high', 'normal'
    status         TEXT      DEFAULT 'pending', -- 'pending', 'sent', 'failed'
    service_source TEXT NOT NULL,               -- e.g. 'auth', 'order'
    event_type     TEXT NOT NULL,               -- e.g. 'asset_updated'
    payload        TEXT,                        -- raw JSON string
    color          TEXT      DEFAULT '#000000',
    click_action   TEXT      DEFAULT 'OPEN_APP',
    icon           TEXT      DEFAULT 'default',
    sound          TEXT      DEFAULT 'default',
    retry_count    INTEGER   DEFAULT 0,
    last_error     TEXT,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at        TIMESTAMP
);

-- Recommended indexes for performance
CREATE INDEX idx_notifications_target_token ON notifications (target_token);
CREATE INDEX idx_notifications_status ON notifications (status);
CREATE INDEX idx_notifications_platform ON notifications (platform);
CREATE INDEX idx_notifications_service_source ON notifications (service_source);
CREATE INDEX idx_notifications_event_type ON notifications (event_type);