INSERT INTO objects(id) VALUES
(1),
(2);

INSERT INTO images(id, name, image_path) VALUES 
(1, 'button-idle'   , 'ressources/sprites/button-idle.png'),
(2, 'button-pressed', 'ressources/sprites/button-pressed.png');

INSERT INTO states(id, name) VALUES 
(1, 'button-idle'),
(2, 'button-pressed');

INSERT INTO entities(id, object_id, state_id, name, x, y, level, visible) VALUES 
(1, 1, 1, 'button_1', 100, 100, 1, true),
(2, 2, 1, 'button_2', 200, 100, 1, true);

INSERT INTO scenes(id, name, x, y, scene_x, scene_y, width, height, scroll_speed) VALUES 
(1, 'menu'  , 100, 100, 0, 0, 400, 400, 10);

INSERT INTO entities_scenes(entity_id, scene_id) VALUES 
(1, 1),
(2, 1);

INSERT INTO sprites(id, image_id, name, relative_x, relative_y, level, width, height) VALUES
(1, 1, 'button-idle'   , 0, 0, 1, 63, 20),
(2, 2, 'button-pressed', 0, 0, 1, 63, 20);

INSERT INTO states_sprites(state_id, sprite_id, animation_group) VALUES
(1, 1, 1),
(2, 2, 1);

INSERT INTO hitboxes(id, entity_id, name, relative_x, relative_y, level, width, height) VALUES
(1, 1, 'button-1-click', 0, 0, 1, 63, 20),
(2, 2, 'button-2-click', 0, 0, 1, 63, 20);