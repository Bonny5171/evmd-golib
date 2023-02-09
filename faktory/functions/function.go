package functions

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"

	"github.com/CognyHub/evmd-golib/db"
	"github.com/CognyHub/evmd-golib/db/dao"
	"github.com/CognyHub/evmd-golib/execlog"
	"github.com/CognyHub/evmd-golib/faktory/push"
	"github.com/CognyHub/evmd-golib/logger"
)

// A map of registered matchers for searching.
var funcs = make(map[string]Function)

type Function interface {
	Handler(ctx worker.Context, args ...interface{}) error
}

type innerFunc = func(conn, connCfg *sqlx.DB, ctx worker.Context, payload Payload, execID int64) error
type innerFuncNoLog = func(conn, connCfg *sqlx.DB, ctx worker.Context, payload Payload) error

func Get() map[string]Function {
	return funcs
}

// Register is called to register a function for use by the program.
func Add(functionName string, function Function) {
	if _, exists := funcs[functionName]; exists {
		log.Fatalln(functionName, "Function already registered")
	}

	log.Println("Register", functionName, "function")
	funcs[functionName] = function
}

func Run(fnName string, fn innerFunc, ctx worker.Context, args ...interface{}) error {
	// Parse payload that come of Faktory
	payload, err := ParsePayload(args...)
	if err != nil {
		return errorHandler(err, "ParsePayload()")
	}

	logger.Tracef("[%s][%s] Executing job function '%s'...\n", payload.StackName, ctx.Jid(), fnName)

	quit := make(chan struct{})
	defer close(quit)

	go pingJob(quit)

	// Queue
	queue := os.Getenv("GOWORKER_QUEUE_NAME")
	if len(queue) == 0 {
		queue = "default"
	}

	// Get connection with Config DB
	var connCfg *sqlx.DB
	if _, exists := db.Connections.List["CONFIG"]; exists {
		logger.Tracef("[%s][%s] Get connection with Config DB", payload.StackName, ctx.Jid())
		connCfg, err = db.GetConnection("CONFIG")
		if err != nil {
			return errorHandler(err, "db.GetConnection('CONFIG')")
		}
	}

	// Get connection with Data DB
	logger.Tracef("[%s][%s] Get connection with Data DB", payload.StackName, ctx.Jid())
	connData, err := db.GetConnection(payload.StackName)
	if err != nil {
		return errorHandler(err, fmt.Sprintf("db.GetConnection('%s')", payload.StackName))
	}

	// Create log execution on itgr.execution table
	logger.Tracef("[%s][%s] Create log execution on itgr.execution table", payload.StackName, ctx.Jid())
	exec, err := execlog.NewExec(connData, ctx.Jid(), payload.JobID, payload.JobName, payload.TenantID, 0, dao.EnumTypeStatusExec)
	if err != nil {
		return errorHandler(err, "execlog.NewExec()")
	}
	logger.Debugf("[%s][%s] Execution write ok", payload.StackName, ctx.Jid())

	logger.Debugf("[%s][%s] Get Job Info", payload.StackName, ctx.Jid())
	jobInfo, err := dao.GetJobByFuncQueue(connCfg, payload.TenantID, payload.StackName, fnName, queue, payload.JobName)
	if err != nil {
		return errorHandler(err, "dao.GetJobByFuncQueue()")
	}

	// Verifying concurrency
	if payload.AllowsConcurrency == false {
		//Checking if this job is executing
		logger.Debugf("[%s][%s] CheckJobs Execution", payload.StackName, ctx.Jid())
		executing, err := dao.ExecSFCheckJobsExecution(connData, payload.TenantID, jobInfo.ID, "processing")
		if err != nil {
			return exec.LogError(errorHandler(err, "dao.ExecSFCheckJobsExecution()"))
		}

		if executing {
			if payload.AllowsSchedule {
				// Get DSN from context
				dsn := cast.ToString(ctx.Value("DSN"))

				// push this job as a scheduled job on faktory
				if err := push.RetryLater(ctx.JobType(), queue, payload.StackName, dsn, args, 5*time.Minute); err != nil {
					return exec.LogError(errorHandler(err, "retryLater()"))
				}

				exec.LogExecution(dao.EnumStatusExecScheduled)
				logger.Tracef("[%s][%s] Job scheduled", payload.StackName, ctx.Jid())
			} else {
				exec.LogExecution(dao.EnumStatusExecOverrided)
				logger.Tracef("[%s][%s] Job overrided", payload.StackName, ctx.Jid())
			}

			return nil
		}
	}

	// Start log execution on itgr.execution table
	logger.Tracef("[%s][%s] Start log execution on itgr.execution table", payload.StackName, ctx.Jid())
	exec.LogExecution(dao.EnumStatusExecProcessing)

	if e := fn(connData, connCfg, ctx, payload, exec.ID); e != nil {
		return exec.LogError(errorHandler(e, "fn(conn, connCfg, payload, exec.ID)"))
	}

	// Log success on itgr.execution table
	logger.Tracef("[%s][%s] Logging success on itgr.execution table", payload.StackName, ctx.Jid())
	exec.LogExecution(dao.EnumStatusExecSuccess)

	logger.Tracef("[%s][%s] '%s' job function done!\n", payload.StackName, ctx.Jid(), fnName)

	return nil
}

// RunNoLog func
func RunNoLog(fnName string, fn innerFuncNoLog, ctx worker.Context, args ...interface{}) error {
	// Parse payload that come of Faktory
	payload, err := ParsePayload(args...)
	if err != nil {
		return errorHandler(err, "ParsePayload()")
	}

	logger.Tracef("[%s][%s] Executing job function '%s'...", payload.StackName, ctx.Jid(), fnName)

	quit := make(chan struct{})
	defer close(quit)

	go pingJob(quit)

	// Queue
	queue := os.Getenv("GOWORKER_QUEUE_NAME")
	if len(queue) == 0 {
		queue = "default"
	}

	// Get connection with Config DB
	var connCfg *sqlx.DB
	if _, exists := db.Connections.List["CONFIG"]; exists {
		logger.Tracef("[%s][%s] Get connection with Config DB", payload.StackName, ctx.Jid())
		connCfg, err = db.GetConnection("CONFIG")
		if err != nil {
			return errorHandler(err, "db.GetConnection('CONFIG')")
		}
	}

	// Get connection with Data DB
	logger.Tracef("[%s][%s] Get connection with Data DB", payload.StackName, ctx.Jid())
	connData, err := db.GetConnection(payload.StackName)
	if err != nil {
		return errorHandler(err, fmt.Sprintf("db.GetConnection('%s')", payload.StackName))
	}

	if e := fn(connData, connCfg, ctx, payload); e != nil {
		return errorHandler(e, "fn(conn, connCfg, payload, exec.ID)")
	}

	logger.Tracef("[%s][%s] '%s' job function done!", payload.StackName, ctx.Jid(), fnName)

	return nil
}

func pingJob(quit <-chan struct{}) {
	for {
		select {
		case <-quit:
			logger.Infoln("Quiting from pingJob loop")
			return
		default:
		}

		appEngineName := os.Getenv("GAE_SERVICE")
		cloudProject := os.Getenv("GOOGLE_CLOUD_PROJECT")

		if len(appEngineName) == 0 || len(cloudProject) == 0 {
			return
		}

		var sb strings.Builder
		sb.WriteString("http://")
		sb.WriteString(appEngineName)
		sb.WriteString("-dot-")
		sb.WriteString(cloudProject)
		sb.WriteString(".appspot.com/")

		logger.Infoln("ping: " + sb.String())

		response, err := http.Get(sb.String())
		if err != nil {
			logger.Errorln(fmt.Errorf("http.Get(): %w", err))
		}

		if response.StatusCode/100 != 2 {
			err := fmt.Errorf("job %s unavaliable", appEngineName)
			logger.Errorln(err)
		}

		// Sleep time
		pingSleepTime := cast.ToInt64(os.Getenv("PINGSLEEPTIME"))
		if pingSleepTime == 0 {
			pingSleepTime = 30
		}
		time.Sleep(time.Duration(pingSleepTime) * time.Second)
	}
}

func errorHandler(err error, stack string) error {
	if err != nil {
		err = fmt.Errorf("%s: %w", stack, err)
		logger.Errorln(err)
		return err
	}
	return nil
}
