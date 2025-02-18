package webapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var api string = "https://sensors.bieda.it"

func LastData(id string) (Measurement, error) {
	req, err := http.Get(fmt.Sprintf("%s/latest/%s", api, id))
	if err != nil {
		return Measurement{}, err
	}
	if req.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Status Code: %v", req.StatusCode))
		log.Println(err)
		return Measurement{}, err
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return Measurement{}, err
	}
	var m Measurement
	err = json.Unmarshal(b, &m)
	if err != nil {
		return Measurement{}, err
	}
	return m, nil
}

func Data(skip int) ([]Measurement, error) {
	req, err := http.Get(fmt.Sprintf("%s/data/%v", api, skip))
	if err != nil {
		return nil, err
	}
	if req.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Status Code: %v", req.StatusCode))
		log.Println(err)
		return nil, err
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var m []Measurement
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
