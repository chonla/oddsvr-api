package database

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Database struct {
	connection string
	session    *mgo.Session
	db         *mgo.Database
}

func NewDatabase(connection, dbname string) (*Database, error) {
	session, e := mgo.Dial(connection)
	if e != nil {
		return nil, e
	}

	return &Database{
		connection: connection,
		session:    session,
		db:         session.DB(dbname),
	}, nil
}

func (db *Database) Get(collection string, filter, output interface{}) error {
	q := db.db.C(collection).Find(filter).Limit(1)
	return q.Select(bson.M{}).One(output)
}

func (db *Database) List(collection string, filter interface{}, order []string, output interface{}) error {
	return db.db.C(collection).Find(filter).Sort(order...).All(output)
}

func (db *Database) InsertBulk(collection string, data []interface{}) error {
	return db.db.C(collection).Insert(data...)
}

func (db *Database) Upsert(collection string, filter, data interface{}) error {
	_, e := db.db.C(collection).Upsert(filter, data)
	return e
}

func (db *Database) Insert(collection string, data interface{}) error {
	e := db.db.C(collection).Insert(data)
	return e
}

func (db *Database) Replace(collection string, filter, data interface{}) error {
	e := db.db.C(collection).Update(filter, data)
	return e
}

func (db *Database) Push(collection string, filter, data interface{}) error {
	e := db.db.C(collection).Update(filter, bson.M{
		"$push": data,
	})
	return e
}

func (db *Database) Has(collection string, filter interface{}) bool {
	q := db.db.C(collection).Find(filter).Limit(1)
	count, e := q.Count()
	if e != nil {
		return false
	}
	return (count > 0)
}
