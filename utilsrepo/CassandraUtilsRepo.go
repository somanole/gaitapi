// CassandraUtilsRepo
package utilsrepo

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
	CassandraUtilsRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraUtilsRepo() UtilsRepo {
	return &CassandraUtilsRepo{}
}

func init() {
	New = NewCassandraUtilsRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraActivityRepo) CheckUserPassword(userId string, rpassword string) error {
	// check password for user
	log.Printf(fmt.Sprintf("CheckUserPassword - Received userId: %v", userId))
	log.Printf(fmt.Sprintf("CheckUserPassword - Received password: %v", rpassword))

	var password string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf("SELECT password FROM users WHERE user_id = %v LIMIT 1", userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&password); err != nil {
				log.Printf(fmt.Sprintf("CheckUserPassword - Error: %v", err.Error()))
		} else if password != rpassword{
			log.Printf(fmt.Sprintf("CheckUserPassword - password received does not match password on record: %v != %v", rpassword, password))
			err = errors.New("401")
		}
	} else {
		log.Printf(fmt.Sprintf("CheckUserPassword - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return err
}

func (repo *CassandraActivityRepo) CheckIfUserExists(userId string) error {
	// check if user exists
	log.Printf(fmt.Sprintf("CheckIfUserExists - Received userId: %v", userId))

	var password string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf("SELECT password FROM users WHERE user_id = %v LIMIT 1", userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&password); err != nil {
				log.Printf(fmt.Sprintf("CheckIfUserExists - Error: %v", err.Error()))
		}
	} else {
		log.Printf(fmt.Sprintf("CheckIfUserExists - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return err
}

func (repo *CassandraActivityRepo) GetUserUsername(userId string) string, error {
	// get user username
	log.Printf(fmt.Sprintf("GetUserUsername - Received userId: %v", userId))

	var username string
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		sql := fmt.Sprintf("SELECT username FROM users WHERE user_id = %v LIMIT 1", userId)
		
		log.Printf(sql)
		
		if err = session.Query(sql).Scan(&username); err != nil {
				log.Printf(fmt.Sprintf("GetUserUsername - Error: %v", err.Error()))
		} else if username == "" {
			log.Printf(fmt.Sprintf("GetUserUsername - username on record is blank: %v", username))
			err = errors.New("404")
		}
	} else {
		log.Printf(fmt.Sprintf("GetUserUsername - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return username, err
}
