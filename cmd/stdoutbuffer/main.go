/*
Allows stdin input to be buffered if writing to stdout fails
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"unicode/utf8"
)

func main() {
	bufferFile := flag.String("bufferfile", "stdoutbuffer.db", "Database file for buffer")
	bucketName := flag.String("bucket", "buffer", "Name of bucket in which to persist values")

	db, err := bolt.Open(*bufferFile, 0600, nil)
	defer db.Close()

	tx, err := db.Begin(true)
	b, err := tx.CreateBucketIfNotExists([]byte(*bucketName))
	defer tx.Commit()

	replayExistingLines(b)

	if err != nil {
		log.Fatal(err)
	}

	bio := bufio.NewReader(os.Stdin)
	for {
		if processLine(bio, b) != nil {
			log.Println("Error reading from stdin. Stopping consumption")

			break
		}

	}
}

func processLine(bio *bufio.Reader, b *bolt.Bucket) error {
	line, _, readErr := bio.ReadLine()
	if readErr != nil {
		return readErr
	}
	_, outErr := fmt.Println(string(line))

	if outErr != nil {
		writeToBuffer(b, line)
	}
	return nil
}

func writeToBuffer(bucket *bolt.Bucket, line []byte) {
	log.Println("Buffering line to file")
	buf := make([]byte, 3)
	key, _ := bucket.NextSequence()
	utf8.EncodeRune(buf, rune(key))
	bucket.Put(buf, line)
}

func replayExistingLines(bucket *bolt.Bucket) {
	log.Println("Replaying buffered lines...")
	bucket.ForEach(func(k []byte, v []byte) error {
		_, outErr := fmt.Println(string(v))

		if outErr == nil {
			bucket.Delete(k)
			return nil
		} else {
			log.Println("Could not replay log line properly. Not popping")
			return outErr
		}
	})
}
