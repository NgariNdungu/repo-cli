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
// var endpoint = "https://repo-cli.free.beeceptor.com"

// TODO: validate order flag
var ordFields = map[string]string {
	"creation" : "CREATED_AT",
	"update" : "UPDATED_AT",
	"push" : "PUSHED_AT",
	"name" : "NAME",
}

func List(t string, orderBy string, desc bool, limit int) []string {
	sort := "ASC"
	if desc {
		sort = "DESC"
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("{\n\"query\":")
	queryBuilder.WriteString("\"query { viewer { repositories(")
	repoParams := fmt.Sprintf("first:%d,orderBy:{field:%s,direction:%s}",
	limit, ordFields[orderBy], sort)
	queryBuilder.WriteString(repoParams)
	queryBuilder.WriteString(") { edges { node { name url id } } } } }\"\n}")

	query := queryBuilder.String()
	// queryO := `
	// {
	// 	"query": "query { viewer { repositories(first:10,orderBy:{field:CREATED_AT,direction:DESC}) { edges { node { name url id } } } } }"
	// }
	// `
	// fmt.Print(queryO)

	// generated with https://mholt.github.io/json-to-go/
	type Node struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		URL  string `json:"url"`
	}
	type Edges struct {
		Node Node `json:"node"`
	}
	type Repositories struct {
		Edges []Edges `json:"edges"`
	}
	type Viewer struct {
		Repositories Repositories `json:"repositories"`
	}
	type Data struct {
		Viewer Viewer `json:"viewer"`
	}
	type queryResponse struct {
		Data Data `json:"data"`
	}

	var out queryResponse
	req,err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(query))
	req.Header.Add("authorization", "bearer " + t)
	resp,err := client.Do(req)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
	// second argument should be a pointer
	err = json.Unmarshal(body, &out)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
	// m := out.(map[string]interface{})
	for k,v := range out.Data.Viewer.Repositories.Edges{
		fmt.Printf("%d: %s | %s | %s\n", k,v.Node.Name, v.Node.ID, v.Node.URL)
	}
	return nil
}

