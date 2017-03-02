package cassandra

import (
	"fmt"
	"log"
	"strings"

	"github.com/dollarshaveclub/go-lib/set"
	"github.com/gocql/gocql"
)

// CTable is a Cassandra table definition
type CTable struct {
	Name    string
	Columns []string
	Options string
}

// UDT is a user-defined type
type UDT struct {
	Name    string
	Columns []string
}

// CreateTable creates a table
func CreateTable(c *gocql.ClusterConfig, t CTable) error {
	qs := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v ( %v ) %v ;", t.Name, strings.Join(t.Columns, ", "), t.Options)
	s, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("Creating table %v\n", t.Name)
	return s.Query(qs).Exec()
}

// CreateUDT creates a user-defined type
func CreateUDT(c *gocql.ClusterConfig, u UDT) error {
	qs := fmt.Sprintf("CREATE TYPE IF NOT EXISTS %v ( %v );", u.Name, strings.Join(u.Columns, ", "))
	s, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("Creating UDT %v\n", u.Name)
	return s.Query(qs).Exec()
}

// CreateRequiredTypes ensures all the types passed in are created if necessary
func CreateRequiredTypes(c *gocql.ClusterConfig, rt []UDT) error {
	rtn := []string{}
	etn := []string{}
	rtm := map[string]UDT{}
	for _, u := range rt {
		rtn = append(rtn, u.Name)
		rtm[u.Name] = u
	}
	rts := set.NewStringSet(rtn)
	s, err := c.CreateSession()
	if err != nil {
		return err
	}
	q := `SELECT type_name FROM system.schema_usertypes WHERE keyspace_name = '%v';`
	q = fmt.Sprintf(q, c.Keyspace)
	iter := s.Query(q).Iter()
	for n := ""; iter.Scan(&n); {
		etn = append(etn, n)
	}
	if err := iter.Close(); err != nil {
		return err
	}
	ets := set.NewStringSet(etn)
	missing := rts.Difference(ets).Items()
	if len(missing) > 0 {
		for _, mt := range missing {
			err := CreateUDT(c, rtm[mt])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateRequiredTables ensures all the tables passed in are created if necessary
func CreateRequiredTables(c *gocql.ClusterConfig, rt []CTable) error {
	tl, err := GetTables(c)
	if err != nil {
		return err
	}

	tm := map[string]CTable{}
	rtl := []string{}
	for _, v := range rt {
		tm[v.Name] = v
		rtl = append(rtl, v.Name)
	}
	tset := set.NewStringSet(tl)
	rset := set.NewStringSet(rtl)
	diff := rset.Difference(tset)
	missing := diff.Items()
	if len(missing) > 0 {
		for _, t := range missing {
			ts := tm[t]
			err = CreateTable(c, ts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateKeyspace creates a keyspace if necessary
// ks -> keyspace name
// rs -> replication strategy class
// rf -> replication factor
func CreateKeyspace(c *gocql.ClusterConfig, ks string, rs string, rf int) error {
	kis, err := GetKeyspaces(c)
	if err != nil {
		return err
	}
	kss := set.NewStringSet(kis)
	if !kss.Contains(ks) {
		log.Printf("Creating keyspace: %v\n", ks)
		c.Keyspace = ""
		s, err := c.CreateSession()
		if err != nil {
			return err
		}
		defer s.Close()
		if rs == "" {
			rs = "SimpleStrategy"
		}
		err = s.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH REPLICATION = {'class': '%v', 'replication_factor': %v};", ks, rs, rf)).Exec()
		if err != nil {
			return err
		}
	}
	c.Keyspace = ks
	return nil
}

// CreateKeyspaceWithNetworkTopologyStrategy creates a keyspace if necesary with NetworkTopologyStrategy
// ks -> keyspace name
// rfmap -> map of datacenter name to replication factor for that DC
func CreateKeyspaceWithNetworkTopologyStrategy(c *gocql.ClusterConfig, ks string, rfmap map[string]uint) error {
	kis, err := GetKeyspaces(c)
	if err != nil {
		return err
	}
	kss := set.NewStringSet(kis)
	if !kss.Contains(ks) {
		log.Printf("Creating keyspace: %v\n", ks)
		c.Keyspace = ""
		s, err := c.CreateSession()
		if err != nil {
			return err
		}
		defer s.Close()
		q := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH REPLICATION = {'class': '%v', ", ks, "NetworkTopologyStrategy")
		rfsl := []string{}
		for dc, rf := range rfmap {
			rfsl = append(rfsl, fmt.Sprintf("'%v' : %v", dc, rf))
		}
		q = fmt.Sprintf("%v%v};", q, strings.Join(rfsl, ", "))
		err = s.Query(q).Exec()
		if err != nil {
			return err
		}
	}
	c.Keyspace = ks
	return nil
}

// DropKeyspace deletes a keyspace and all data associated with it
func DropKeyspace(c *gocql.ClusterConfig, ks string) error {
	c.Keyspace = ""
	s, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("Dropping keyspace: %v\n", ks)
	err = s.Query(fmt.Sprintf("DROP KEYSPACE IF EXISTS %v\n", ks)).Exec()
	if err != nil {
		return err
	}
	return nil
}
