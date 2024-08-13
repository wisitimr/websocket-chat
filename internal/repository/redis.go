package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"websocket-chat/internal/model"
	mRepo "websocket-chat/internal/model/repository"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type redisRepository struct {
	redisClient *redis.Client
	logger      *logrus.Logger
}

func InitRedisRepository(redisClient *redis.Client, logger *logrus.Logger) mRepo.RedisRepository {
	return &redisRepository{
		redisClient: redisClient,
		logger:      logger,
	}
}

func (r redisRepository) RegisterNewUser(username, password string) error {
	// redis-cli
	// SYNTAX: SET key value
	// SET username password
	// register new username:password key-value pair
	err := r.redisClient.Set(context.Background(), username, password, 0).Err()
	if err != nil {
		log.Println("error while adding new user", err)
		return err
	}

	// redis-cli
	// SYNTAX: SADD key value
	// SADD users username
	err = r.redisClient.SAdd(context.Background(), userSetKey(), username).Err()
	if err != nil {
		log.Println("error while adding user in set", err)
		// redis-cli
		// SYNTAX: DEL key
		// DEL username
		// drop the registered user
		r.redisClient.Del(context.Background(), username)

		return err
	}

	return nil
}

func (r redisRepository) IsUserExist(username string) bool {
	// redis-cli
	// SYNTAX: SISMEMBER key value
	// SISMEMBER users username
	return r.redisClient.SIsMember(context.Background(), userSetKey(), username).Val()
}

func (r redisRepository) IsUserAuthentic(username, password string) error {
	// redis-cli
	// SYNTAX: GET key
	// GET username
	p := r.redisClient.Get(context.Background(), username).Val()

	if !strings.EqualFold(p, password) {
		return fmt.Errorf("invalid username or password")
	}

	return nil
}

// UpdateContactList add contact to username's contact list
// if not present or update its timestamp as last contacted
func updateContactList(redisClient *redis.Client, username, contact string) error {
	zs := &redis.Z{Score: float64(time.Now().Unix()), Member: contact}

	// redis-cli SCORE is always float or int
	// SYNTAX: ZADD key SCORE MEMBER
	// ZADD contacts:username 1661360942123 contact
	err := redisClient.ZAdd(context.Background(),
		contactListZKey(username),
		zs,
	).Err()

	if err != nil {
		log.Println("error while updating contact list. username: ",
			username, "contact:", contact, err)
		return err
	}

	return nil
}

func (r redisRepository) CreateChat(c *model.Chat) (string, error) {
	chatKey := chatKey()
	fmt.Println("chat key", chatKey)

	by, _ := json.Marshal(c)

	// redis-cli
	// SYNTAX: JSON.SET key $ json_in_string
	// JSON.SET chat#1661360942123 $ '{"from": "sun", "to":"earth","message":"good morning!"}'
	res, err := r.redisClient.Do(
		context.Background(),
		"JSON.SET",
		chatKey,
		"$",
		string(by),
	).Result()

	if err != nil {
		log.Println("error while setting chat json", err)
		return "", err
	}

	log.Println("chat successfully set", res)

	// add contacts to both user's contact list
	err = updateContactList(r.redisClient, c.From, c.To)
	if err != nil {
		log.Println("error while updating contact list of", c.From)
	}

	err = updateContactList(r.redisClient, c.To, c.From)
	if err != nil {
		log.Println("error while updating contact list of", c.To)
	}

	return chatKey, nil
}

func (r redisRepository) CreateFetchChatBetweenIndex() {
	res, err := r.redisClient.Do(context.Background(),
		"FT.CREATE",
		chatIndex(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	fmt.Println(res, err)
}

func (r redisRepository) FetchChatBetween(username1, username2, fromTS, toTS string) ([]model.Chat, error) {
	// redis-cli
	// SYNTAX: FT.SEARCH index query
	// FT.SEARCH idx#chats '@from:{user2|user1} @to:{user1|user2} @timestamp:[0 +inf] SORTBY timestamp DESC'
	query := fmt.Sprintf("@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]",
		username1, username2, username1, username2, fromTS, toTS)

	res, err := r.redisClient.Do(context.Background(),
		"FT.SEARCH",
		chatIndex(),
		query,
		"SORTBY", "timestamp", "DESC",
	).Result()

	if err != nil {
		return nil, err
	}

	// deserialise redis data to map
	data := Deserialise(res)

	// deserialise data map to chat
	chats := DeserialiseChat(data)
	return chats, nil
}

// FetchContactList of the user. It includes all the messages sent to and received by contact
// It will return a sorted list by last activity with a contact
func (r redisRepository) FetchContactList(username string) ([]model.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListZKey(username),
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}

	// redis-cli
	// SYNTAX: ZRANGE key from_index to_index REV WITHSCORES
	// ZRANGE contacts:username 0 -1 REV WITHSCORES
	res, err := r.redisClient.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()

	if err != nil {
		log.Println("error while fetching contact list. username: ",
			username, err)
		return nil, err
	}

	contactList := DeserialiseContactList(res)

	return contactList, nil
}

func Deserialise(res interface{}) []model.Document {
	switch v := res.(type) {
	case []interface{}:
		if len(v) > 1 {
			total := len(v) - 1
			var docs = make([]model.Document, 0, total/2)

			for i := 1; i <= total; i = i + 2 {
				arrOfValues := v[i+1].([]interface{})
				value := arrOfValues[len(arrOfValues)-1].(string)

				// add _id in the response
				doc := model.Document{
					ID:      v[i].(string),
					Payload: []byte(value),
					Total:   v[0].(int64),
				}

				docs = append(docs, doc)
			}
			return docs
		}
	default:
		log.Printf("different response type otherthan []interface{}. type: %T", res)
		return nil
	}

	return nil
}

func DeserialiseChat(docs []model.Document) []model.Chat {
	chats := []model.Chat{}
	for _, doc := range docs {
		var c model.Chat
		json.Unmarshal(doc.Payload, &c)

		c.ID = doc.ID
		chats = append(chats, c)
	}

	return chats
}

func DeserialiseContactList(contacts []redis.Z) []model.ContactList {
	contactList := make([]model.ContactList, 0, len(contacts))

	// improvement tip: use switch to get type of contact.Member
	// handle unknown type accordingly
	for _, contact := range contacts {
		contactList = append(contactList, model.ContactList{
			Username:     contact.Member.(string),
			LastActivity: int64(contact.Score),
		})
	}

	return contactList
}

func userSetKey() string {
	return "users"
}

func sessionKey(client string) string {
	return "session#" + client
}

func chatKey() string {
	return fmt.Sprintf("chat#%d", time.Now().UnixMilli())
}

func chatIndex() string {
	return "idx#chats"
}

func contactListZKey(username string) string {
	return "contacts:" + username
}
