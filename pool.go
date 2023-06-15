package neo4jx

import (
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	pool  = make(map[string]neo4j.Driver)
	mutex sync.Mutex
)

func GetDriverFromPool(dsn string) (neo4j.Driver, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if v, ok := pool[dsn]; ok {
		return v, nil
	}
	addr, username, password, e := ParseDsn(dsn)
	if e != nil {
		return nil, e
	}

	v, e := neo4j.NewDriver(addr, neo4j.BasicAuth(username, password, ""))
	if e != nil {
		return nil, e
	}
	pool[dsn] = v
	return v, nil
}
