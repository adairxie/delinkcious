package service

import (
	"log"
	"net/http"
	"os"

	"github.com/adairxie/delinkcious/pkg/db_util"

	sgm "github.com/adairxie/delinkcious/pkg/user_manager"
	httptransport "github.com/go-kit/kit/transport/http"
)

func Run() {
	dbHost, dbPort, err := db_util.GetDbEndpoint("user")
	if err != nil {
		log.Fatal(err)
	}

	store, err := sgm.NewDbUserStore(dbHost, dbPort, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}

	svc, err := sgm.NewUserManager(store)
	if err != nil {
		log.Fatal(err)
	}

	registerHandler := httptransport.NewServer(
		makeRegisterEndpoint(svc),
		decodeRegisterRequest,
		encodeResponse,
	)

	LoginHandler := httptransport.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest,
		encodeResponse,
	)

	LogoutHandler := httptransport.NewServer(
		makeLogoutEndpoint(svc),
		decodeLogoutRequest,
		encodeResponse,
	)

	http.Handle("/register", registerHandler)
	http.Handle("/login", LoginHandler)
	http.Handle("/logout", LogoutHandler)

	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
