package ghiccup

import (
	"encoding/json"
	"log"
	"unicode/utf8"

	"github.com/boltdb/bolt"
)

/* MessageBuffer is a buffer containing messages
 */
type MessageBuffer struct {
	bucketName string
	db         *bolt.DB
}

/* Create new BoltDB based buffer
 */
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

func (buffer MessageBuffer) Close() error {
	return buffer.db.Close()
}

func (buffer MessageBuffer) nextKey(b *bolt.Bucket) []byte {
	buf := make([]byte, 3)
	key, _ := b.NextSequence()
	utf8.EncodeRune(buf, rune(key))
	return buf
}
func serialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}
func deserialize(bytes []byte, objMaker func() interface{}) (interface{}, error) {
	obj := objMaker()
	err := json.Unmarshal(bytes, obj)
	return obj, err
}
func (buffer MessageBuffer) Add(obj interface{}) error {
	tx, err := buffer.db.Begin(true)
	defer tx.Commit()
	b := tx.Bucket([]byte(buffer.bucketName))
	line, err := serialize(obj)
	if err != nil {
		return err
	} else {
		key := buffer.nextKey(b)
		log.Println("Adding entry with key", key)
		err := b.Put(key, line)
		tx.Commit()
		return err
	}
}

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
