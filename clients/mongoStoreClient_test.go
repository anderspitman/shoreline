package clients

import (
	"encoding/json"
	mongo "github.com/tidepool-org/go-common/clients/mongo"
	api "github.com/tidepool-org/shoreline/api"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"testing"
)

func TestMongoStoreUserOperations(t *testing.T) {

	type Config struct {
		Mongo *mongo.Config `json:"mongo"`
	}

	var config Config

	if jsonConfig, err := ioutil.ReadFile("../config/server.json"); err == nil {

		if err := json.Unmarshal(jsonConfig, &config); err != nil {
			log.Fatal(err)
		}

		log.Println("config is", config.Mongo)

		mc := NewMongoStoreClient(config.Mongo)

		/*
		 * INIT THE TEST - we use a clean copy of the collection before we start
		 */

		if err := mc.usersC.DropCollection(); err != nil {
			t.Fatalf("We couldn't drop the users collection and start the tests fresh ", err)
		}

		if err := mc.usersC.Create(&mgo.CollectionInfo{}); err != nil {
			t.Fatalf("We couldn't created the users collection for these tests ", err)
		}

		/*
		 * THE TESTS
		 */
		user, _ := api.NewUser("test user", "myT35t", []string{"test@foo.bar"})

		if err := mc.UpsertUser(user); err != nil {
			t.Fatalf("we could not create the user %v", err)
		}

		user.Name = "test user updated"

		if err := mc.UpsertUser(user); err != nil {
			t.Fatalf("we could not update the user %v", err)
		}

		toFindByName := &api.User{Name: user.Name}

		if found, err := mc.GetUser(toFindByName); err != nil {
			t.Fatalf("we could not find the the user bu name %v", toFindByName)
		} else {
			if found.Name != toFindByName.Name {
				t.Fatalf("the user we found doesn't match what we asked for %v", found)
			}
		}

		toFindById := &api.User{Id: user.Id}

		if found, err := mc.GetUser(toFindById); err != nil {
			t.Fatalf("we could not find the the user by id %v", toFindById)
		} else {
			if found.Id != toFindById.Id {
				t.Fatalf("the user we found doesn't match what we asked for %v", found)
			}
		}

		toFindByEmails := &api.User{Emails: user.Emails}

		if found, err := mc.GetUser(toFindByEmails); err != nil {
			t.Fatalf("we could not find the the user by id %v", toFindByEmails)
		} else {
			if found.Emails[0] != toFindByEmails.Emails[0] {
				t.Fatalf("the user we found doesn't match what we asked for %v", found)
			}
		}

	} else {
		t.Fatalf("wtf - failed parsing the config %v", err)
	}
}
