// This is a client module to support server-side use of the Tidepool
// service called user-api.
package oauth2

import (
	"errors"
	"log"
)

type (
	OAuth2Client struct{ oauth2api *Api }
	//Generic client interface that we will implement
	Data   map[string]interface{}
	Client interface {
		CheckToken(token string) (Data, error)
	}
)

func NewOAuth2Client(api *Api) *OAuth2Client {
	return &OAuth2Client{oauth2api: api}
}

func (client *OAuth2Client) CheckToken(token string) (Data, error) {

	if data, err := client.oauth2api.storage.LoadAccess(token); err != nil {
		log.Print(OAUTH2_API_PREFIX, "OAuth2Client CheckToken ", err.Error())
		return nil, err
	} else if data == nil {
		log.Print(OAUTH2_API_PREFIX, "OAuth2Client CheckToken ", "nothing found")
		return nil, errors.New("nothing found")
	} else {
		returnData := Data{}
		returnData["userId"] = data.Client.GetId()
		ud := data.Client.GetUserData().(map[string]interface{})

		//so that we have the details of the user that has given authorization via the oauth2 authorize process
		if ud["AppUser"] != nil {
			log.Print(OAUTH2_API_PREFIX, "OAuth2Client CheckToken setting the authUserId")
			returnData["authUserId"] = ud["AppUser"]
		} else {
			log.Print(OAUTH2_API_PREFIX, "OAuth2Client CheckToken no authUserId found")
		}

		log.Print(OAUTH2_API_PREFIX, "OAuth2Client CheckToken ", returnData)
		return returnData, nil
	}
}
