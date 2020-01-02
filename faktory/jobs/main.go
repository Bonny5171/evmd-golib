package jobs

import (
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"bitbucket.org/everymind/evmd-golib/db"
	fn "bitbucket.org/everymind/evmd-golib/faktory/functions"
	"bitbucket.org/everymind/evmd-golib/faktory/middleware"
	"bitbucket.org/everymind/evmd-golib/logger"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/gorilla/mux"
)

type DBVars struct {
	ConfigDSN    string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifeTime  int
}

type Job struct {
	Concurrency int
	DB          DBVars
	Labels      []string
	ProcessWID  string
	Queues      []string
}

// NewJob returns a new job with default values.
func NewJob() *Job {
	return &Job{
		Concurrency: 5,
		DB: DBVars{
			ConfigDSN:    "",
			MaxOpenConns: 5,
			MaxIdleConns: 1,
			MaxLifeTime:  10,
		},
		Labels: []string{"golang"},
		Queues: []string{"default"},
	}
}

func (j *Job) Run() {
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	wg := sync.WaitGroup{}

	go func() {
		sig := <-gracefulStop
		wg.Wait()
		logger.Warningf("Signal %s sended, graceful shutdown." + sig.String())
	}()

	fn.SetWG(&wg)

	// Starting web server
	startWebServer()

	// Setting config DB connection
	configDB := db.PostgresDB{
		ConnectionStr: j.DB.ConfigDSN,
		MaxOpenConns:  j.DB.MaxOpenConns,
		MaxIdleConns:  j.DB.MaxIdleConns,
		MaxLifetime:   j.DB.MaxLifeTime,
	}

	// Starting config DB connection
	if len(j.DB.ConfigDSN) > 0 {
		if err := db.Connections.Connect("CONFIG", &configDB); err != nil {
			logger.Infof("DSN: %s\n", j.DB.ConfigDSN)
			logger.Errorln(err)
		}
	}

	// New worker manager
	mgr := worker.NewManager()
	logger.Infof("Worker manager created")

	// Middleware to set Stack name on context
	mgr.Use(middleware.SetStackNameOnCtx)
	logger.Traceln("Middleware 'SetStackNameOnCtx' configured")

	// Middleware to extract DSNDB in job custom property and store on context
	mgr.Use(middleware.ExtractDSN)
	logger.Traceln("Middleware 'ExtractDSN' configured")

	// Do anything when this job is starting up
	mgr.On(worker.Startup, func() {
		logger.Infoln("Starting JOB, waiting for processing jobs...")
	})

	// Do anything when this job is quite
	mgr.On(worker.Quiet, func() {
		mgr.Terminate()
	})

	// Do anything when this job is shutdown
	mgr.On(worker.Shutdown, func() {
		logger.Warningln("This JOB is shutting down!")
	})

	// register job types and the function to execute them
	for n, f := range fn.Get() {
		mgr.Register(n, f.Handler)
		logger.Infof("Job '%s' registered on Faktory.", n)
	}

	// use up to N goroutines to execute jobs
	mgr.Concurrency = j.Concurrency

	if len(j.ProcessWID) > 0 {
		j.Labels = append(j.Labels, j.ProcessWID)
		rand.Seed(time.Now().UnixNano())
		mgr.ProcessWID = j.ProcessWID + "-" + strconv.FormatInt(rand.Int63(), 32)
	}

	// Labels to be displayed in the UI
	for _, q := range j.Queues {
		j.Labels = append(j.Labels, "queue:"+q)
	}
	mgr.Labels = j.Labels

	// pull jobs from these queues, in this order of precedence
	if len(j.Queues) == 0 {
		mgr.ProcessStrictPriorityQueues("default")
	} else {
		mgr.ProcessStrictPriorityQueues(j.Queues...)
	}

	// Start processing jobs, this method does not return
	mgr.Run()
}

func startWebServer() {
	go func() {
		router := mux.NewRouter().StrictSlash(true)

		router.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
			logger.Infoln("health check done")
			w.WriteHeader(200)
		}).Methods("GET")

		router.HandleFunc("/_ah/warmup", func(w http.ResponseWriter, r *http.Request) {
			logger.Infoln("warmup done")
			w.WriteHeader(200)
		}).Methods("GET")

		router.HandleFunc("/_ah/start", func(w http.ResponseWriter, r *http.Request) {
			logger.Infoln("start done")
			w.WriteHeader(200)
		}).Methods("GET")

		router.HandleFunc("/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}).Methods("GET")

		logger.Traceln("Starting HTTP server...")
		if err := http.ListenAndServe(":8080", router); err != nil {
			logger.Errorln(err)
		}
	}()
}
