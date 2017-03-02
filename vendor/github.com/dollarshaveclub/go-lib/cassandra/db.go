package cassandra

import "github.com/gocql/gocql"

// GetTables returns all tables in configured keyspace
func GetTables(c *gocql.ClusterConfig) ([]string, error) {
	tables := []string{}
	s, err := c.CreateSession()
	if err != nil {
		return tables, err
	}
	defer s.Close()
	q := s.Query("SELECT columnfamily_name FROM system.schema_columnfamilies WHERE keyspace_name = ?;", c.Keyspace).Iter()
	var tn string
	for q.Scan(&tn) {
		tables = append(tables, tn)
	}
	return tables, q.Close()
}

// GetKeyspaces returns all extant keyspaces
func GetKeyspaces(c *gocql.ClusterConfig) ([]string, error) {
	kss := []string{}
	s, err := c.CreateSession()
	if err != nil {
		return kss, err
	}
	defer s.Close()
	q := s.Query("SELECT keyspace_name FROM system.schema_keyspaces;").Iter()
	var kn string
	for q.Scan(&kn) {
		kss = append(kss, kn)
	}
	return kss, q.Close()
}
