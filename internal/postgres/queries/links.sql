-- name: InsertLink :exec
insert into links
       (url, title, tags, comments_count, views_count, rating, published_at, author, complexity, status)
values ($1,  $2,    $3,   $4,             $5,          $6,     $7,           $8,     $9,         $10);

-- name: UpdateLink :exec
update links set status = $1;

-- name: EnqueueUrl :exec
insert into links_queue (url)
values ($1);

-- name: DequeueUrl :one
with candidate as (
    select * from links_queue
    where picked_at is null
    order by added_at
    limit 1
    for update skip locked
)
update links_queue as q
set picked_at = now()
from candidate
where q.added_at = candidate.added_at and q.url = candidate.url
returning q.*;

-- name: DeleteQueuedUrl :exec
delete from links_queue where url = $1;
