CREATE TABLE layouts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    subject     VARCHAR(200) NOT NULL,
    body        TEXT NOT NULL,
    type        VARCHAR(20) NOT NULL,
    variables   TEXT[],
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version     INTEGER DEFAULT 1
);
