CREATE TABLE IF NOT EXISTS songs(
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    lyrics TEXT,
    release_date VARCHAR(255),
    link VARCHAR(255)
    );