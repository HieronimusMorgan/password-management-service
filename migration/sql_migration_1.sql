CREATE TABLE password_groups
(
    group_id   SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255)
);
CREATE INDEX idx_password_groups_group_id ON password_groups (group_id);
CREATE INDEX idx_password_groups_user_id ON password_groups (user_id);

CREATE TABLE password_entries
(
    entry_id           SERIAL PRIMARY KEY,
    user_id            INT          NOT NULL,
    group_id           INT          REFERENCES password_groups (group_id) ON DELETE SET NULL,
    title              VARCHAR(255) NOT NULL,
    username           TEXT         NOT NULL,
    encrypted_password TEXT         NOT NULL,
    encrypted_notes    TEXT,
    url                TEXT,
    tags               TEXT[],    -- Optional categorization
    expires_at         TIMESTAMP, -- Optional for temporary credentials
    last_accessed_at   TIMESTAMP,
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by         VARCHAR(255),
    updated_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by         VARCHAR(255),
    deleted_at         TIMESTAMP NULL,
    deleted_by         VARCHAR(255)
);
CREATE INDEX idx_password_entries_entry_id ON password_entries (entry_id);
CREATE INDEX idx_password_entries_user_id ON password_entries (user_id);
CREATE INDEX idx_password_entries_group_id ON password_entries (group_id);

CREATE TABLE password_entry_keys
(
    entry_id                INT PRIMARY KEY REFERENCES password_entries (entry_id) ON DELETE CASCADE,
    encrypted_symmetric_key TEXT NOT NULL
);
CREATE INDEX idx_password_entry_keys_entry_id ON password_entry_keys (entry_id);

CREATE TABLE shared_passwords
(
    share_id                SERIAL PRIMARY KEY,
    entry_id                INT REFERENCES password_entries (entry_id) ON DELETE CASCADE,
    from_user_id            INT  NOT NULL,
    to_user_id              INT  NOT NULL,
    encrypted_symmetric_key TEXT NOT NULL,
    shared_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_share UNIQUE (entry_id, to_user_id)
);
CREATE INDEX idx_shared_passwords_entry_id ON shared_passwords (entry_id);
CREATE INDEX idx_shared_passwords_from_user_id ON shared_passwords (from_user_id);
CREATE INDEX idx_shared_passwords_to_user_id ON shared_passwords (to_user_id);
CREATE INDEX idx_shared_passwords_share_id ON shared_passwords (share_id);

CREATE TABLE password_history
(
    history_id         SERIAL PRIMARY KEY,
    entry_id           INT REFERENCES password_entries (entry_id) ON DELETE CASCADE,
    encrypted_password TEXT NOT NULL,
    changed_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    changed_by         VARCHAR(255)
);
CREATE INDEX idx_password_history_entry_id ON password_history (entry_id);
CREATE INDEX idx_password_history_history_id ON password_history (history_id);