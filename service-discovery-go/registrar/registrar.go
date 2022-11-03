package registrar

import (
	"context"
	"fmt"
	dockerclient "service-discovery/docker-client"
	serviceregistry "service-discovery/service-registry"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type Registrar struct {
	DockerClient *dockerclient.DockerClient
	SRegistry    *serviceregistry.ServiceRegistry
}

func (r *Registrar) Init() error {
	cList, err := r.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("ancestor", dockerclient.HelloServiceImageName),
			filters.Arg("status", dockerclient.ContainerRunningState),
		),
	})
	if err != nil {
		return err
	}

	for _, c := range cList {
		r.SRegistry.Add(c.ID, findContainerAddress(c.Ports[0].PublicPort))
	}

	return nil
}

func (r *Registrar) Observe() {
	msgCh, errCh := r.DockerClient.Events(context.Background(), types.EventsOptions{
		Filters: filters.NewArgs(
			filters.Arg("type", "container"),
			filters.Arg("image", dockerclient.HelloServiceImageName),
			filters.Arg("event", "start"),
			filters.Arg("event", "kill"),
		),
	})

	for {
		select {
		case c := <-msgCh:
			fmt.Printf("State of the container %s is %s\n", c.ID, c.Status)
			if c.Status == dockerclient.ContainerKillState {
				r.SRegistry.RemoveByContainerID(c.ID)
			} else if c.Status == dockerclient.ContainerStartState {
				port, err := r.DockerClient.GetContainerPort(context.Background(), c.ID)
				if err != nil {
					fmt.Printf("err getting newly started container port %s\n", err.Error())
					continue
				}
				r.SRegistry.Add(c.ID, findContainerAddress(port))
			}
		case err := <-errCh:
			fmt.Println("Error Docker Event Chan", err.Error())
		}
	}
}

func findContainerAddress(cPort uint16) string {
	return fmt.Sprintf("http://localhost:%d", cPort)
}
