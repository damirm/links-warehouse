package postgres_test

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/damirm/links-warehouse/internal/pkg/projectpath"
	"github.com/damirm/links-warehouse/internal/postgres"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDatabaseConfig = &postgres.Config{}

func TestMain(m *testing.M) {
	config, closer := runPostgresContainer(context.Background())
	testDatabaseConfig = config
	exitCode := m.Run()
	closer()
	os.Exit(exitCode)
}

func runPostgresContainer(ctx context.Context) (*postgres.Config, func()) {
	var (
		user = "myuser"
		pass = "secret"
		db   = "mydb"
	)

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
			"POSTGRES_DB":       db,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
	}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
	closer := func() {
		if err := pgC.Terminate(ctx); err != nil {
			log.Printf("failed to terminate postgres container: %#v", err)
		}
	}

	endpoint, err := pgC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	host, portString, err := net.SplitHostPort(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)

	log.Printf("current project path: %s\n", projectpath.Root)

	return &postgres.Config{
		Host:           host,
		Port:           uint(port),
		User:           user,
		Password:       pass,
		Database:       db,
		TimeZone:       "UTC",
		Schema:         "public",
		MigrationsPath: projectpath.Root + "/internal/postgres/migrations",
	}, closer
}
