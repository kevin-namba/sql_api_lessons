use lesson1

drop table if exists users;
create table users(
id varchar(20) ,
token varchar(20),
name varchar(20),
primary key(id)
);

insert into users(id,token,name)values("samplei1","sampletoken1","samplename1"),("samplei2","sampletoken2","samplename2");

drop table if exists characters;
create table characters(
       characterid varchar(20),
       name varchar(20),
       primary key(characterid));

insert into characters (characterid,name) values ('chara1','a'),('chara2','b'),('chara3','c');

drop table if exists gachatable;
create table gachatable(
       id int unsigned,
       characterid varchar(20),
       primary key(id));

insert into gachatable (id,characterid) values (1,'chara1'),(2,'chara1'),(3,'chara1'),(4,'chara2'),(5,'chara2'),(6,'chara3');

drop table if exists usercharacter;
create table usercharacter(
       usercharacterid varchar(20),
       characterid varchar(20),
       usertoken varchar(20),
       primary key(usercharacterid)
       );

insert into usercharacter (usercharacterid,characterid,usertoken) values ('1','chara1','sampletoken1'),(2,'chara1','sampletoken1'),(3,'chara1','sampletoken1'),(4,'chara2','sampletoken1'),(5,'chara2','sampletoken2'),(6,'chara3','sampletoken3');


