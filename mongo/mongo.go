package mongo

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoSVC interface {
	Create(string, string, string, string) error
	Update(string, string, string) error
	Close() error
}

func New(host, port, dbName string, collectionName string) (MongoSVC, error) {
	log.Printf("Connect to Mongo DB %s %s", host, port)
	db, err := mgo.Dial(host + ":" + port)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &mongo{db, dbName, collectionName}, nil
}

func (m *mongo) Close() error {
	m.db.Close()
	return nil
}

func (m *mongo) Create(name string, status string, org string, jobType string) error {
	log.Print("Creating mongo record.....")
	log.Print(name, status)
	job := &Job{bson.NewObjectId(), name, status, "", org, jobType, time.Now()}
	c := m.db.DB(m.dbName).C(m.collectionName)
	err := c.Insert(job)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongo) Update(name string, status string, jobLog string) error {
	log.Print("Update job instance")
	log.Print(name, status, jobLog)
	c := m.db.DB(m.dbName).C(m.collectionName)
	err := c.Update(bson.M{"name": name}, bson.M{"$set": bson.M{"status": status, "logs": jobLog}})
	if err != nil {
		return err
	}
	return nil
}

type mongo struct {
	db             *mgo.Session
	dbName         string
	collectionName string
}

type Job struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	NAME         string        `json:"name"`
	STATUS       string        `json:"status"`
	LOGS         string        `json:"logs"`
	ORGANIZATION string        `json:"organization"`
	TYPE         string        `json:"type"`
	CREATEDAT    time.Time     `json:"createdAt"`
}
