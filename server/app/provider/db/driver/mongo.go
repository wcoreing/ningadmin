package driver

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"wnadmin/app/provider/db/model"
)

// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#pkg-examples
/*
collection := global.DB.GameLog.Collection("GameLog202403")
	result := bson.M{}
	err := collection.FindOne(context.Background(), bson.M{"LogID": "1709256625010210072010001"}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("cur", result)
*/
func MongoByConf(conf model.Mongo) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Uri()))
	if err != nil {
		panic(err)
	}

	return client
}
