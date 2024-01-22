CREATE TABLE todo (
    id SERIAL PRIMARY KEY,
    user_id UUID NULL,
    text TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create a trigger to update the updated_at column on every update
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_updated_at
BEFORE UPDATE ON todo
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
