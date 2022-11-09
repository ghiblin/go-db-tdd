package docker_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestDocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Suite")
}

const (
	dbName = "dbname"
	passwd = "secret"
	uname  = "user_name"
)

var Db *sql.DB
var cleanupDocker func()

var _ = BeforeSuite(func() {
	// setup postgres db with docker
	Db, cleanupDocker = setupDbWithDocker()
})

var _ = AfterSuite(func() {
	// clean up resource
	cleanupDocker()
})

var _ = BeforeEach(func() {
	// clear db tables before each test
	_, err := Db.Exec(`DROP SCHEMA public CASCADE;CREATE SCHEMA public;`)
	Ω(err).To(Succeed())
})

func setupDbWithDocker() (*sql.DB, func()) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres", // image
		Tag:        "14",       // version
		Env: []string{"POSTGRES_PASSWORD=" + passwd,
			"POSTGRES_USER=" + uname,
			"POSTGRES_DB=" + dbName,
			"listen_addresses = '*'"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", uname, passwd, hostAndPort, dbName)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	var db *sql.DB
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	fnCleanup := func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resources: %s", err)
		}
	}

	return db, fnCleanup
}
