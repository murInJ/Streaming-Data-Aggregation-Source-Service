package utils

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
)

func PostJson(url string, jsonData []byte) error {
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("http status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
