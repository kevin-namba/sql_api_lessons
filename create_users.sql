use lesson1
drop table if exists users;
create table users(
id text ,
token text,
name text,
primary key(id(8))
);

insert into users(id,token,name)values("samplei1","sampletoken1","samplename1"),("samplei2","sampletoken2","samplename2");
