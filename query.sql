create table users(
id int primary key,
email varchar(50) NOT NULL,
first_name varchar(25),
last_name varchar(25),
password varchar(250) NOT NULL,
active int,
created_at time DEFAULT CURRENT_TIMESTAMP,
updated_at time DEFAULT CURRENT_TIMESTAMP)


alter table users
rename column update_at to updated_at


insert into users
(id, email, first_name, password, active)
values
(1, 'admin@example.com', 'admin', 'verysecret', 0);

select * from users;
select id, email, first_name, last_name, password, active, created_at, updated_at from users where email = 'admin@example.com';
