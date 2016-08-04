package common

import (
	"log"
	"os"

	"gopkg.in/mgo.v2/bson"
)

type Action struct {
	ID                    bson.ObjectId `bson:"_id,omitempty"`
	NameSpace, User, Code string
}

// MustGetenv gets an env variable if it exists, else it panics.
func MustGetenv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		log.Panicln("Missing env variable:", name)
	}
	return env
}
