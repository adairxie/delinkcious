package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/adairxie/delinkcious/pkg/db_util"
	om "github.com/adairxie/delinkcious/pkg/object_model"
	"github.com/gorilla/mux"

	lm "github.com/adairxie/delinkcious/pkg/link_manager"
	nats "github.com/adairxie/delinkcious/pkg/link_manager_events"
	sgm "github.com/adairxie/delinkcious/pkg/social_graph_client"
	httptransport "github.com/go-kit/kit/transport/http"
)

type EventSink struct {
}

func (s *EventSink) OnLinkAdded(username string, link *om.Link) {
	//log.Println("Link added")
}

func (s *EventSink) OnLinkUpdated(username string, link *om.Link) {
	//log.Println("Link updated")
}

func (s *EventSink) OnLinkDeleted(username string, url string) {
	//log.Println("Link deleted")
}

func Run() {
	dbHost, dbPort, err := db_util.GetDbEndpoint("social_graph")
	if err != nil {
		log.Fatal(err)
	}

	store, err := lm.NewDbLinkStore(dbHost, dbPort, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	sgHost := os.Getenv("SOCIAL_GRAPH_SERVICE_HOST")
	if sgHost == "" {
		sgHost = "localhost"
	}

	sgPort := os.Getenv("SOCIAL_GRAPH_SERVICE_PORT")
	if sgPort == "" {
		sgPort = "9090"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	maxLinksPerUserStr := os.Getenv("MAX_LINKS_PER_USER")
	if maxLinksPerUserStr == "" {
		maxLinksPerUserStr = "10"
	}

	maxLinksPerUser, err := strconv.ParseInt(maxLinksPerUserStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	socialGraphClient, err := sgm.NewClient(fmt.Sprintf("%s:%s", sgHost, sgPort))
	if err != nil {
		log.Fatal(err)
	}

	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")

	var eventSink om.LinkManagerEvents
	if natsHostname != "" {
		natsUrl := natsHostname + ":" + natsPort
		eventSink, err = nats.NewEventSender(natsUrl)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		eventSink = &EventSink{}
	}

	svc, err := lm.NewLinkManager(store, socialGraphClient, eventSink, maxLinksPerUser)
	if err != nil {
		log.Fatal(err)
	}

	getLinksHandler := httptransport.NewServer(
		makeGetLinksEndpoint(svc),
		decodeGetLinksRequest,
		encodeResponse,
	)

	addLinkHandler := httptransport.NewServer(
		makeAddLinkEndpoint(svc),
		decodeAddLinkRequest,
		encodeResponse,
	)

	updateLinkHandler := httptransport.NewServer(
		makeUpdateLinkEndpoint(svc),
		decodeUpdateLinkRequest,
		encodeResponse,
	)

	deleteLinkHandler := httptransport.NewServer(
		makeDeleteLinkEndpoint(svc),
		decodeDeleteLinkRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("GET").Path("/links").Handler(getLinksHandler)
	r.Methods("POST").Path("/links").Handler(addLinkHandler)
	r.Methods("PUT").Path("/links").Handler(updateLinkHandler)
	r.Methods("DELETE").Path("/links").Handler(deleteLinkHandler)

	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
