CREATE TABLE member (
    id INTEGER NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    birthday VARCHAR(5),
    blood_type CHAR(2),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8;

CREATE TABLE hobby (
    id INTEGER NOT NULL AUTO_INCREMENT,
    member_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id),
    CONSTRAINT FOREIGN KEY (member_id) REFERENCES member(id)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8;

