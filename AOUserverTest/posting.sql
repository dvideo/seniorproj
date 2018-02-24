CREATE DATABASE allofusdbmysql2;
CREATE TABLE allofusdbmysql2.userPhotos (
  Userid VARCHAR(255) NULL AUTO_INCREMENT,
  Profilepic VARCHAR(45) NULL,
  Cover1 VARCHAR(255) NULL,
  Cover2 VARCHAR(255) NULL,
  Cover3 VARCHAR(255) NULL,
  Cover4 VARCHAR(255) NULL,
  Cover5 VARCHAR(255) NULL,
  Cover6 VARCHAR(255) NULL,
  TimeAccountCreated INT NULL,
  PRIMARY KEY (Userid)
  );

CREATE TABLE allofusdbmysql2.userPost (
  PostID VARCHAR(255) NULL,
  Author VARCHAR(255) NULL,
  Recipient VARCHAR(255) NULL,
  Photo VARCHAR(255) NULL,
  Video VARCHAR(255) NULL,
  Meme1 VARCHAR(255) NULL,
  Meme2 VARCHAR(255) NULL,
  TimeAccountCreated INT NULL,
  PRIMARY KEY (Userid)
  );

CREATE TABLE allofusdbmysql2.statPost (
  PostID VARCHAR(255) NULL,
  StatAvg INT NULL,
  TimeCreated INT NULL,
  PRIMARY KEY (PostID)
  );

CREATE TABLE allofusdbmysql2.stats (
  PostID VARCHAR(255) NULL,
  UseridStat VARCHAR(255) NULL,
  StatValue VARCHAR(255) NULL,
  TimeCreated INT NULL,
  PRIMARY KEY (PostID)
  );