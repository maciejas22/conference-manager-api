CREATE TABLE News (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

