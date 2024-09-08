package nebulaorm

import (
	"fmt"
	"github.com/haysons/nebulaorm/resolver"
	"github.com/haysons/nebulaorm/statement"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"net"
	"strconv"
	"time"
)

// DB will use statement.Statement to construct the nGQL statement, and then hand it over to nebula.SessionPool to execute
// the statement, and you can eventually get the result of the execution through methods such as Find Exec Pluck.  DB is
// concurrency-safe, and multiple statements can be executed by a single DB object at the same time. The nebula graph
// officially provides a SessionPool, which eliminates the need for the application layer to implement a connection pool.
// So in most cases, the application layer only needs to use a single DB instance.
// However, statement.Statement is not concurrency-safe, so don't concurrently build nGQL statements.
// NOTE: No embedded field is supported for struct, so do not use embedded field when declaring struct.
type DB struct {
	Statement   *statement.Statement
	conf        *Config
	sessionPool *nebula.SessionPool
	clone       int
}

func Open(conf *Config, opts ...ConfigOption) (*DB, error) {
	for _, o := range opts {
		o.apply(conf)
	}

	if conf.TimezoneName != "" {
		loc, err := time.LoadLocation(conf.TimezoneName)
		if err != nil {
			return nil, fmt.Errorf("nebulaorm: load timezone failed: %v", err)
		}
		conf.timezone = loc
	} else {
		conf.timezone = time.Local
	}
	resolver.SetTimezone(conf.timezone)

	hostAddr, err := parseServerAddr(conf.Addresses)
	if err != nil {
		return nil, err
	}
	poolConf, err := nebula.NewSessionPoolConf(conf.Username, conf.Password, hostAddr, conf.SpaceName, parseSessionOptions(conf)...)
	if err != nil {
		return nil, fmt.Errorf("nebulaorm: build session pool conf failed: %v", err)
	}
	pool, err := nebula.NewSessionPool(*poolConf, nebula.DefaultLogger{})
	if err != nil {
		return nil, fmt.Errorf("nebulaorm: create session pool failed: %v", err)
	}

	db := &DB{
		Statement:   statement.New(),
		conf:        conf,
		sessionPool: pool,
		clone:       1, // when clone is 1, the Statement object will be copied to ensure that the same singleton build statement does not affect each other.
	}
	return db, nil
}

func parseServerAddr(addrList []string) ([]nebula.HostAddress, error) {
	hostAddr := make([]nebula.HostAddress, 0, len(addrList))
	for _, addr := range addrList {
		host, portTmp, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("nebulaorm: parse server addr failed: %w", err)
		}
		port, err := strconv.Atoi(portTmp)
		if err != nil {
			return nil, fmt.Errorf("nebulaorm: convert server addr port failed: %w", err)
		}
		hostAddr = append(hostAddr, nebula.HostAddress{
			Host: host,
			Port: port,
		})
	}
	return hostAddr, nil
}

func parseSessionOptions(conf *Config) []nebula.SessionPoolConfOption {
	poolOptions := make([]nebula.SessionPoolConfOption, 0)
	if conf.MaxOpenConns > 0 {
		poolOptions = append(poolOptions, nebula.WithMaxSize(conf.MaxOpenConns))
	}
	if conf.MinOpenConns > 0 {
		poolOptions = append(poolOptions, nebula.WithMinSize(conf.MinOpenConns))
	}
	if conf.ConnTimeout > 0 {
		poolOptions = append(poolOptions, nebula.WithTimeOut(conf.ConnTimeout))
	}
	if conf.ConnMaxIdleTime > 0 {
		poolOptions = append(poolOptions, nebula.WithIdleTime(conf.ConnMaxIdleTime))
	}
	poolOptions = append(poolOptions, conf.nebulaSessionOpts...)
	return poolOptions
}

func (db *DB) getInstance() *DB {
	if db.clone > 0 {
		tx := &DB{conf: db.conf, sessionPool: db.sessionPool, clone: 0}
		tx.Statement = statement.New()
		return tx
	}
	return db
}

func (db *DB) Close() error {
	db.sessionPool.Close()
	return nil
}
