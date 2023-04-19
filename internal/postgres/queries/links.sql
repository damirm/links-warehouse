-- name: CreateLink :exec
insert into links
       (url, title, tags, comments_count, views_count, rating, published_at, author, complexity, status)
values ($1,  $2,    $3,   $4,             $5,          $6,     $7,           $8,     $9,         $10);
