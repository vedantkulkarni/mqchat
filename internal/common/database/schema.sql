drop table if exists sessions;
drop table if exists users;
drop table if exists messages;


create table users(
    user_id varchar not null primary key, 
    first_name varchar,
    last_name varchar,
    email varchar, 
    phone varchar 
);


create table sessions(
    user_id varchar not null primary key,
    session_id varchar not null unique,
    is_active boolean,
    last_active timestamp,

    foreign key (user_id) references users (user_id)

);


create table messages(
    message_id varchar not null primary key,
    from_user varchar not null,
    to_user varchar not null,
    created_at timestamp
);



