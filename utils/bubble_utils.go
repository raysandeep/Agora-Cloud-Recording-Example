package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/raysandeep/Estimator-App/schemas"
	"github.com/spf13/viper"
)

func PatchRecord(id string, payload map[string]interface{}) {

	url := fmt.Sprintf("https://terrierapp.com/version-test/api/1.1/obj/job/%s?api_key=%s", id, viper.GetString("BUBBLE_API_KEY"))
	method := "PATCH"

	requestBody, _ := json.Marshal(&payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(requestBody)))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func FetchRecord(shortId string) (*schemas.BubbleRecord, error) {

	constraints := []map[string]interface{}{
		{
			"key":             "short_id",
			"constraint_type": "equals",
			"value":           shortId,
		},
	}

	constraintsString, err := json.Marshal(constraints)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("api_key", "this will be esc@ped!")
	params.Add("constraints", string(constraintsString))

	url := fmt.Sprintf("https://terrierapp.com/version-test/api/1.1/obj/job?%s", params.Encode())
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result schemas.BubbleResponse
	err = json.Unmarshal(body,&result)
	if err != nil {
		return nil, err
	}
	return &result.Response.Results[0], nil
}
