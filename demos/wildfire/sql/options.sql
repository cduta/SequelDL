INSERT INTO integer_options(id, name, value) VALUES 
(1, 'window width'     , 1024),
(2, 'window height'    ,  768),
(3, 'default font size',   30),
(4, 'fps'              ,   60),
(5, 'extinguish chance',   10), -- 0 <= x <= 100
(6, 'spread chance'    ,    4); -- 0 <= x <= 100

INSERT INTO text_options(id, name, value) VALUES 
(1, 'window title' , 'Wildfire Demo'),
(2, 'default font' , 'ressources/font/DejaVuSansMono.ttf');

INSERT INTO boolean_options(id, name, value) VALUES 
(1, 'show fps'     , false);