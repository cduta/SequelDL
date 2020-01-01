-- Apply Bresenham's line drawing algorithm to a line connecting l.here and l.there
WITH RECURSIVE
preparation(id, object_id, x0, y0, there_x, there_y, dx, dy, sx, sy) AS (
  SELECT l.id, l.object_id, l.here_x, l.here_y, l.there_x, l.there_y, 
    ABS(l.here_x - l.there_x), 
    -ABS(l.here_y - l.there_y),
    CASE WHEN l.here_x < l.there_x THEN 1 ELSE -1 END, 
    CASE WHEN l.here_y < l.there_y THEN 1 ELSE -1 END 
  FROM   lines AS l 
),
bresenham(id, object_id, x, y, error, there_x, there_y, dx, dy, sx, sy) AS (
  SELECT p.id, p.object_id, p.x0, p.y0, p.dx+p.dy, p.there_x, p.there_y, p.dx, p.dy, p.sx, p.sy
  FROM   preparation AS p 
    UNION ALL 
  SELECT b.id, b.object_id,
  CASE 
    WHEN 2*b.error >= b.dy THEN b.x + b.sx
    ELSE b.x
  END,
  CASE 
    WHEN 2*b.error <= b.dx THEN b.y + b.sy
    ELSE b.y
  END, 
  CASE 
    WHEN 2*b.error BETWEEN b.dy AND b.dx THEN b.error + b.dx + b.dy
    WHEN 2*b.error >= b.dy               THEN b.error + b.dy
    WHEN 2*b.error <= b.dx               THEN b.error + b.dx 
    ELSE b.error                     
  END, b.there_x, b.there_y, b.dx, b.dy, b.sx, b.sy
  FROM   bresenham AS b
  WHERE  NOT (b.x = b.there_x AND b.y = b.there_y)
)
SELECT b.id, b.object_id, b.x, b.y, b.there_x, b.there_y, c.r, c.g, c.b, c.a
FROM   bresenham AS b, colors AS c 
WHERE  b.object_id = c.object_id
ORDER BY b.object_id, b.x = b.there_x AND b.y = b.there_y;