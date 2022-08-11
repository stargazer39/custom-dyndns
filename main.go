package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	for {
		if err := start(); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second * 5)
	}
}

func start() error {
	godotenv.Load()
	// Get port environment
	port := os.Getenv("PORT")
	mongo_uri := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_uri))

	if err != nil {
		return err
	}

	ddns := client.Database("ddns").Collection("ddns")

	r := gin.Default()

	r.GET("/nic/update", func(c *gin.Context) {
		hostname, ok := c.GetQuery("hostname")

		if !ok {
			c.Status(404)
			return
		}

		myip, ok := c.GetQuery("myip")

		if !ok {
			c.Status(404)
			return
		}

		var record DDNSEntry

		err := ddns.FindOne(c, bson.D{{Key: "hostname", Value: hostname}}).Decode(&record)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				record.ID = primitive.NewObjectID()
				record.Added = time.Now()
				record.Hostname = hostname
				record.IP = myip
				record.Updated = record.Added

				if _, err := ddns.InsertOne(c, record); err != nil {
					c.Status(400)
					return
				}

				c.Status(200)
				return
			}
		}

		update := bson.D{{
			Key: "$set", Value: bson.D{{Key: "updated", Value: time.Now()}, {Key: "ip", Value: myip}},
		}}

		if res, err := ddns.UpdateByID(c, record.ID, update); err != nil || res.MatchedCount != 1 {
			log.Println(err)
			c.Status(500)
			return
		}

		c.Status(200)
	})

	return r.Run(":" + port)
}
