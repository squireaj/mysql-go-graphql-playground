package main

import (
        _ "github.com/go-sql-driver/mysql"
				"github.com/graphql-go/graphql"
				
        "database/sql"
				"log"
				"encoding/json"
				"fmt"
				"math/rand"
				"net/http"
				"time"
			
)

type Document struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name,omitempty"`
	File  string  `json:"file,omitempty"`
	// Zones []Zone `json:"zones,omitempty"`
}

var documents = []Document{
	{
		ID:    1,
		Name:  "Document one",
		File:  "a23hkjhl03209n2lh34sd009f92h3h4120098fwejk13h342h...",
		// Zones: ["user": "Steve Price", "posX": 12.32, "posY": 34.23],
	},
	{
		ID:    2,
		Name:  "Document 2",
		File:  "a23hkjhl03209n2lh34sd009f92h3h4120098fwejk13h342h...",
		// Zones: ["user": "Dave Hall", "posX": 62.33, "posY": 84.27],
	},
	{
		ID:    3,
		Name:  "Document 3",
		File:  "a23hkjhl03209n2lh34sd009f92h3h4120098fwejk13h342h...",
		// Zones: ["user": "Shawn Chambless", "posX": 52.72, "posY": 64.20],
	},
}

var documentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Document",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"file": &graphql.Field{
				Type: graphql.String,
			},
			// "zones": &graphql.Field{
			// 	Type: Zone,
			// },
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single document by id
			   http://localhost:8080/document?query={document(id:1){name,file}}
			*/
			"document": &graphql.Field{
				Type:        documentType,
				Description: "Get document by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						// Find document
						for _, document := range documents {
							if int(document.ID) == id {
								return document, nil
							}
						}
					}
					return nil, nil
				},
			},
			/* Get (read) document list
			   http://localhost:8080/document?query={list{id,name,file,zone}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(documentType),
				Description: "Get document list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return documents, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {
        sdb, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/test_db")
        if err != nil {
                log.Fatal(err)
        }

        id := 1
        var name string

        if err := sdb.QueryRow("SELECT name FROM test_tb WHERE id = ? LIMIT 1", id).Scan(&name); err != nil {
                log.Fatal(err)
        }

				fmt.Println(id, name)

				http.HandleFunc("/newEmptyDocument", func(w http.ResponseWriter, r *http.Request) {
					result := executeQuery(r.URL.Query().Get("query"), schema)
					json.NewEncoder(w).Encode(result)
				})
			
				fmt.Println("Server is running on port 8080")
				http.ListenAndServe(":8080", nil)
}