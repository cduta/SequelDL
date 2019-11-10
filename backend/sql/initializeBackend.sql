CREATE TABLE objecttype (
  id   integer PRIMARY KEY NOT NULL, 
  name text                NOT NULL
);

CREATE TABLE objects (
  id integer PRIMARY KEY AUTOINCREMENT NOT NULL, 
  objecttype_id integer                NOT NULL REFERENCES objecttype(id), 
  x integer                            NOT NULL, 
  y integer                            NOT NULL
);

INSERT INTO objecttype VALUES
  (1,"line"),
  (2,"square");
