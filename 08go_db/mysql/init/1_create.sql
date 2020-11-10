CREATE DATABASE golang_db;
use golang_db;

CREATE TABLE talks (
    id INT(11) AUTO_INCREMENT NOT NULL,
    talk VARCHAR(250) NOT NULL,
    creat_at TIMESTAMP,
    update_at TIMESTAMP,
    delete_at TIMESTAMP,
    PRIMARY KEY (id)
);
