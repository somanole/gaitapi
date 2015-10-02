// CassandraMatchRepo
package matchrepo

import (
	"errors"
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
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
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraMatchRepo) GetUserMatch(userId string) (types.Match, error) {
	// get match for user by id
	log.Printf(fmt.Sprintf("GetUserMatch - Received userId: %v", userId))
	
	var match types.Match
	var user_id string
	var matched_user_id string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf(`SELECT user_id, matched_user_id, matched_username, 
		is_match_active, is_chat_active, timestamp 
		FROM matches WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, &matched_user_id, 
		&match.MatchedUsername, &match.IsMatchActive, &match.IsChatActive, &match.Timestamp); err != nil {
				log.Printf(fmt.Sprintf("GetUserMatch - Error: %v", err.Error()))
		} else {
			match.UserId = uuid.Parse(user_id)
			match.MatchedUserId = uuid.Parse(matched_user_id)
		}
	} else {
		log.Printf(fmt.Sprintf("GetUserMatch - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return match, err
}
