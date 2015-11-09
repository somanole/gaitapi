// CassandraMatchRepo
package matchrepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/constants"
	"github.com/gocql/gocql"
	"code.google.com/p/go-uuid/uuid"
)

type (
	CassandraMatchRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraMatchRepo() MatchRepo {
	return &CassandraMatchRepo{}
}

func init() {
	New = NewCassandraMatchRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster(constants.CASSANDRA_MASTER)
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraMatchRepo) GetUserMatch(userId string) (types.Match, error) {
	// get match for user by id
	log.Printf(fmt.Sprintf("CassandraMatchRepo.GetUserMatch() - Received userId: %v", userId))
	
	var match types.Match
	var user_id string
	var matched_user_id string
	var err error
	err = nil
	
	sql := fmt.Sprintf(`SELECT user_id, matched_user_id, 
	matched_username, timestamp 
	FROM matches WHERE user_id = %v LIMIT 1`, userId)
		
	log.Printf(sql)
		
	if err = session.Query(sql).Scan(&user_id, &matched_user_id, 
	&match.MatchedUsername, &match.Timestamp); err != nil {
		log.Printf(fmt.Sprintf("CassandraMatchRepo.GetUserMatch() - Error: %v", err.Error()))
	} else {
		match.UserId = uuid.Parse(user_id)
		match.MatchedUserId = uuid.Parse(matched_user_id)
	}

	return match, err
}

func (repo *CassandraMatchRepo) CreateMatch(m types.Match) error {
    // insert match in matches
	var err error
	err = nil
	
	sql := fmt.Sprintf(`INSERT INTO matches (user_id, matched_user_id, matched_username, 
	timestamp) VALUES (%v, %v, '%v', %v)`, 
	m.UserId, m.MatchedUserId, m.MatchedUsername, m.Timestamp)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraMatchRepo.CreateMatch() - Error: %v", err.Error()))
	} else {
		sql = fmt.Sprintf(`INSERT INTO matches_by_matched_user_id (user_id, matched_user_id, matched_username, 
		timestamp) VALUES (%v, %v, '%v', %v)`, 
		m.UserId, m.MatchedUserId, m.MatchedUsername, m.Timestamp)
							
		log.Printf(sql)
		if err = session.Query(sql).Exec(); err != nil {
			log.Printf(fmt.Sprintf("CassandraMatchRepo.CreateMatch() - Error: %v", err.Error()))
		} 
	} 
	
    return err
}
