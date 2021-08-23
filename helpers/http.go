package helpers

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Get(url string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("error while connecting to api")
	} else {
		if res.Status == "200 OK" {
			defer res.Body.Close()
			reqBody, _ := ioutil.ReadAll(res.Body)
			return reqBody
		} else {
			return []byte{}
		}
	}
	return []byte{}
}
