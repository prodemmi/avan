package network

import (
	"avan/internal/config"
	"fmt"
	"github.com/hashicorp/memberlist"
	"go.uber.org/zap"
	"io"
	"time"
)

type (
	// The GossipService type represents a single node in the cluster, it contains
	// the list of other members in the cluster and an HTTP client for
	// directly messaging other nodes.
	GossipService struct {
		memberlist *memberlist.Memberlist
	}
)

func (g GossipService) PrintMembers() {
	for {
		for _, node := range g.memberlist.Members() {
			fmt.Printf("Name: %s, Host: %s Port: %d", node.Name, node.Addr, node.Port)
		}
		time.Sleep(5 * time.Second)
	}
}

func (g GossipService) GetNodeHealth() (int, int, int) {
	nodes := g.memberlist.Members()
	healthy := 0
	unhealthy := 0
	for _, node := range nodes {
		if node.State == memberlist.StateAlive {
			healthy++
		} else {
			unhealthy++
		}
	}
	return len(nodes), healthy, unhealthy
}

func NewGossipService(logger zap.Logger, appConfig config.Config) GossipService {
	cfg := memberlist.DefaultLocalConfig()
	cfg.LogOutput = io.Discard
	cfg.BindPort = 7947

	// Create a new memberlist instance
	memberlistInstance, err := memberlist.Create(cfg)
	if err != nil {
		logger.Fatal("failed to create gossip memberlist", zap.Error(err))
	}

	// Join the cluster (you can specify other node addresses here)
	_, err = memberlistInstance.Join(appConfig.Hosts) // Replace with actual node IPs/addresses
	if err != nil {
		logger.Fatal("failed to join gossip memberlist", zap.Error(err))
	}

	return GossipService{
		memberlist: memberlistInstance,
	}
}
