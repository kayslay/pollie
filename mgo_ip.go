package pollie

import (
	"context"
	"log"
	"pollie/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var (
	DefaultMgoIpInterfacer IPInterfacer
	initialized            bool
)

// MgoIPInterfacer implements IPInterfacer
type mgoIPInterfacer struct {
	mgo config.MgoFn
}

// InitMgoIPer initialize the default mongo ip interfacer
func InitMgoIPer(mgo config.MgoFn) {

	if initialized {
		return
	}

	// set defaultMgoIPInterfacer
	DefaultMgoIpInterfacer = mgoIPInterfacer{mgo: mgo}
	initialized = true
}

// Get get the ip from a mongo db
func (m mgoIPInterfacer) Get(ip string) (IPInfo, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	c := m.mgo(MongoIPCollection)

	var ipInfo IPInfo

	err := c.FindOne(ctx, bson.M{"ip": ip}).Decode(&ipInfo)
	log.Println("using mongo as interfacer")
	return ipInfo, err
}

// Set set the ip info for an ip
func (m mgoIPInterfacer) Set(ip IPInfo) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	c := m.mgo(MongoIPCollection)

	u, err := c.UpdateOne(ctx, bson.M{"ip": ip.IP}, bson.M{"$set": ip},
		(&options.UpdateOptions{}).SetUpsert(true))
	if u.UpsertedCount == 0 && u.ModifiedCount == 0 {
		log.Println("no change made to db", err)
	}
	return err
}
