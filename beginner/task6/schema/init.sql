CREATE TABLE post
(
    id INT PRIMARY KEY ,
    user_id INT,
    title VARCHAR(256),
    body TEXT
);

CREATE TABLE comment
(
    id INT PRIMARY KEY,
    post_id INT,
    name VARCHAR(256),
    email VARCHAR(256),
    body TEXT,
    FOREIGN KEY (post_id)  REFERENCES post(id)
);