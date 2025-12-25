package dbutils

const train = `create table if not exists train (
	id integer primary key autoincrement,
	driver_name varchar(64) null,
	operating_status boolean)`

const station = `create table if not exists station (
	id integer primary key autoincrement,
	name varchar(64) null,
	opening_time time null,
	closing_time time null)`

const schedule = `create table if not exists schedule (
	id integer primary key autoincrement,
	train_id int,
	station_id int,
	arrival_time time,
	foreign key(train_id) references train(id),
	foreign key(station_id) references station(id))`