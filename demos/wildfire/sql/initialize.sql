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

CREATE TABLE states (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name       text    NOT NULL UNIQUE,
  next_state integer NOT NULL REFERENCES states(id),
  ticks_left integer          CHECK (ticks BETWEEN 0 AND 2147483647), -- Golang's unsigned int32 constraint
  ticks      integer          CHECK (ticks BETWEEN 0 AND 2147483647)  -- Golang's unsigned int32 constraint
);

CREATE TABLE entities (
  id            integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  state_id      integer NOT NULL REFERENCES states(id),
  name          text    NOT NULL UNIQUE,
  x             integer NOT NULL CHECK (x     BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  y             integer NOT NULL CHECK (y     BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  level         integer NOT NULL CHECK (level BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  visible       boolean NOT NULL
);

CREATE TABLE colors (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  r  integer NOT NULL,
  g  integer NOT NULL,
  b  integer NOT NULL,
  a  integer NOT NULL
);

CREATE TABLE color_ranges (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  color_from integer NOT NULL REFERENCES colors(id),
  color_to   integer NOT NULL REFERENCES colors(id)
);

CREATE TABLE particles (
  id             integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  color_range_id integer NOT NULL REFERENCES color_ranges(id),
  name           text    NOT NULL,
  relative_x     integer NOT NULL CHECK (relative_x BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  relative_y     integer NOT NULL CHECK (relative_y BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  level          integer NOT NULL CHECK (level      BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  width          integer NOT NULL CHECK (width      BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  height         integer NOT NULL CHECK (height     BETWEEN -2147483648 AND 2147483647)  -- Golang's int32 constraint
);

CREATE TABLE states_particles (
  state_id    integer NOT NULL REFERENCES states(id),
  particle_id integer NOT NULL REFERENCES particles(id)
);