package ghiccup

import (
	"encoding/json"
	"log"
	"unicode/utf8"

	"github.com/boltdb/bolt"
)

// MessageBuffer provides a queue-like interface for locally buffered messages
type MessageBuffer struct {
	bucketName string
	db         *bolt.DB
}

// NewBuffer creates a new MessageBuffer object backed by a BoltDB database bucket
func NewBuffer(fileName string) (*MessageBuffer, error) {
	buf := new(MessageBuffer)
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, err
	}
	buf.db = db
	tx, err := db.Begin(true)
	buf.bucketName = "buffer"
	_, err = tx.CreateBucketIfNotExists([]byte(buf.bucketName))
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return buf, nil
}

// Close releases resources held by the buffer.
func (buffer MessageBuffer) Close() error {
	return buffer.db.Close()
}

// Add pushes a message to the queue
func (buffer MessageBuffer) Add(obj interface{}) error {
	tx, err := buffer.db.Begin(true)
	defer tx.Commit()
	b := tx.Bucket([]byte(buffer.bucketName))
	line, err := serialize(obj)
	if err != nil {
		return err
	}
	key := buffer.nextKey(b)
	log.Println("Adding entry with key", key)
	err = b.Put(key, line)
	if err != nil {
		return err
	}
	tx.Commit()
	return err

}

// ForEach allows consumers to iterate over all messages currently in the queue
func (buffer MessageBuffer) ForEach(fn func(v interface{}) error, objMaker func() interface{}) error {
	tx, err := buffer.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Commit()
	b := tx.Bucket([]byte(buffer.bucketName))

	return b.ForEach(func(k []byte, v []byte) error {
		obj, err := deserialize(v, objMaker)

		if err == nil {
			fn(obj)
			b.Delete(k)
		} else {
			log.Fatal(err)
		}
		return err
	})
}

// Generates a new key in the bucket
func (buffer MessageBuffer) nextKey(b *bolt.Bucket) []byte {
	buf := make([]byte, 3)
	key, _ := b.NextSequence()
	utf8.EncodeRune(buf, rune(key))
	return buf
}

// Returns serialization of an object
func serialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// Deserialize a []byte into the specified object
func deserialize(bytes []byte, objMaker func() interface{}) (interface{}, error) {
	obj := objMaker()
	err := json.Unmarshal(bytes, obj)
	return obj, err
}
