/*
  https://www.sqlite.org/pragma.html#pragma_synchronous
  With synchronous OFF (0), SQLite continues without syncing as soon as it has handed data off to the operating system. If the 
  application running SQLite crashes, the data will be safe, but the database might become corrupted if the operating system crashes or 
  the computer loses power before that data has been written to the disk surface. On the other hand, commits can be orders of magnitude 
  faster with synchronous OFF. 
*/
PRAGMA synchronous=OFF; 

CREATE TABLE integer_options (
  id    integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name  text    NOT NULL UNIQUE,
  value integer NOT NULL CHECK (value BETWEEN -2147483648 AND 2147483647) -- Golang's int32 constraint
);

CREATE TABLE text_options (
  id    integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name  text    NOT NULL UNIQUE,
  value text    NOT NULL
);

CREATE TABLE boolean_options (
  id    integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name  text    NOT NULL UNIQUE,
  value boolean NOT NULL
);

CREATE TABLE real_options (
  id    integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name  text    NOT NULL UNIQUE,
  value real    NOT NULL
);
