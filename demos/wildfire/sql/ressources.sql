INSERT INTO 
colors(id, r  , g  , b , a  ) VALUES 
      (1 , 255,  45, 10, 255),
      (2 , 255, 189, 46, 255);

INSERT INTO 
color_ranges(id, color_from, color_to) VALUES
            (1 , 1         , 2       );

INSERT INTO 
particles(id, color_range_id, name     , relative_x, relative_y, level, width, height) VALUES
         (1 , 1             , 'Fire'   , 0         , 0         , 1    , 1    , 1     );

INSERT INTO 
states(id, name     , next_state, ticks_left, ticks) VALUES
      (1 , 'Paper'  , NULL      , NULL      , NULL ),
      (2 , 'Burning', 2         , 3         , 3    ),
      (3 , 'Burnt'  , NULL      , NULL      , NULL );

INSERT INTO
entities(id, name   , x  , y  , level, visible) VALUES
        (1 , 'torch', 200, 200, 1    , true   );

INSERT INTO 
states_particles(state_id, particle_id) VALUES
                (2       , 1          );
                
INSERT INTO 
entities_states(entity_id, state_id) VALUES
               (1        , 2       );