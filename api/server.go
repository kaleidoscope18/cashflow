package api

import (
	"cashflow/api/graph"
	"cashflow/api/graph/generated"
	"cashflow/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func Run(app *models.App) {
	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{App: app},
	})

	server := handler.NewDefaultServer(schema)

	browserTabTitle := "Cashflow GraphQL Playground"
	http.Handle("/playground", playground.Handler(browserTabTitle, "/graphql"))
	http.Handle("/graphql", server)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "GraphQL playground:\thttp://localhost:%s/playground\t\n", port)
	fmt.Fprintf(w, "GraphQL endpoint:\thttp://localhost:%s/graphql\t\n", port)
	w.Flush()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
