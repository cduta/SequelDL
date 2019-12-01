PRAGMA read_uncommitted = true;
PRAGMA synchronous=OFF;

CREATE TABLE objects (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT
);

CREATE TABLE dots (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
  x         integer NOT NULL, 
  y         integer NOT NULL,
  object_id integer NOT NULL REFERENCES objects(id)
);

CREATE TABLE lines (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id), 
  here_x    integer NOT NULL,
  here_y    integer NOT NULL,
  there_x   integer NOT NULL, 
  there_y   integer NOT NULL
);

CREATE TABLE rectangles (
  id             integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id      integer NOT NULL REFERENCES objects(id),
  top_left_x     integer NOT NULL,
  top_left_y     integer NOT NULL,
  bottom_right_x integer NOT NULL,
  bottom_right_y integer NOT NULL
);

CREATE TABLE triangles (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id),
  x1        integer NOT NULL,
  y1        integer NOT NULL,
  x2        integer NOT NULL,
  y2        integer NOT NULL,
  x3        integer NOT NULL,
  y3        integer NOT NULL
);

CREATE TABLE polygons (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id)
);

CREATE TABLE polygon_vertices (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  polygon_id integer NOT NULL REFERENCES objects(id),
  i          integer NOT NULL CHECK(i >= 0), 
  x          integer NOT NULL,
  y          integer NOT NULL,
  UNIQUE(id, i)
);

CREATE UNIQUE INDEX pk_polygon_vertices ON polygon_vertices(id, i);

CREATE TABLE colors (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,  
  object_id  integer NOT NULL REFERENCES objects(id),
  r          integer NOT NULL CHECK(0 <= r AND r <= 255),
  g          integer NOT NULL CHECK(0 <= g AND g <= 255), 
  b          integer NOT NULL CHECK(0 <= b AND b <= 255),
  a          integer NOT NULL CHECK(0 <= a AND a <= 255)
);