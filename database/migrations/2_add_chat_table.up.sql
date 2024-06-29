DROP TABLE IF EXISTS chats;
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    to_user_id INT NOT NULL,
    from_user_id INT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE chats ADD CONSTRAINT fk_to_user_id FOREIGN KEY (to_user_id) REFERENCES users(user_id);
ALTER TABLE chats ADD CONSTRAINT fk_from_user_id FOREIGN KEY (from_user_id) REFERENCES users(user_id);

