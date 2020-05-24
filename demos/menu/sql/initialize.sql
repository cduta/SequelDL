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

CREATE TABLE objects (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT
);

CREATE TABLE dots (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
  object_id integer NOT NULL REFERENCES objects(id),
  x         integer NOT NULL, 
  y         integer NOT NULL
);

CREATE UNIQUE INDEX dots_object_id_idx ON dots(object_id);
CREATE        INDEX dots_position_idx  ON dots(x,y);

CREATE TABLE lines (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id), 
  here_x    integer NOT NULL,
  here_y    integer NOT NULL,
  there_x   integer NOT NULL, 
  there_y   integer NOT NULL
);

CREATE UNIQUE INDEX lines_object_id_idx ON lines(object_id);
CREATE        INDEX lines_here_idx      ON lines(here_x,here_y);
CREATE        INDEX lines_there_idx     ON lines(there_x,there_y);

CREATE TABLE colors (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,  
  object_id  integer NOT NULL REFERENCES objects(id),
  r          integer NOT NULL CHECK(r BETWEEN 0 AND 255),
  g          integer NOT NULL CHECK(g BETWEEN 0 AND 255), 
  b          integer NOT NULL CHECK(b BETWEEN 0 AND 255),
  a          integer NOT NULL CHECK(a BETWEEN 0 AND 255)
);

CREATE TABLE states (
  id   integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name text    NOT NULL UNIQUE
);

CREATE TABLE entities (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id), 
  state_id  integer NOT NULL REFERENCES states(id),
  name      text    NOT NULL UNIQUE,
  x         integer NOT NULL CHECK (x     BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  y         integer NOT NULL CHECK (y     BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  level     integer NOT NULL CHECK (level BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  visible   boolean NOT NULL
);

CREATE TABLE scenes (
  id           integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name         text    NOT NULL UNIQUE,
  x            integer NOT NULL CHECK (x            BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  y            integer NOT NULL CHECK (y            BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  scroll_x     integer NOT NULL CHECK (scroll_x     BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  scroll_y     integer NOT NULL CHECK (scroll_y     BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  width        integer NOT NULL CHECK (width        BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  height       integer NOT NULL CHECK (height       BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  scroll_speed integer NOT NULL CHECK (scroll_speed BETWEEN           0 AND 2147483647)  -- Golang's unsigned int32 constraint
);

CREATE TABLE sprites (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,  
  image_id   integer NOT NULL REFERENCES images(id),
  name       text    NOT NULL UNIQUE,
  relative_x integer NOT NULL CHECK (relative_x BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  relative_y integer NOT NULL CHECK (relative_y BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  level      integer NOT NULL CHECK (level      BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  width      integer NOT NULL CHECK (width      BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  height     integer NOT NULL CHECK (height     BETWEEN -2147483648 AND 2147483647)  -- Golang's int32 constraint
);

CREATE TABLE images (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
  name       text    NOT NULL UNIQUE,
  image_path text    NOT NULL UNIQUE 
);

CREATE TABLE states_sprites (
  state_id        integer NOT NULL REFERENCES states(id),
  sprite_id       integer NOT NULL REFERENCES sprites(id),
  animation_group integer NOT NULL
);

CREATE TABLE entities_scenes (
  entity_id integer NOT NULL REFERENCES entities(id),      
  scene_id  integer NOT NULL REFERENCES scenes(id),
  PRIMARY KEY(entity_id, scene_id)
);

CREATE TABLE hitboxes (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,  
  entity_id  integer NOT NULL REFERENCES entities(id), 
  name       text    NOT NULL UNIQUE,
  relative_x integer NOT NULL CHECK (relative_x BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  relative_y integer NOT NULL CHECK (relative_y BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  level      integer NOT NULL CHECK (level      BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  width      integer NOT NULL CHECK (width      BETWEEN -2147483648 AND 2147483647), -- Golang's int32 constraint
  height     integer NOT NULL CHECK (height     BETWEEN -2147483648 AND 2147483647)  -- Golang's int32 constraint
); 
