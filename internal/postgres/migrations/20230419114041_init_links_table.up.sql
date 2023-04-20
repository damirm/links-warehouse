create table "links" (
    "url" varchar primary key,
    "title" varchar not null,
    "tags" jsonb,
    "comments_count" integer not null,
    "views_count" integer not null,
    "rating" integer not null,
    "published_at" timestamp with time zone not null,
    "author" jsonb,
    "complexity" smallint not null,
    "status" smallint not null
);

create table "links_queue" (
    "added_at" timestamp default(now()),
    "url" varchar not null,
    "picked_at" timestamp,
    primary key ("added_at", "url")
);
create index links_queue_url_idx on links_queue (url);
create index links_queue_added_at_idx on links_queue (added_at, picked_at);
