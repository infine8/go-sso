CREATE TABLE IF NOT EXISTS roles
(
    id          INTEGER PRIMARY KEY,
    name        TEXT    NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_roles
(
    user_id     INT     NOT NULL,
    role_id     INT     NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_id ON user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_role_id ON user_roles (role_id);
