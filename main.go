package main

import (
	"log"

	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	conf := api.DefaultConfig()

	// Create client
	client, err := api.NewClient(conf)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	consul_agent := client.Agent()

	// 测试崩溃节点
	check_unpass := &api.AgentServiceCheck{
		HTTP:     "http://192.168.10.10:8080",
		Interval: "10s",
		Timeout:  "1s",
	}

	service2 := &api.AgentServiceRegistration{
		ID:      "redis8081",
		Name:    "redis_svc",
		Tags:    []string{"123", "abc"},
		Port:    8081,
		Address: "192.168.10.10",
	}
	service1 := &api.AgentServiceRegistration{
		ID:      "redis8080",
		Name:    "redis_svc",
		Tags:    []string{"123", "abc"},
		Port:    8080,
		Address: "192.168.10.10",
	}
	err = consul_agent.ServiceRegister(service1)
	if err != nil {
		log.Fatalf("reg service err: %s", err)
	}
	err = consul_agent.ServiceRegister(service2)
	if err != nil {
		log.Fatalf("reg service err: %s", err)
	}

	redis_services, err := consul_agent.Services()
	if err != nil {
		log.Fatalf("get service err: %s", err)
	}

	err = consul_agent.CheckRegister(&api.AgentCheckRegistration{
		Name:              "node1_check",
		AgentServiceCheck: *check_unpass,
	})
	if err != nil {
		log.Fatalf("check register err: %s", err)
	}
	for name, redis_service := range redis_services {
		log.Println("name: ", name)
		spew.Dump(redis_service)
	}
	router := gin.Default()
	router.GET("/alive", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})
	router.Run(":8081")
}
