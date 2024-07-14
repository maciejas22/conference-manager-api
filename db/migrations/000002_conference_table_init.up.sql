CREATE TABLE conferences (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  date TIMESTAMP NOT NULL,
  location VARCHAR(255) NOT NULL,
  additional_info TEXT
);

CREATE TABLE agenda (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conference_id UUID NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  event TEXT NOT NULL,
  speaker VARCHAR(255) NOT NULL,
  FOREIGN KEY (conference_id) REFERENCES conferences(id) ON DELETE CASCADE
);

CREATE TABLE conference_organizers (
    user_id UUID NOT NULL,
    conference_id UUID NOT NULL,
    PRIMARY KEY (user_id, conference_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (conference_id) REFERENCES conferences(id) ON DELETE CASCADE
);

CREATE TABLE conference_participants (
    user_id UUID NOT NULL,
    conference_id UUID NOT NULL,
    PRIMARY KEY (user_id, conference_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (conference_id) REFERENCES conferences(id) ON DELETE CASCADE
);