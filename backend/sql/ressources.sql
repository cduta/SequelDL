INSERT INTO objects(id) VALUES
(1);

INSERT INTO images(id, name, image_path) VALUES 
(1, 'button', 'ressources/sprites/button.png'),
(2, 'button-pressed', 'ressources/sprites/button-pressed.png');

INSERT INTO entities(id, object_id, name, x, y, level, visible) VALUES 
(1, 1, 'generic-button', 200, 200, 1, true);

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