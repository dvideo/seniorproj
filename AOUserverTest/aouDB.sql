CREATE DATABASE allofusdbmysql2;
drop table allofusdbmysql2.userTable;
CREATE TABLE allofusdbmysql2.userTable (
  Userid INT NOT NULL AUTO_INCREMENT,
  fName VARCHAR(45) NULL,
  lName VARCHAR(45) NULL,
  Username VARCHAR(45) NULL COMMENT 'Optional ',
  Password VARCHAR(100) NULL,
  MaritalStatus VARCHAR(45) NULL,
  DateOfBirth INT NULL,
  Email VARCHAR(45) NULL,
  Location VARCHAR(45) NULL,
  Device VARCHAR(45) NULL,
  IP VARCHAR(45) NULL,
  Privacy VARCHAR(45) NULL,
  TimeAccountCreated INT NULL,
  OS VARCHAR(45) NULL,
  Browser VARCHAR(45) NULL,
  PRIMARY KEY (Userid)
  );

CREATE TABLE allofusdbmysql2.userLocation (
  Userid INT NOT NULL AUTO_INCREMENT,
  Location VARCHAR(45) NULL,

  PRIMARY KEY (Userid)
  );

CREATE TABLE allofusdbmysql2.userDevice (
  Userid INT NOT NULL AUTO_INCREMENT,
  Device VARCHAR(45) NULL,

  PRIMARY KEY (Userid)
  );

CREATE TABLE allofusdbmysql2.userIPAddress (
  Userid INT NOT NULL AUTO_INCREMENT,
  IP VARCHAR(45) NULL,
  OS VARCHAR(45) NULL,
  Browser VARCHAR(45) NULL,

  PRIMARY KEY (Userid)
  );

CREATE TABLE allofusdbmysql2.profileTable (
  Userid INT NOT NULL AUTO_INCREMENT,
  Profilepic VARCHAR(45) NULL,
  Cover1 VARCHAR(255) NULL,
  Cover2 VARCHAR(255) NULL,
  Cover3 VARCHAR(255) NULL,
  Cover4 VARCHAR(255) NULL,
  Cover5 VARCHAR(255) NULL,
  Cover6 VARCHAR(255) NULL,
  Username VARCHAR(45) NULL,
  Password VARCHAR(45) NULL,
  Privacy VARCHAR(45) NULL,
  TimeAccountCreated INT NULL,
  PRIMARY KEY (Userid)
  );

CREATE TABLE allofusdbmysql2.userPost (
  PostID INT NOT NULL AUTO_INCREMENT,
  Userid INT NOT NULL,
  Author VARCHAR(255) NULL,
  Recipient VARCHAR(255) NULL,
  Photo VARCHAR(255) NULL,
  Video VARCHAR(255) NULL,
  Meme1 VARCHAR(255) NULL,
  Meme2 VARCHAR(255) NULL,
  TimeAccountCreated INT NULL,
  PRIMARY KEY (PostID)
  );

CREATE TABLE allofusdbmysql2.statPost (
  PostID INT NOT NULL,
  StatAvg INT NULL,
  TimeCreated INT NULL,
  PRIMARY KEY (PostID)
  );

CREATE TABLE allofusdbmysql2.stats (
  PostID INT NOT NULL,
  UseridStat VARCHAR(255) NULL,
  StatValue VARCHAR(255) NULL,
  TimeCreated INT NULL,
  PRIMARY KEY (PostID)
  );
  
  
INSERT INTO allofusdbmysql2.UserTable (Userid,fName,lName,Username,Password,MaritalStatus,DateOfBirth,Email,Privacy,TimeAccountCreated)
VALUES (00001,'Joshua', 'Gitter', 'jgitter', 'password', 'Single', 06151996, 'joshuacod4@yahoo.com', 'Private', 1200);

INSERT INTO allofusdbmysql2.userLocation(Userid,Location)
VAlUES(00001,'NY');

INSERT INTO allofusdbmysql2.userDevice(Userid,Device)
VAlUES(00001,'Computer');

INSERT INTO allofusdbmysql2.userIPAddress(Userid,IP,OS,Browser)
VAlUES(00001,'.127.0.0.1', 'Mac', 'Chrome');

select * from allofusdbmysql2.UserTable;
