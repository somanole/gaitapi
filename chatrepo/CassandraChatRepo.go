// CassandraChatRepo
package chatrepo

import (
	"fmt"
	"log"
	"time"
	"github.com/somanole/gaitapi/types"
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

func init() {
	New = NewCassandraChatRepo
}

func getCqlSession() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "gait"
	
	s,_ := cluster.CreateSession()
	
	return s
}

func (repo *CassandraChatRepo) CreateChat(c types.Chat) error {
	//insert chat in chats
	
	var err error
	err = nil
	
	sql := fmt.Sprintf(`INSERT INTO chats (sender_id, receiver_id, 
	is_chat_active, is_chat_blocked_by_sender, is_chat_blocked_by_receiver, receiver_username, 
	last_message, timestamp) VALUES (%v, %v, %v, %v, %v, '%v', '%v', %v)`, 
	c.SenderId, c.ReceiverId, c.IsChatActive, c.IsChatBlockedBySender, c.IsChatBlockedByReceiver, c.ReceiverUsername, c.LastMessage, c.Timestamp)
						
	log.Printf(sql)
	
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.CreateChat() - Error: %v", err.Error()))
	} 
	
    return err
}

func (repo *CassandraChatRepo) GetUserActiveChats(userId string) (types.Chats, error) {
	// get active chats for user by id
	log.Printf(fmt.Sprintf("CassandraChatRepo.GetUserActiveChats() - Received userId: %v", userId))
	
	var chats types.Chats
	var err error
	err = nil
	
	var sender_id, receiver_id, receiver_username, last_message string
	var is_chat_active, is_chat_blocked_by_sender, is_chat_blocked_by_receiver  bool
	var timestamp int64
			
	sql := fmt.Sprintf(`SELECT sender_id, receiver_id, receiver_username, 
	is_chat_active, is_chat_blocked_by_sender, is_chat_blocked_by_receiver, last_message, timestamp 
	FROM chats WHERE sender_id = %v`, userId)
			
	log.Printf(sql)
			
	iter := session.Query(sql).Iter()
	for iter.Scan(&sender_id, &receiver_id, &receiver_username,
	&is_chat_active, &is_chat_blocked_by_sender, &is_chat_blocked_by_receiver, &last_message, &timestamp) {
		if is_chat_active {
			chats = append(chats, types.Chat{SenderId: uuid.Parse(sender_id), ReceiverId: uuid.Parse(receiver_id), 
			ReceiverUsername: receiver_username, IsChatActive: is_chat_active,
			IsChatBlockedBySender: is_chat_blocked_by_sender, IsChatBlockedByReceiver: is_chat_blocked_by_receiver,
			LastMessage: last_message, Timestamp: timestamp})
		}
	}
	if err = iter.Close(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.GetUserActiveChats() - Error: %v", err.Error()))
	}
		
	if chats != nil {
		sort.Sort(ByTimestamp(chats))
	}
	
	return chats, err
}

func (repo *CassandraChatRepo) BlockChat(senderId string, receiverId string) error {
    // block chat in chats
	var err error
	err = nil
	
	timestamp := int64(time.Now().UTC().Unix())
			
	sql := fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_sender = true, timestamp = %v WHERE sender_id = %v
	AND receiver_id = %v`, timestamp, senderId, receiverId)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.BlockChat() - Error: %v", err.Error()))
	} 
		
	sql = fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_receiver = true, timestamp = %v WHERE sender_id = %v
	AND receiver_id = %v`, timestamp, receiverId, senderId)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.BlockChat() - Error: %v", err.Error()))
	} 

    return err
}

func (repo *CassandraChatRepo) UnblockChat(senderId string, receiverId string) error {
    // unblock chat in chats
	var err error
	err = nil
	
	timestamp := int64(time.Now().UTC().Unix())
			
	sql := fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_sender = false, timestamp = %v WHERE sender_id = %v
	AND receiver_id = %v`, timestamp, senderId, receiverId)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.UnblockChat() - Error: %v", err.Error()))
	} 
		
	sql = fmt.Sprintf(`UPDATE chats SET is_chat_blocked_by_receiver = false, timestamp = %v WHERE sender_id = %v
	AND receiver_id = %v`, timestamp, receiverId, senderId)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.UnblockChat() - Error: %v", err.Error()))
	} 

    return err
}

func (repo *CassandraChatRepo) DeleteChat(senderId string, receiverId string) error {
    // delete chat in chats
	var err error
	err = nil
	
	timestamp := int64(time.Now().UTC().Unix())
			
	sql := fmt.Sprintf(`UPDATE chats SET is_chat_active = false, timestamp = %v WHERE sender_id = %v
	AND receiver_id = %v`, timestamp, senderId, receiverId)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.DeleteChat() - Error: %v", err.Error()))
	} 

    return err
}

func (repo *CassandraChatRepo) UpdateLastMessageChat(senderId string, receiverId string, message string) error {
    // update last message in chat
	var err error
	err = nil
	
	timestamp := int64(time.Now().UTC().Unix())
			
	sql := fmt.Sprintf(`UPDATE chats SET last_message = '%v', timestamp = %v WHERE sender_id = %v
	AND receiver_id = %v`, message, timestamp, senderId, receiverId)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.UpdateLastMessageChat() - Error: %v", err.Error()))
	} 
		
	sql = fmt.Sprintf(`UPDATE chats SET last_message = '%v', timestamp = %v WHERE sender_id = %v
	AND receiver_id = %v`, message, timestamp, receiverId, senderId)
						
	log.Printf(sql)
	if err = session.Query(sql).Exec(); err != nil {
		log.Printf(fmt.Sprintf("CassandraChatRepo.UpdateLastMessageChat() - Error: %v", err.Error()))
	} 

    return err
}
