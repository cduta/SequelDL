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
  x         integer NOT NULL,
  y         integer NOT NULL
);

CREATE TABLE rectangle (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id),
  x         integer NOT NULL,
  y         integer NOT NULL
);

CREATE TABLE triangle (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id),
  x1        integer NOT NULL,
  y1        integer NOT NULL,
  x2        integer NOT NULL,
  y2        integer NOT NULL
);

CREATE TABLE polygon (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES object(id)
);

CREATE TABLE polygon_vertex (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  polygon_id integer NOT NULL REFERENCES object(id),
  i          integer NOT NULL CHECK(i > 0), 
  x          integer NOT NULL,
  y          integer NOT NULL,
  UNIQUE(id, i)
);

CREATE UNIQUE INDEX pk_polygon_vertex ON polygon_vertex(id, i);

CREATE TABLE color (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,  
  object_id  integer NOT NULL REFERENCES object(id),
  r          integer NOT NULL CHECK(0 <= r AND r <= 255),
  g          integer NOT NULL CHECK(0 <= g AND g <= 255), 
  b          integer NOT NULL CHECK(0 <= b AND b <= 255),
  a          integer NOT NULL CHECK(0 <= a AND a <= 255)
);