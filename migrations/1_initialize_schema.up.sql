CREATE TABLE searchhistory
(
	id varchar(36) not null ,
	userid varchar(33) not null,
	keyword varchar(255) not null,
	time bigint not null,
	constraint story_pk
		primary key (id)
);
