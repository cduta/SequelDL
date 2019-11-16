CREATE TABLE object (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
  x  integer NOT NULL, 
  y  integer NOT NULL
);

CREATE TABLE dot (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id)
);

CREATE TABLE line (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id), 
  x         integer,
  y         integer
);

CREATE TABLE rectangle (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id),
  x         integer,
  y         integer
);

CREATE TABLE triangle (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id),
  x_a       integer,
  y_a       integer,
  x_b       integer,
  y_b       integer
);

CREATE TABLE polygon (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id)
);

CREATE TABLE polygon_vertex (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  polygon_id integer NOT NULL REFERENCES object(id),
  i          integer, 
  x          integer,
  y          integer
);

