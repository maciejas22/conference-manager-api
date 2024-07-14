CREATE TABLE terms_of_service (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    introduction TEXT NOT NULL,
    acknowledgement TEXT NOT NULL
);

CREATE TABLE sections (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    terms_of_service_id UUID NOT NULL REFERENCES terms_of_service(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT
);

CREATE TABLE subsections (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    section_id UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);
