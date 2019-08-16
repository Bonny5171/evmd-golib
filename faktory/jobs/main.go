package jobs

import (
	"time"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/faktory/middleware"
	"bitbucket.org/everymind/evmd-golib/logger"
	worker "github.com/contribsys/faktory_worker_go"
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
	Queues      []string
	funcs       map[string]worker.Perform
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
		Queues: []string{"default"},
		funcs:  make(map[string]worker.Perform),
	}
}

func (j *Job) Register(name string, fn worker.Perform) {
	j.funcs[name] = fn
}

func (j *Job) Run() {
	// Setting config DB connection
	configDB := db.PostgresDB{
		ConnectionStr: j.DB.ConfigDSN,
		MaxOpenConns:  j.DB.MaxOpenConns,
		MaxIdleConns:  j.DB.MaxIdleConns,
		MaxLifetime:   j.DB.MaxLifeTime,
	}

	// Starting config DB connection
	if _, ok := db.Connections.List["CONFIG"]; ok {
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

	// Do anything when this job is quite
	mgr.On(worker.Quiet, func() {
		time.Sleep(30 * time.Second)
		logger.Warningln("Bye...")
		mgr.Terminate()
	})

	// Do anything when this job is shutdown
	mgr.On(worker.Shutdown, func() {
		logger.Warningln("ALERT, THIS JOB IS SHUTING DOWN!")
	})

	// register job types and the function to execute them
	for n, f := range j.funcs {
		mgr.Register(n, f)
		logger.Infof("Job '%s' registered on Faktory.", n)
	}

	// use up to N goroutines to execute jobs
	mgr.Concurrency = j.Concurrency

	// pull jobs from these queues, in this order of precedence
	if len(j.Queues) == 0 {
		mgr.ProcessStrictPriorityQueues("default")
	} else {
		mgr.ProcessStrictPriorityQueues(j.Queues...)
	}

	// Start processing jobs, this method does not return
	logger.Infoln("Waiting for processing jobs...")
	mgr.Run()
}
