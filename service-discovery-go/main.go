package main

import (
	"fmt"
	"log"
	"net/http"
	dockerclient "service-discovery/docker-client"
	"service-discovery/registrar"
	serviceregistry "service-discovery/service-registry"
	"sync/atomic"
)

type Application struct {
	RequestCount uint64
	SRegistry    *serviceregistry.ServiceRegistry
}

func (a *Application) Handle(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&a.RequestCount, 1)

	if a.SRegistry.Len() == 0 {
		w.Write([]byte(`No backend entry in the service registry`))
		return
	}

	backendIndex := int(atomic.LoadUint64(&a.RequestCount) % uint64(a.SRegistry.Len()))
	fmt.Printf("Request routing to instance %d\n", backendIndex)

	backend := a.SRegistry.GetByIndex(backendIndex)
	backend.Proxy().ServeHTTP(w, r)
}

func main() {
	registry := serviceregistry.ServiceRegistry{}
	registry.Init()

	dockerClient, err := dockerclient.NewDockerClient()
	if err != nil {
		panic(err)
	}

	registrar := registrar.Registrar{SRegistry: &registry, DockerClient: dockerClient}

	if err = registrar.Init(); err != nil {
		panic(err)
	}
	go registrar.Observe()

	app := Application{SRegistry: &registry}
	http.HandleFunc("/reverse-proxy", app.Handle)

	log.Fatalln(http.ListenAndServe(":8000", nil))
}
