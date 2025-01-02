package api

import (
	"cashflow/api/graph"
	"cashflow/api/graph/generated"
	"cashflow/api/middlewares"
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

	config := generated.Config{
		Resolvers:  &graph.Resolver{App: app},
		Directives: graph.GetDirectives(app),
	}
	schema := generated.NewExecutableSchema(config)

	server := handler.NewDefaultServer(schema)

	http.Handle("/playground", playground.Handler("Cashflow GraphQL Playground", "/graphql"))
	http.Handle("/graphql", middlewares.AuthMiddleware(server))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			port = defaultPort
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "GraphQL playground:\thttp://localhost:%s/playground\t\n", port)
	fmt.Fprintf(w, "GraphQL endpoint:\thttp://localhost:%s/graphql\t\n", port)
	w.Flush()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
