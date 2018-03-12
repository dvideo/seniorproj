CREATE DATABASE allofusdbmysql2;
CREATE TABLE allofusdbmysql2.userTable (
  Userid INT NOT NULL AUTO_INCREMENT,
  fName VARCHAR(45) NULL,
  lName VARCHAR(45) NULL,
  Username VARCHAR(45) NULL COMMENT 'Optional ',
  Password VARCHAR(100) NULL,
  MaritalStatus VARCHAR(45) NULL,
  DateOfBirth INT NULL,
  Email VARCHAR(45) NULL,
  Privacy VARCHAR(45) NULL,
  TimeAccountCreated INT NULL,
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
  
  
INSERT INTO allofusdbmysql2.UserTable (Userid,fName,lName,Username,Password,MaritalStatus,DateOfBirth,Email,Privacy,TimeAccountCreated)
VALUES (00001,'Joshua', 'Gitter', 'jgitter', 'password', 'Single', 06151996, 'joshuacod4@yahoo.com', 'Private', 1200);

INSERT INTO allofusdbmysql2.userLocation(Userid,Location)
VAlUES(00001,'NY');

INSERT INTO allofusdbmysql2.userDevice(Userid,Device)
VAlUES(00001,'Computer');

INSERT INTO allofusdbmysql2.userIPAddress(Userid,IP,OS,Browser)
VAlUES(00001,'.127.0.0.1', 'Mac', 'Chrome');

select * from allofusdbmysql2.UserTable;