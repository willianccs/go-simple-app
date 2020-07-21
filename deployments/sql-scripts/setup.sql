-- create tables
CREATE TABLE breeds (
    id MEDIUMINT NOT NULL AUTO_INCREMENT,
    name varchar(100),
    origin varchar(50),
    temperament  varchar(100),
    description varchar(1000),
    PRIMARY KEY (id)
);

CREATE TABLE breeds_images (
    id varchar(10),
    img  varchar(100)
);

CREATE TABLE cats_hat (
    id varchar(10),
    img  varchar(100)
);

CREATE TABLE cats_glasses (
    id varchar(10),
    img  varchar(100)
);

-- create the user with privileges
CREATE USER 'topcat'@'%' identified with mysql_native_password by 'Zaq!@wsx34';
GRANT ALL PRIVILEGES ON cats.* TO 'topcat'@'%';
UPDATE mysql.user SET Super_Priv='Y' WHERE user='topcat';
FLUSH PRIVILEGES;