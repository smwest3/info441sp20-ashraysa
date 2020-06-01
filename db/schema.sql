create table if not exists Users (
    id bigint not null auto_increment primary key,
    email varchar(80) not null,
    passHash char(60) not null,
    username varchar(80) not null,
    firstName varchar(64) not null,
    lastName varchar(128) not null,
    bio varchar(500),
    points int not null,
    photoUrl varchar(255) not null
);

insert into Users(email, passHash, username, firstName, lastName, bio, points, photoUrl) values ("root@root.com", "NA", "rootuser", "root", "user", "this is the root user", 0, "NA");

CREATE UNIQUE INDEX email_index
ON Users (email);

CREATE UNIQUE INDEX username_index
ON Users (username);

create table if not exists Progress (
    progressID int primary key auto_increment not null,
    daysSober int not null,
    dateLogged datetime default CURRENT_TIMESTAMP not null,
    userID bigint not null,
    foreign key (userID) references Users(id)
);

create table if not exists Marketplace (
    badgeID int primary key auto_increment not null,
    cost int not null,
    badgeName varchar(30) not null,
    badgeDescription varchar(255) not null,
    imgURL varchar(255) not null
);

insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (0, "New Member", "Add this to your profile for free! Think of this as a gift from us to you for having the courage to log onto our website.", "images/award.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (100, "One Day Sober", "Congratulations for your first step! It may not seem like much, but we are now your fans! It's time to Eradicate The Vape.", "images/one_day.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (700, "One Week Sober", "You're keeping up well! The first step towards getting somewhere is to decide that you are not going to stay where you are", "images/one_week.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (3000, "One Month Sober", "Congratulations! This is a huge step. Persistence is key! You deserve this badge.", "images/one_month.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (18000, "Six Months Sober", "Congratulations for your first step! It may not seem like much, but we are now your fans! It's time to Eradicate The Vape.", "images/six_months.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (36500, "One Year Sober", "Congratulations for your first step! It may not seem like much, but we are now your fans! It's time to Eradicate The Vape.", "images/one_year.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (1500, "Vape-Tastic", "Congratulations for your first step! It may not seem like much, but we are now your fans! It's time to Eradicate The Vape.", "images/peach.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (10500, "Crown Juul", "Congratulations for your first step! It may not seem like much, but we are now your fans! It's time to Eradicate The Vape.", "images/apple.png");
insert into Marketplace (cost, badgeName, badgeDescription, imgURL) values (25000, "Planet of the Vapes", "Congratulations for your first step! It may not seem like much, but we are now your fans! It's time to Eradicate The Vape.", "images/heart.png");

create table if not exists Badges (
    badgeID int not null,
    userID bigint not null,
    foreign key (badgeID) references Marketplace(badgeID),
    foreign key (userID) references Users(id)
);

CREATE UNIQUE INDEX badge_index
ON Badges (badgeID, userID);

create table if not exists Threads (
    threadID int primary key auto_increment not null,
    threadName varchar(80) not null,
    threadDescription varchar(255),
    userWhoCreatedID bigint not null,
    timeCreated datetime not null,
    editedAt datetime not null,
    foreign key (userWhoCreatedID) references Users(id)
);

create table if not exists Posts (
    postID int primary key auto_increment not null,
    threadID int not null,
    content varchar (500) not null,
    userWhoCreatedID bigint not null,
    timeCreated datetime not null,
    editedAt datetime not null,
    foreign key (userWhoCreatedID) references Users(id),
    foreign key (threadID) references Threads(threadID)
);