DROP TABLE IF EXISTS chats;
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    user_id_1 INT NOT NULL,
    user_id_2 INT NOT NULL,
    chat_id INT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE chats ADD CONSTRAINT fk_to_user_id FOREIGN KEY (user_id_1) REFERENCES users(user_id);
ALTER TABLE chats ADD CONSTRAINT fk_from_user_id FOREIGN KEY (user_id_2) REFERENCES users(user_id);

CREATE INDEX idx_chat_id ON chats (chat_id);
