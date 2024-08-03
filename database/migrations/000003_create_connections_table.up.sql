DROP TABLE IF EXISTS connections;

CREATE TABLE connections (
    user_id_1 INT NOT NULL,
    user_id_2 INT NOT NULL,
    chat_id SERIAL PRIMARY KEY,
    FOREIGN KEY (user_id_1) REFERENCES users (user_id),
    FOREIGN KEY (user_id_2) REFERENCES users (user_id)
);

-- Create an index on the combination of user_id_1 and user_id_2
CREATE INDEX idx_user_ids ON connections (user_id_1, user_id_2);
