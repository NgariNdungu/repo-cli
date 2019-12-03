// Package repo provides functions for working with github repos
package	repo

import	(
	"os"
	"net/http"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"encoding/json"
)

var client = &http.Client{}

var endpoint = "https://api.github.com/graphql"

func List(t string) []string {
	query := `
	{
		"query": "query {
			viewer {
				repositories(first:10) {
					edges {
						node {
							name
							url
							id
						}
					}
				}
			}
		}"
	}
	`
	type data struct {
		Viewer map[string]interface{}
	}
	type viewer struct {
		Repositories map[string]interface{}
	}
	type node struct {
		Name string
		Url string
		Id string
	}
	type repositories struct {
		Edges []node
	}


	var out interface{}
	req,err := http.NewRequest(http.MethodGet, endpoint, strings.NewReader(query))
	req.Header.Add("authorization", "bearer " + t)
	resp,err := client.Do(req)
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
	err = json.Unmarshal(body, out)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
	m := out.(map[string]interface{})
	for k,v := range m {
		fmt.Printf("%s is %v", k,v)
	}
	return nil
}

func Schema(t string) error {
	req,err := http.NewRequest("GET", endpoint, nil)
	fmt.Printf("Token: %s", t)
	req.Header.Add("authorization", "bearer " + t)
	resp,err := client.Do(req)
	if err != nil {
		return err
	}
	return resp.Write(os.Stdout)
}

