package leaderboardstorer

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboardrecalculator"
)

const (
	leaderboardKey = "leaderboard"
)

var (
	_ leaderboardrecalculator.LeaderboardStorer = (*RedisLeaderboardStoarer)(nil)
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
			Member: entry,
		})
	}

	return members
}
