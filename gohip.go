package gohip

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	base_url = "https://api.hipchat.com/v2/"
)

type Hip struct {
	Auth_token string
}

type Room struct {
	Name             string
	Id               string
	Link             string
	Created          time.Time
	Guest_access_url string
	Is_archived      bool
	Last_active      time.Time
	Webhook          string
	Privacy          string
	Topic            string
	Xmpp_jid         string
	//owner            user
	//participants     []user
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(2*time.Second))
}

func (hip *Hip) GetAllRooms() (map[string]Room, error) {
	rooms := make(map[string]Room)
	transport := http.Transport{Dial: dialTimeout}
	client := http.Client{
		Transport: &transport,
	}
	resp, err := client.Get(base_url + "room" + "?auth_token=" + hip.Auth_token)
	if err != nil {
		return rooms, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var jsonresponse map[string]interface{}
	err = json.Unmarshal(body, &jsonresponse)
	if err != nil {
		return rooms, err
	}
	for _, v := range jsonresponse["items"].([]interface{}) {
		vv := v.(map[string]interface{})
		room := Room{Id: strconv.FormatFloat(vv["id"].(float64), 'f', 0, 64), Link: vv["links"].(map[string]interface{})["self"].(string)}
		rooms[room.Id] = room
	}
	return rooms, err
}

func (hip *Hip) PostToRoom(room Room, body string) error {
	transport := http.Transport{Dial: dialTimeout}
	client := http.Client{
		Transport: &transport,
	}
	payload := struct {
		Message        string `json:"message"`
		Color          string `json:"color"`
		Message_format string `json:"message_format"`
	}{body, "red", "text"}
	buf, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error converting the message to json payload:%v\n", err)
		return err
	}
	resp, err := client.Post(base_url+"room/"+room.Id+"/notification?auth_token="+hip.Auth_token, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	if resp.StatusCode > 399 {
		return errors.New(resp.Status)
	}
	return nil
}
