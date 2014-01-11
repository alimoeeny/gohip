package gohip

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	auth_token = "xNGZ9zP8mwTVop90sDxJdePIsRniuQb7j9PuDGuU"
	base_url   = "https://api.hipchat.com/v2/"
)

type Hip struct {
	auth_token string
}

type Room struct {
	id   string
	link string
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(2*time.Second))
}

func instanceId() string {
	transport := http.Transport{Dial: dialTimeout}
	client := http.Client{
		Transport: &transport,
	}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return "unknown"
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "unknown"
		} else {
			return string(body)
		}
	}
}

func (hip *Hip) GetAllRooms() (map[string]Room, error) {
	rooms := make(map[string]Room)
	transport := http.Transport{Dial: dialTimeout}
	client := http.Client{
		Transport: &transport,
	}
	resp, err := client.Get(base_url + "room" + "?auth_token=" + auth_token)
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
		room := Room{id: strconv.FormatFloat(vv["id"].(float64), 'f', 0, 64), link: vv["links"].(map[string]interface{})["self"].(string)}
		rooms[room.id] = room
	}
	return rooms, err
}
