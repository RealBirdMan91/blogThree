package main

import (
	"blogThree/internal/interfaces/authctx"
	"blogThree/internal/interfaces/graph"
	"blogThree/internal/interfaces/graph/resolvers"
	"blogThree/internal/interfaces/httpctx"
	db "blogThree/internal/platform/postgres"
	"blogThree/internal/platform/postgres/migrations"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/crypto/bcrypt"

	userApp "blogThree/internal/user/app"
	userPolicies "blogThree/internal/user/app/policies"
	userPostgres "blogThree/internal/user/infra/postgres"
	userSecurity "blogThree/internal/user/infra/security"

	authApp "blogThree/internal/auth/app"
	jwt "blogThree/internal/auth/infra/jwt"
	authPostgres "blogThree/internal/auth/infra/postgres"
)

func mustGetEnv() (string, string) {
	_ = godotenv.Load() // optional, nur lokal wichtig; Fehler ignorieren

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return port, dsn
}
func main() {
	port, databaseUrl := mustGetEnv()

	pgDB, err := db.Open(databaseUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pgDB.Close()
	err = db.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	//USER
	userRepo := userPostgres.NewPostgresUserRepo(pgDB)
	policy := userPolicies.NewSimplePasswordPolicy(8, true)
	hasher := userSecurity.NewBcryptHasher(bcrypt.DefaultCost)
	userService := userApp.NewService(userRepo, policy, hasher)

	//AUTH
	authRepo := authPostgres.NewRefreshTokenRepo(pgDB)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	encoder := jwt.New([]byte(jwtSecret))
	authService := authApp.NewService(authRepo, encoder)

	//ROOT RESOLVER INITIALIZATION
	res := &resolvers.Resolver{UserSvc: userService, AuthSvc: authService}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: res}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query",
		httpctx.Inject( // Request+Response in den Context legen (für Cookies etc.)
			authctx.Middleware(encoder)( // Access-Token aus Authorization prüfen, UserID in Context
				srv,
			),
		),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
