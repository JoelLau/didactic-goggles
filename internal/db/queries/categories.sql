-- name: ListCategories :many
SELECT name
FROM categories 
ORDER BY name ASC;
