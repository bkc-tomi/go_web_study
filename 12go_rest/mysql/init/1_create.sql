CREATE DATABASE golang_db;
use golang_db;

CREATE TABLE users (
    id INT(11) AUTO_INCREMENT NOT NULL,
    name VARCHAR(250) NOT NULL,
    profile VARCHAR(250),
    date_of_birth VARCHAR(250),
    create_at TIMESTAMP,
    update_at TIMESTAMP,
    PRIMARY KEY (id)
);

INSERT INTO users (name, profile, create_at, update_at) VALUES ("Subaru", "エミリアたんマジ天使！", NOW(), NOW());
INSERT INTO users (name, profile, create_at, update_at) VALUES ("Emilia", "もう、スバルのオタンコナス！", NOW(), NOW());
INSERT INTO users (name, profile, create_at, update_at) VALUES ("Ram", "いいえお客様、きっと生まれて来たのが間違いだわ", NOW(), NOW());
INSERT INTO users (name, profile, create_at, update_at) VALUES ("Rem", "はい、スバルくんのレムです。", NOW(), NOW());
INSERT INTO users (name, profile, create_at, update_at) VALUES ("Roswaal", "君は私になーぁにを望むのかな？", NOW(), NOW());