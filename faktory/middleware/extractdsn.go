package middleware

import (
	"context"
	"os"
	"strconv"

	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/logger"
)

// ExtractDSN is a worker middleware get the DSNDB in job custom property
func ExtractDSN(perform worker.Handler) worker.Handler {
	// DB conn variables
	var (
		dbMaxOpenConns int = 5
		dbMaxIdleConns int = 1
		dbMaxLifeTime  int = 10
	)

	if len(os.Getenv("GOWORKER_DB_MAXOPENCONNS")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GOWORKER_DB_MAXOPENCONNS")); e != nil {
			dbMaxOpenConns = v
		}
	}

	if len(os.Getenv("GOWORKER_DB_MAXIDLECONNS")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GOWORKER_DB_MAXIDLECONNS")); e != nil {
			dbMaxIdleConns = v
		}
	}

	if len(os.Getenv("GOWORKER_DB_MAXLIFETIME")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GOWORKER_DB_MAXLIFETIME")); e != nil {
			dbMaxLifeTime = v
		}
	}

	return func(ctx worker.Context, job *faktory.Job) (err error) {
		if dsn, ok := job.Custom["dsn"]; ok {
			if _, exists := db.Connections.List[job.Queue]; !exists {
				pgDB := db.PostgresDB{
					ConnectionStr: dsn.(string),
					MaxOpenConns:  dbMaxOpenConns,
					MaxIdleConns:  dbMaxIdleConns,
					MaxLifetime:   dbMaxLifeTime,
				}

				var stack string
				if s, ok := job.Custom["stack"]; ok {
					stack = s.(string)
				} else {
					stack = "default"
				}

				if err := db.Connections.Connect(stack, &pgDB); err != nil {
					logger.Infof("DSN: %s\n", dsn)
					logger.Errorln(err)
				}
			}

			ctx = &worker.DefaultContext{
				Context: context.WithValue(ctx, "DSN", dsn),
				JID:     job.Jid,
				Type:    job.Type,
			}
		}

		return perform(ctx, job)
	}
}
