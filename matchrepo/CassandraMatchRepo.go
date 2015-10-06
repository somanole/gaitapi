// CassandraMatchRepo
package matchrepo

import (
	"errors"
	"fmt"
	"log"
	"time"
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
		sql := fmt.Sprintf(`SELECT user_id, matched_user_id, 
		matched_username, timestamp 
		FROM matches WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&user_id, &matched_user_id, 
		&match.MatchedUsername, &match.Timestamp); err != nil {
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

func (repo *CassandraMatchRepo) CreateMatch(userId string, matchedUserId string) error {
    // insert match in matches
	var m types.Match
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil && uuid.Parse(matchedUserId) != nil {
		sql := fmt.Sprintf(`SELECT email from users WHERE user_id = %v LIMIT 1`, userId)
		
		log.Printf(sql)
		var email string
		if err = session.Query(sql).Scan(&email); err != nil {
				log.Printf(fmt.Sprintf("CreateMatch - Error: %v", err.Error()))
		}
		
		sql = fmt.Sprintf(`SELECT username from users WHERE user_id = %v LIMIT 1`, matchedUserId)
		
		log.Printf(sql)
		var matchedUsername string
		if err = session.Query(sql).Scan(&matchedUsername); err != nil {
				log.Printf(fmt.Sprintf("CreateMatch - Error: %v", err.Error()))
		}
		
		if err == nil {
			m.UserId = uuid.Parse(userId)
			m.MatchedUserId = uuid.Parse(matchedUserId)
			m.MatchedUsername = matchedUsername
			m.Timestamp = int64(time.Now().UTC().Unix())
			
			sql := fmt.Sprintf(`INSERT INTO matches (user_id, matched_user_id, matched_username, 
			timestamp) VALUES (%v, %v, '%v', %v)`, 
			m.UserId, m.MatchedUserId, m.MatchedUsername, m.Timestamp)
						
			log.Printf(sql)
			if err = session.Query(sql).Exec(); err != nil {
				log.Printf(fmt.Sprintf("CreateMatch - Error: %v", err.Error()))
			} 
		}
	} else {
		log.Printf(fmt.Sprintf("CreateMatch - User | MatchedUser Id: %v | %v is not UUID", userId, matchedUserId))
		err = errors.New("not uuid")
	}
	
    return err
}