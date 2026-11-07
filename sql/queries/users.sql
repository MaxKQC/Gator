-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name AS feed_name, feeds.url, users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT inserted_feed_follow.id,
       inserted_feed_follow.created_at,
       inserted_feed_follow.updated_at,
       inserted_feed_follow.user_id,
       inserted_feed_follow.feed_id,
       users.name AS user_name,
       feeds.name AS feed_name
FROM inserted_feed_follow
JOIN users ON inserted_feed_follow.user_id = users.id
JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.id,
       feed_follows.created_at,
       feed_follows.updated_at,
       feed_follows.user_id,
       feed_follows.feed_id,
       users.name AS user_name,
       feeds.name AS feed_name
FROM feed_follows
JOIN users ON feed_follows.user_id = users.id
JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO NOTHING;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC NULLS LAST, posts.created_at DESC
LIMIT $2;
