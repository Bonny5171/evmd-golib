package functions

import (
	"fmt"
	"log"
	"os"
	"time"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/dao"
	"bitbucket.org/everymind/evmd-golib/execlog"
	"bitbucket.org/everymind/evmd-golib/faktory/push"
	"bitbucket.org/everymind/evmd-golib/logger"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// A map of registered matchers for searching.
var funcs = make(map[string]Function)

type Function interface {
	Handler(ctx worker.Context, args ...interface{}) error
}

type innerFunc = func(conn, connCfg *sqlx.DB, ctx worker.Context, payload Payload, execID int64) error

func Get() map[string]Function {
	return funcs
}

// Register is called to register a function for use by the program.
func Register(functionName string, function Function) {
	if _, exists := funcs[functionName]; exists {
		log.Fatalln(functionName, "Function already registered")
	}

	log.Println("Register", functionName, "function")
	funcs[functionName] = function
}

func Run(fnName string, fn innerFunc, ctx worker.Context, args ...interface{}) error {
	logger.Tracef("Executing '%s' job function...\n", fnName)

	// Get stack name from context
	stack := cast.ToString(ctx.Value("Stack"))

	// Parse payload that come of Faktory
	payload, err := ParsePayload(args...)
	if err != nil {
		return errorHandler(err, "ParsePayload()")
	}

	// Get connection with Config DB
	var connCfg *sqlx.DB
	if _, exists := db.Connections.List["CONFIG"]; exists {
		logger.Traceln("Get connection with Config DB")
		connCfg, err = db.GetConnection("CONFIG")
		if err != nil {
			return errorHandler(err, "db.GetConnection('CONFIG')")
		}
	}

	// Get connection with Data DB
	logger.Traceln("Get connection with Data DB")
	conn, err := db.GetConnection(stack)
	if err != nil {
		return errorHandler(err, fmt.Sprintf("db.GetConnection('%s')", stack))
	}

	// Create log execution on itgr.execution table
	logger.Traceln("Create log execution on itgr.execution table")
	exec, err := execlog.NewExec(conn, ctx.Jid(), payload.ID, payload.TenantID, 0, dao.EnumTypeStatusExec)
	if err != nil {
		return errorHandler(err, "execlog.NewExec()")
	}

	// Verifying concurrency
	if !payload.AllowsConcurrency {
		//Checking if this job is executing
		notok, err := dao.ExecSFCheckJobsExection(conn, payload.TenantID, fnName, "processing")
		if err != nil {
			return exec.LogError(errorHandler(err, "dao.ExecSFCheckJobsExection()"))
		}

		if notok {
			queue := os.Getenv("GOWORKER_QUEUE_NAME")
			if len(queue) == 0 {
				queue = "default"
			}

			if payload.AllowsSchedule {
				// Get DSN from context
				dsn := cast.ToString(ctx.Value("DSN"))

				// push this job as a scheduled job on faktory
				if err := push.RetryLater(ctx.JobType(), queue, stack, dsn, args, 5*time.Minute); err != nil {
					return exec.LogError(errorHandler(err, "retryLater()"))
				}

				exec.LogExecution(dao.EnumStatusExecScheduled)
				logger.Tracef("[%s] Job scheduled", ctx.Jid())
			} else {
				exec.LogExecution(dao.EnumStatusExecOverrided)
				logger.Tracef("[%s] Job overrided", ctx.Jid())
			}

			return nil
		}
	}

	// Start log execution on itgr.execution table
	logger.Traceln("Start log execution on itgr.execution table")
	exec.LogExecution(dao.EnumStatusExecProcessing)

	if e := fn(conn, connCfg, ctx, payload, exec.ID); e != nil {
		return exec.LogError(errorHandler(e, "fn(conn, connCfg, payload, exec.ID)"))
	}

	// Log success on itgr.execution table
	logger.Traceln("Logging success on itgr.execution table")
	exec.LogExecution(dao.EnumStatusExecSuccess)

	logger.Tracef("'%s' job function done!\n", fnName)

	return nil
}

func errorHandler(err error, stack string) error {
	if err != nil {
		err = errors.Wrap(err, stack)
		logger.Errorln(err)
		return err
	}
	return nil
}
