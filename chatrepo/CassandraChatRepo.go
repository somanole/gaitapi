// CassandraChatRepo
package chatrepo

import (
	"errors"
	"fmt"
	"log"
	"time"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/utilsrepo"
	"github.com/gocql/gocql"
	"code.google.com/p/go-uuid/uuid"
	"sort"
)

// ByTimestamp implements sort.Interface for []Chat based on
// the Timestamp field.
type ByTimestamp []types.Chat

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }

type (
	CassandraChatRepo struct {}
)

var session *gocql.Session = getCqlSession()

func NewCassandraChatRepo() ChatRepo {
	return &CassandraChatRepo{}
}

var utilsRepo utilsrepo.UtilsRepo

func init() {
	New = NewCassandraChatRepo
	utilsRepo = utilsrepo.New()
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraChatRepo) GetUserActiveChats(userId string) (types.Chats, error) {
	// get active chats for user by id
	log.Printf(fmt.Sprintf("GetUserActiveChats - Received userId: %v", userId))
	
	var chats types.Chats
	var err error
	err = nil
	
	if uuid.Parse(userId) != nil {
		var sender_id, receiver_id, receiver_username string
		var is_chat_active, is_chat_blocked bool
		var timestamp int64
			
		sql := fmt.Sprintf(`SELECT sender_id, receiver_id, receiver_username, 
		is_chat_active, is_chat_blocked, timestamp 
		FROM chats WHERE user_id = %v`, userId)
			
		log.Printf(sql)
			
		iter := session.Query(sql).Iter()
		for iter.Scan(&user_id, &matched_user_id, &matched_username,
		&is_match_active, &is_chat_active, &is_chat_blocked, &timestamp) {
			if is_chat_active {
				chats = append(chats, types.Chat{UserId: uuid.Parse(user_id), MatchedUserId: uuid.Parse(matched_user_id), 
				MatchedUserName: matched_username, IsMatchActive: is_match_active, IsChatActive: is_chat_active,
				IsChatBlocked: is_chat_blocked Timestamp: timestamp})
			}
		}
		if err = iter.Close(); err != nil {
			log.Printf(fmt.Sprintf("GetUserActiveChats - Error: %v", err.Error()))
		}
	} else {
		log.Printf(fmt.Sprintf("GetUserMatch - Received userId: %v is not UUID", userId))
		err = errors.New("not uuid")
	}
	
	return chats, err
}

func (repo *CassandraChatRepo) BlockChat(senderId string, receiverId string) error {
    // block chat in chats
	var m types.Match
	var err error
	var receiverUsername string
	err = nil
	
	err = utilsRepo.CheckIfUserExists(senderId)
	
	if err == nil {
		err = utilsRepo.CheckIfUserExists(receiverId)
	}

	if err == nil {
		timestamp := int64(time.Now().UTC().Unix())
			
		sql := fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_sender = true, timestamp = %v WHERE sender_id = %v
		AND receiver_id = %v`, timestamp, senderId, receiverId)
						
		log.Printf(sql)
		if err = session.Query(sql).Exec(); err != nil {
			log.Printf(fmt.Sprintf("BlockChat - Error: %v", err.Error()))
		} 
		
		sql = fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_receiver = true, timestamp = %v WHERE sender_id = %v
		AND receiver_id = %v`, timestamp, receiverId, senderId)
						
		log.Printf(sql)
		if err = session.Query(sql).Exec(); err != nil {
			log.Printf(fmt.Sprintf("BlockChat - Error: %v", err.Error()))
		} 
	}

    return err
}

func (repo *CassandraChatRepo) UnblockChat(senderId string, receiverId string) error {
    // unblock chat in chats
	var m types.Match
	var err error
	var receiverUsername string
	err = nil
	
	err = utilsRepo.CheckIfUserExists(senderId)
	
	if err == nil {
		err = utilsRepo.CheckIfUserExists(receiverId)
	}

	if err == nil {
		timestamp := int64(time.Now().UTC().Unix())
			
		sql := fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_sender = false, timestamp = %v WHERE sender_id = %v
		AND receiver_id = %v`, timestamp, senderId, receiverId)
						
		log.Printf(sql)
		if err = session.Query(sql).Exec(); err != nil {
			log.Printf(fmt.Sprintf("UnblockChat - Error: %v", err.Error()))
		} 
		
		sql = fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_receiver = false, timestamp = %v WHERE sender_id = %v
		AND receiver_id = %v`, timestamp, receiverId, senderId)
						
		log.Printf(sql)
		if err = session.Query(sql).Exec(); err != nil {
			log.Printf(fmt.Sprintf("UnblockChat - Error: %v", err.Error()))
		} 
	}

    return err
}

func (repo *CassandraChatRepo) DeleteChat(senderId string, receiverId string) error {
    // delete chat in chats
	var m types.Match
	var err error
	var receiverUsername string
	err = nil
	
	err = utilsRepo.CheckIfUserExists(senderId)
	
	if err == nil {
		err = utilsRepo.CheckIfUserExists(receiverId)
	}

	if err == nil {
		timestamp := int64(time.Now().UTC().Unix())
			
		sql := fmt.Sprintf(`UPDATE chats SET is_chat_active = false, timestamp = %v WHERE sender_id = %v
		AND receiver_id = %v`, timestamp, senderId, receiverId)
						
		log.Printf(sql)
		if err = session.Query(sql).Exec(); err != nil {
			log.Printf(fmt.Sprintf("DeleteChat - Error: %v", err.Error()))
		} 
	}

    return err
}
