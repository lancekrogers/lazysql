//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

const (
	postgresPort = "5432/tcp"
)

type TestContainer struct {
	container testcontainers.Container
	ctx       context.Context
	t         *testing.T

	host     string
	port     string
	user     string
	password string
	database string
}

func NewSharedContainer() (*TestContainer, error) {
	ctx := context.Background()
	user := "postgres"
	password := "postgres"
	database := "lazysql_test"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{postgresPort},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       database,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
		AutoRemove: true,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(postgresPort))
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	return &TestContainer{
		container: container,
		ctx:       ctx,
		t:         nil,
		host:      host,
		port:      mappedPort.Port(),
		user:      user,
		password:  password,
		database:  database,
	}, nil
}

func (tc *TestContainer) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", tc.user, tc.password, tc.host, tc.port, tc.database)
}

func (tc *TestContainer) Database() string {
	return tc.database
}

func (tc *TestContainer) Cleanup() {
	if tc.container == nil {
		return
	}
	_ = tc.container.Terminate(tc.ctx)
}

func (tc *TestContainer) Reset() error {
	ctx, cancel := context.WithTimeout(tc.ctx, 15*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", tc.DSN())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	_, err = db.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		return fmt.Errorf("reset schema: %w", err)
	}

	return nil
}
