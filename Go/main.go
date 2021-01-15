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
					sdb, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/test_db")
					if err != nil {
						log.Fatal(err)
					}
					id, ok := p.Args["id"].(int)
					if ok {
						var name, file string
						doc := Document{}
						if err := sdb.QueryRow("SELECT * FROM test_tb WHERE id = ? LIMIT 1", id).Scan(&id, &name, &file); err != nil {
							log.Fatal(err)
						}
						doc.ID = 	int64(id)
						doc.Name = name
						doc.File = file
						return doc, nil
					}
					return nil, nil
				},
			},
			/* Get (read) document list
			   http://localhost:8080/document?query={list{id,name,file,zone}}
			*/
			// "list": &graphql.Field{
			// 	Type:        graphql.NewList(documentType),
			// 	Description: "Get document list",
			// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// 		// return SELECT * 
			// 		// FROM Reflow  
			// 		// WHERE ReflowProcessID = somenumber
			// 		// ORDER BY ID DESC
			// 		// LIMIT 20
			// 		return nil, nil
			// 	},
			// },
			"newDocument": &graphql.Field{
				Type:        documentType,
				Description: "Get New Document with UUID included",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return "2344312323", nil
				},
			},
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new document item
		http://localhost:8080/document?query=mutation+_{create(name:"Inca Kola",file:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",zone:1.99){id,name,file,zone}}
		*/
		"create": &graphql.Field{
			Type:        documentType,
			Description: "Create new document",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"file": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				// "zones": &graphql.ArgumentConfig{
				// 	Type: graphql.NewNonNull(graphql.Float),
				// },
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())
				document := Document{
					ID:    int64(rand.Intn(100000)), // generate random ID
					Name:  params.Args["name"].(string),
					File:  params.Args["file"].(string),
					// Zones: params.Args["zones"].(Zone),
				}
				// documents = append(documents, document)
				return document, nil
			},
		},

		/* Update document by id
		   http://localhost:8080/document?query=mutation+_{update(id:1,zones:3.95){id,name,file,zone}}
		*/
		"update": &graphql.Field{
			Type:        documentType,
			Description: "Update document by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"file": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				// "zones": &graphql.ArgumentConfig{
				// 	Type: Zone[],
				// },
			},
			// Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// 	id, _ := params.Args["id"].(int)
			// 	name, nameOk := params.Args["name"].(string)
			// 	file, fileOk := params.Args["file"].(string)
			// 	// zone, zoneOk := params.Args["zones"].(Zone)
			// 	document := Document{}
			// 	return document, nil
			// },
		},

		/* Delete document by id
		   http://localhost:8080/document?query=mutation+_{delete(id:1){id,name,file,zones}}
		*/
		"delete": &graphql.Field{
			Type:        documentType,
			Description: "Delete document by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			// Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// 	id, _ := params.Args["id"].(int)
			// 	document := Document{}
			// 	return document, nil
			// },
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
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
	var name, file string
	doc := Document{}
	if err := sdb.QueryRow("SELECT * FROM test_tb WHERE id = ? LIMIT 1", id).Scan(&id, &name, &file); err != nil {
					log.Fatal(err)
	}
	doc.ID = 	int64(id)
	doc.Name = name
	doc.File = file
	fmt.Println(doc)

	http.HandleFunc("/document", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	http.HandleFunc("/newEmptyDocument", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}