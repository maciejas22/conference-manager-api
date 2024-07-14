ALTER TABLE conferences
ADD COLUMN registration_deadline TIMESTAMP NULL,
ADD COLUMN participants_limit INT NULL;
