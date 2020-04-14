/*
  https://www.sqlite.org/pragma.html#pragma_synchronous
  With synchronous OFF (0), SQLite continues without syncing as soon as it has handed data off to the operating system. If the 
  application running SQLite crashes, the data will be safe, but the database might become corrupted if the operating system crashes or 
  the computer loses power before that data has been written to the disk surface. On the other hand, commits can be orders of magnitude 
  faster with synchronous OFF. 
*/
PRAGMA synchronous=OFF; 

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
CREATE INDEX dots_position_idx ON dots(x,y);

CREATE TABLE lines (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id), 
  here_x    integer NOT NULL,
  here_y    integer NOT NULL,
  there_x   integer NOT NULL, 
  there_y   integer NOT NULL
);

CREATE UNIQUE INDEX lines_object_id_idx ON lines(object_id);
CREATE INDEX lines_here_idx ON lines(here_x,here_y);
CREATE INDEX lines_there_idx ON lines(there_x,there_y);

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

CREATE TABLE rectangles (
  id             integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id      integer NOT NULL REFERENCES objects(id),
  top_left_x     integer NOT NULL,
  top_left_y     integer NOT NULL,
  bottom_right_x integer NOT NULL,
  bottom_right_y integer NOT NULL
);

CREATE TABLE polygons (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id)
);

CREATE TABLE polygon_vertices (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  polygon_id integer NOT NULL REFERENCES polygons(id),
  i          integer NOT NULL CHECK(i >= 0), 
  x          integer NOT NULL,
  y          integer NOT NULL,
  UNIQUE(id, i)
);

CREATE TABLE colors (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,  
  object_id  integer NOT NULL REFERENCES objects(id),
  r          integer NOT NULL CHECK(r BETWEEN 0 AND 255),
  g          integer NOT NULL CHECK(g BETWEEN 0 AND 255), 
  b          integer NOT NULL CHECK(b BETWEEN 0 AND 255),
  a          integer NOT NULL CHECK(a BETWEEN 0 AND 255)
);

CREATE TABLE entities (
  id        integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  object_id integer NOT NULL REFERENCES objects(id), 
  name      text    NOT NULL UNIQUE,
  x         integer NOT NULL,
  y         integer NOT NULL,
  level     integer NOT NULL,
  visible   boolean NOT NULL
);

CREATE TABLE scenes (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name       text    NOT NULL UNIQUE
);

CREATE TABLE images (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
  name       text    NOT NULL UNIQUE,
  image_path text    NOT NULL UNIQUE 
);

CREATE TABLE entities_scenes (
  entity_id integer NOT NULL REFERENCES entities(id),      
  scene_id  integer NOT NULL REFERENCES scenes(id),
  PRIMARY KEY(entity_id, scene_id)
);

CREATE TABLE images_scenes (
  image_id integer NOT NULL REFERENCES images(id),      
  scene_id integer NOT NULL REFERENCES scenes(id),
  PRIMARY KEY(image_id, scene_id)
);

CREATE TABLE sprites (
  îd         integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  entity_id  integer NOT NULL REFERENCES entities(id),      
  image_id   integer NOT NULL REFERENCES images(id),
  name       text    NOT NULL UNIQUE,
  relative_x integer NOT NULL,
  relative_y integer NOT NULL,
  level      integer NOT NULL,
  width      integer NOT NULL,
  height     integer NOT NULL
);

CREATE TABLE hitboxes (
  id         integer NOT NULL PRIMARY KEY AUTOINCREMENT,  
  entity_id  integer NOT NULL REFERENCES entities(id), 
  name       text    NOT NULL UNIQUE,
  relative_x integer NOT NULL,
  relative_y integer NOT NULL,
  level      integer NOT NULL,
  width      integer NOT NULL,
  height     integer NOT NULL
); 

INSERT INTO objects(id) VALUES
(1);

INSERT INTO images(id, name, image_path) VALUES 
(1, 'button', 'ressources/sprites/button.png'),
(2, 'button-pressed', 'ressources/sprites/button-pressed.png');

INSERT INTO entities(id, object_id, name, x, y, level, visible) VALUES 
(1, 1, 'generic-button', 50, 50, 1, true);

INSERT INTO scenes(id, name) VALUES 
(1, 'menu');

INSERT INTO entities_scenes(entity_id, scene_id) VALUES 
(1,1);

INSERT INTO images_scenes(image_id, scene_id) VALUES 
(1,1),
(2,1);

INSERT INTO sprites(entity_id, image_id, name, relative_x, relative_y, level, width, height) VALUES
(1, 1, 'button-sprite', 0, 0, 1, 63, 20);

INSERT INTO hitboxes(entity_id, name, relative_x, relative_y, level, width, height) VALUES
(1, 'button-click', 0, 0, 1, 63, 20);