package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ottolauncher/recipes/graph"
	db "github.com/ottolauncher/recipes/graph/db/mongo"
	"github.com/ottolauncher/recipes/graph/generated"
	"github.com/ottolauncher/recipes/graph/model"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	// 	TokenLookup: "header:X-XSRF-TOKEN",
	// }))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	var (
		once sync.Once
		dao  *mongo.Client
	)

	once.Do(func() {
		dao = db.Init()
	})

	src := dao.Database("recipedb")

	defer func() {
		if err := dao.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	rm := db.NewRecipeManager(src)
	im := db.NewIngredientManager(src)

	config := generated.Config{Resolvers: &graph.Resolver{RM: rm, IM: im, Recipes: []*model.Recipe{}, RecipeObservers: map[string]chan []*model.Recipe{}}}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	e.GET("/playground", func(c echo.Context) error {
		playground.Handler("GraphQL playground", "/query").ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.POST("/query", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	h2s := &http2.Server{
		MaxConcurrentStreams: 250,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          10 * time.Second,
	}

	s := http.Server{
		Addr:    ":" + port,
		Handler: h2c.NewHandler(e, h2s),
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
