package models

const ProfileSchema = `
CREATE TABLE Profile (
    ID int NOT NULL AUTO_INCREMENT,
    LastName varchar(255),
    FirstName varchar(255),
    Skills MEDIUMTEXT,
    RawProfile MEDIUMBLOB,
	PRIMARY KEY (ID)
);
`

const PositionSchema = `create TABLE Position (
	ID int NOT NULL AUTO_INCREMENT,
	ProfileID int,
    Title varchar(255),
    CompanyName varchar(255),
	EndDate varchar(255),
	StartDate varchar(255),
	PRIMARY KEY (ID),
	FOREIGN KEY (ProfileID) REFERENCES Profile(ID)
);
`

const EducationSchema = `create TABLE Education (
	ID int NOT NULL AUTO_INCREMENT,
	ProfileID int,
	EndDate varchar(255),
	StartDate varchar(255),
	DegreeName varchar(255),
	FieldOfStudy varchar(255),
	SchoolName varchar(255),
	PRIMARY KEY (ID),
	FOREIGN KEY (ProfileID) REFERENCES Profile(ID)
);

`
