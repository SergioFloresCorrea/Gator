-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
	INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
	VALUES (
    	$1,
    	$2,
    	$3,
    	$4
	)
	RETURNING *
)
SELECT
	inserted_feed_follow.*,
	users.name AS user_name,
	feeds.name AS feed_name
FROM inserted_feed_follow
INNER JOIN users ON users.id = user_id
INNER JOIN feeds ON feeds.id = feed_id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_id
WHERE feed_follows.user_id = $1;

-- name: UnFollowFeed :exec
DELETE FROM feed_follows
WHERE feed_id = (
	SELECT id FROM feeds
	WHERE url = $1
	LIMIT 1
)
AND feed_follows.user_id = $2;
