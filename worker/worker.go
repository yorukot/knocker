package worker

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/knocker/repository"
	"github.com/yorukot/knocker/utils/config"
	"github.com/yorukot/knocker/worker/handler"
	"github.com/yorukot/knocker/worker/tasks"
	"go.uber.org/zap"
)

func Run(db *pgxpool.Pool) {
	zap.L().Info("Starting worker")
	cfg := config.Env()

	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: cfg.RedisPassword,
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10000,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	notifier := asynq.NewClient(redisOpt)
	defer notifier.Close()

	repo := repository.New(db)
	h := handler.NewHandler(repo, notifier)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.GetMonitorPingType(config.Env().AppRegion), h.HandleStartServiceTask)
	mux.HandleFunc(tasks.TypeNotificationDispatch, h.HandleNotificationDispatch)

	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
