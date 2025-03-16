// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: follows.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT 
    iff.id, iff.created_at, iff.updated_at, iff.user_id, iff.feed_id, 
    u.name AS user_name, 
    f.name AS feed_name
FROM inserted_feed_follow iff
LEFT JOIN users u ON u.id = iff.user_id
LEFT JOIN feeds f ON f.id = iff.feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	UserName  sql.NullString
	FeedName  sql.NullString
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) ([]CreateFeedFollowRow, error) {
	rows, err := q.db.QueryContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CreateFeedFollowRow
	for rows.Next() {
		var i CreateFeedFollowRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.UserName,
			&i.FeedName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteAllFeedFollows = `-- name: DeleteAllFeedFollows :exec
DELETE FROM feed_follows
`

func (q *Queries) DeleteAllFeedFollows(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllFeedFollows)
	return err
}

const getFeedFollows = `-- name: GetFeedFollows :many
SELECT id, created_at, updated_at, user_id, feed_id FROM feed_follows
`

func (q *Queries) GetFeedFollows(ctx context.Context) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT u.name AS username, f.name AS feedname FROM feed_follows ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id
WHERE $1::TEXT = u.name
`

type GetFeedFollowsForUserRow struct {
	Username sql.NullString
	Feedname string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, username string) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(&i.Username, &i.Feedname); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedFollowsForUserFeedPair = `-- name: GetFeedFollowsForUserFeedPair :one
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, feed_id, u.id, u.created_at, u.updated_at, u.name, f.id, f.created_at, f.updated_at, f.name, url, f.user_id FROM feed_follows ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id
WHERE $1::TEXT = u.name AND $2::TEXT = f.url
`

type GetFeedFollowsForUserFeedPairParams struct {
	Username string
	Feedurl  string
}

type GetFeedFollowsForUserFeedPairRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID
	FeedID      uuid.UUID
	ID_2        uuid.UUID
	CreatedAt_2 time.Time
	UpdatedAt_2 time.Time
	Name        sql.NullString
	ID_3        uuid.UUID
	CreatedAt_3 time.Time
	UpdatedAt_3 time.Time
	Name_2      string
	Url         string
	UserID_2    uuid.UUID
}

func (q *Queries) GetFeedFollowsForUserFeedPair(ctx context.Context, arg GetFeedFollowsForUserFeedPairParams) (GetFeedFollowsForUserFeedPairRow, error) {
	row := q.db.QueryRowContext(ctx, getFeedFollowsForUserFeedPair, arg.Username, arg.Feedurl)
	var i GetFeedFollowsForUserFeedPairRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.ID_2,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
		&i.Name,
		&i.ID_3,
		&i.CreatedAt_3,
		&i.UpdatedAt_3,
		&i.Name_2,
		&i.Url,
		&i.UserID_2,
	)
	return i, err
}

const unfollowFeedFollow = `-- name: UnfollowFeedFollow :exec
DELETE FROM feed_follows ff
USING users u, feeds f
WHERE ff.user_id = u.id
AND ff.feed_id = f.id
AND u.name = $1::TEXT
AND f.url = $2::TEXT
`

type UnfollowFeedFollowParams struct {
	Username string
	Feedurl  string
}

func (q *Queries) UnfollowFeedFollow(ctx context.Context, arg UnfollowFeedFollowParams) error {
	_, err := q.db.ExecContext(ctx, unfollowFeedFollow, arg.Username, arg.Feedurl)
	return err
}
