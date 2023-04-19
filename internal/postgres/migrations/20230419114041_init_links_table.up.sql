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
    "url" varchar primary key
);
