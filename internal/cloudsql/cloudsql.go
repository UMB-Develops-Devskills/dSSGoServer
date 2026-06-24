package cloudsql

import (
	"fmt"
	"context"
	"os"
	"net"
	"log"
	"github.com/joho/godotenv"
	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This function creates the connection pool to CLOUDSQL
func ConnectDB(ctx context.Context) (*pgxpool.Pool, func() error, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	instanceConnection := os.Getenv("INSTANCE_CONNECTION_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
  dbName := os.Getenv("DB_NAME")

	// dns network connnectio
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)

	dialer, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithLazyRefresh())
	if err != nil {
		return nil, nil, fmt.Errorf("could not create Cloud SQl dialer %w", err)
	}
	config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
      dialer.Close()
      return nil, nil, fmt.Errorf("failed to parse config: %w", err)
  }
	config.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
    return dialer.Dial(ctx, instanceConnection)
  }
	// set pool limit
	config.MaxConns = 10
	config.MinConns = 2
 
	// creating the connection pool
	// do not use context.background() if u care for timeout
	pool, err := pgxpool.NewWithConfig(ctx, config) 
	if err != nil {
		dialer.Close()
		return nil, nil, fmt.Errorf("failed to create pool %w", err)
	}
	// verify the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		dialer.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Connected to Cloud SQL")
	return pool, dialer.Close, nil
}