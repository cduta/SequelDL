INSERT INTO 
colors(id, r  , g  , b  , a  ) VALUES 
      (1 , 200, 200, 200, 255),
      (2 ,  30,  30,  30, 255),
      (3 , 255,  45,  10, 255),
      (4 , 255, 189,  46, 255);

INSERT INTO 
color_ranges(id, color_from, color_to, redraw_delay) VALUES
            (1 , 2         , 1       ,            0),
            (2 , 2         , 2       ,            3),
            (3 , 3         , 4       ,            0);            

INSERT INTO 
particles(id, color_range_id, name     , relative_x, relative_y, level) VALUES
         (1 , 1             , 'Paper'  ,  0        ,  0        , 1    ),
         (2 , 2             , 'Burnt'  ,  0        ,  0        , 1    ),
         (3 , 3             , '↖⚠'     , -1        , -1        , 2    ),
         (4 , 3             , '↑⚠'     ,  0        , -1        , 2    ),
         (5 , 3             , '↗⚠'     ,  1        , -1        , 2    ),
         (6 , 3             , '←⚠'     , -1        ,  0        , 2    ),
         (7 , 3             , '→⚠'     ,  1        ,  0        , 2    ),
         (8 , 3             , '↙⚠'     , -1        ,  1        , 2    ),
         (9 , 3             , '↓⚠'     ,  0        ,  1        , 2    ),
         (10, 3             , '↘⚠'     ,  1        ,  1        , 2    );

INSERT INTO 
states(id, name     ) VALUES
      (1 , 'Paper'  ),
      (2 , 'Burning'),
      (3 , 'Burnt'  );

INSERT INTO
entities(name, x, y, level, visible) 
WITH RECURSIVE
x(pos) AS (
  SELECT 1
    UNION ALL
  SELECT x.pos + 1
  FROM   x
  WHERE  x.pos < 200
),
y(pos) AS (
  SELECT 1 
    UNION ALL
  SELECT y.pos + 1  
  FROM   y
  WHERE  y.pos < 200
)
SELECT 'paper('||x.pos||','||y.pos||')', x.pos, y.pos, 1, true
FROM   x, y;


INSERT INTO 
states_particles(state_id, particle_id) VALUES
                (1       , 1          ),
                (2       , 2          ),
                (3       , 3          );

INSERT INTO
entities_states(entity_id, state_id) 
SELECT e.id, 1
FROM   entities AS e;