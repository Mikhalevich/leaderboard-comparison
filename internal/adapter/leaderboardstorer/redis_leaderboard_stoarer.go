package leaderboardstorer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboardrecalculator"
)

const (
	leaderboardKey = "leaderboard"
	halfRankRange  = 2
)

var (
	_ leaderboardrecalculator.LeaderboardStorer = (*RedisLeaderboardStoarer)(nil)
	_ leaderboard.LeaderboardRepo               = (*RedisLeaderboardStoarer)(nil)
)

type RedisLeaderboardStoarer struct {
	client redis.UniversalClient
}

func New(client redis.UniversalClient) *RedisLeaderboardStoarer {
	return &RedisLeaderboardStoarer{
		client: client,
	}
}

func (r *RedisLeaderboardStoarer) StoreLeaderbord(
	ctx context.Context,
	top []leaderboardrecalculator.LeaderboardEntry,
) error {
	if _, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if err := pipe.Del(ctx, leaderboardKey).Err(); err != nil {
			return fmt.Errorf("delete leaderboard key: %w", err)
		}

		if len(top) == 0 {
			return nil
		}

		if err := pipe.ZAdd(ctx, leaderboardKey, makeZSetMembers(top)...).Err(); err != nil {
			return fmt.Errorf("zadd members: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func makeZSetMembers(entries []leaderboardrecalculator.LeaderboardEntry) []redis.Z {
	members := make([]redis.Z, 0, len(entries))

	for _, entry := range entries {
		members = append(members, redis.Z{
			Score:  float64(entry.Position),
			Member: entry.UserID,
		})
	}

	return members
}

func (r *RedisLeaderboardStoarer) LeaderboardTop(
	ctx context.Context,
	limit int,
) ([]leaderboard.LeaderbordEntry, error) {
	members, err := r.client.ZRangeWithScores(ctx, leaderboardKey, 0, int64(limit)).Result()
	if err != nil {
		return nil, fmt.Errorf("zrange with scores: %w", err)
	}

	entries, err := convertToLeaderboardEntry(members)
	if err != nil {
		return nil, fmt.Errorf("convert to leaderboard entries: %w", err)
	}

	return entries, nil
}

func convertToLeaderboardEntry(members []redis.Z) ([]leaderboard.LeaderbordEntry, error) {
	entries := make([]leaderboard.LeaderbordEntry, 0, len(members))

	for _, member := range members {
		rawUserID, ok := member.Member.(string)
		if !ok {
			return nil, errors.New("invalid leaderboard entry")
		}

		userID, err := strconv.ParseInt(rawUserID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("convert to int: %w", err)
		}

		entries = append(entries, leaderboard.LeaderbordEntry{
			UserID:   userID,
			Position: int(member.Score),
		})
	}

	return entries, nil
}

func (r *RedisLeaderboardStoarer) LeaderboardByUserID(
	ctx context.Context,
	userID int64,
	limit int,
) ([]leaderboard.LeaderbordEntry, error) {
	rank, err := r.client.ZRank(ctx, leaderboardKey, strconv.FormatInt(userID, 10)).Result()
	if err != nil {
		return nil, fmt.Errorf("zrank: %w", err)
	}

	halfRange := limit / halfRankRange

	members, err := r.client.ZRangeWithScores(
		ctx,
		leaderboardKey,
		rank-int64(halfRange),
		rank+int64(halfRange),
	).Result()
	if err != nil {
		return nil, fmt.Errorf("zrange with scores: %w", err)
	}

	entries, err := convertToLeaderboardEntry(members)
	if err != nil {
		return nil, fmt.Errorf("convert to leaderboard entries: %w", err)
	}

	return entries, nil
}
