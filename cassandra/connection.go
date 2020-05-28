package cassandra

import (
	"github.com/gocql/gocql"
)

const Keyspace = "brain_debt"

func NewConnection(hosts ...string) (*gocql.Session, error) {
	cluster := gocql.NewCluster()
	cluster.Keyspace = Keyspace
	cluster.Consistency = gocql.LocalQuorum
	cluster.Hosts = hosts
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
