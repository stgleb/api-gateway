package registration

import (
	. "api-gateway"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func NewRegistrar(config *TomlConfig, logger log.Logger) (registar sd.Registrar) {
	rand.Seed(time.Now().UnixNano())

	var client consul.Client
	{
		consulConfig := api.DefaultConfig()
		consulConfig.Address = fmt.Sprintf("%s:%d",
			config.ServiceDiscovery.ConsulAddress,
			config.ServiceDiscovery.ConsulPort)
		consulClient, err := api.NewClient(consulConfig)

		if err != nil {
			logger.Log("message", "Can not find consul to do service discovery", "err", err)
		}

		client = consul.NewClient(consulClient)
	}

	check := api.AgentServiceCheck{
		HTTP: "http://" +
			fmt.Sprintf("%s:%d", config.ServiceDiscovery.AdvertisedAddress,
				config.ServiceDiscovery.AdvertisedPort) +
			"/health",
		Interval: config.ServiceDiscovery.Interval,
		Timeout:  config.ServiceDiscovery.Timeout,
		Notes:    "Basic health checks",
	}

	hostName, _ := os.Hostname()

	num := rand.Intn(100) // to make service ID unique
	asr := api.AgentServiceRegistration{
		ID:      config.Main.ServiceName + strconv.Itoa(num), //unique service ID
		Name:    config.Main.ServiceName,
		Address: config.ServiceDiscovery.AdvertisedAddress,
		Port:    config.ServiceDiscovery.AdvertisedPort,
		Tags:    []string{config.Main.ServiceName, hostName},
		Check:   &check,
	}
	registar = consul.NewRegistrar(client, &asr, logger)
	return
}
