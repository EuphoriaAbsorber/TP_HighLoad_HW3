package main

import (
	_ "dbproject/docs"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	deliv "dbproject/delivery"
	rep "dbproject/repository"
	usecase "dbproject/usecase"

	conf "dbproject/config"

	"github.com/jackc/pgx"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lib/pq"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, r.Method)

		for header := range conf.Headers {
			w.Header().Set(header, conf.Headers[header])
		}
		next.ServeHTTP(w, r)
		hitsCounter.Inc()
	})
}

var hitsCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "hits_counter",
	Help: "Number of hits to the server",
})

func main() {
	myRouter := mux.NewRouter()
	prometheus.MustRegister(hitsCounter)
	conn, err := pgx.ParseConnectionString("host=localhost user=art password=12345 dbname=dbproject_base sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 1000,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Println("could not connect to database: ", err)
	} else {
		log.Println("database is reachable")
	}
	defer db.Close()

	store := rep.NewStore(db)

	usecase := usecase.NewUsecase(store)

	handler := deliv.NewHandler(usecase)
	//user
	myRouter.HandleFunc(conf.PathCreateUser, handler.CreateUser).Methods(http.MethodPost)
	myRouter.HandleFunc(conf.PathProfile, handler.GetProfile).Methods(http.MethodGet)
	myRouter.HandleFunc(conf.PathProfile, handler.PostProfile).Methods(http.MethodPost)

	//forum
	myRouter.HandleFunc(conf.PathCreateForum, handler.CreateForum).Methods(http.MethodPost)
	myRouter.HandleFunc(conf.PathForumInfo, handler.GetForumInfo).Methods(http.MethodGet)
	myRouter.HandleFunc(conf.PathCreateThread, handler.CreateThread).Methods(http.MethodPost)
	myRouter.HandleFunc(conf.PathGetForumUsers, handler.GetForumUsers).Methods(http.MethodGet)
	myRouter.HandleFunc(conf.PathGetForumThreads, handler.GetForumThreads).Methods(http.MethodGet)

	//post
	myRouter.HandleFunc(conf.PathPost, handler.GetPostById).Methods(http.MethodGet)
	myRouter.HandleFunc(conf.PathPost, handler.UpdatePost).Methods(http.MethodPost)

	//service
	myRouter.HandleFunc(conf.PathGetServiceStatus, handler.ServiceStatus).Methods(http.MethodGet)
	myRouter.HandleFunc(conf.PathServiceClear, handler.ServiceClear).Methods(http.MethodPost)

	//threads
	myRouter.HandleFunc(conf.PathCreatePosts, handler.CreatePosts).Methods(http.MethodPost)
	myRouter.HandleFunc(conf.PathThreadInfo, handler.GetThreadInfo).Methods(http.MethodGet)
	myRouter.HandleFunc(conf.PathThreadInfo, handler.UpdateThreadInfo).Methods(http.MethodPost)
	myRouter.HandleFunc(conf.PathThreadVote, handler.VoteForThread).Methods(http.MethodPost)
	myRouter.HandleFunc(conf.PathGetThreadPosts, handler.GetThreadPosts).Methods(http.MethodGet)

	myRouter.PathPrefix(conf.PathDocs).Handler(httpSwagger.WrapHandler)

	myRouter.Use(loggingMiddleware)
	myRouter.Path("/metrics").Handler(promhttp.Handler())

	err = http.ListenAndServe(conf.Port, myRouter)
	if err != nil {
		log.Println("cant serve", err)
	}
}
