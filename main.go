package main

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"ocb.amot.io/docs"
	"ocb.amot.io/internal/adapters/clients"
	"ocb.amot.io/internal/adapters/repositories"
	"ocb.amot.io/internal/adapters/services"
	"ocb.amot.io/internal/core/domain"
	"ocb.amot.io/internal/presentation/rest"
)

// @securityDefinitions.apikey ApiKeyAuth
func main() {
	// required environment variables
	appHost := "localhost"
	appPort := "9200"

	// setup swagger config
	docs.SwaggerInfo.Title = "Auth Microservice"
	docs.SwaggerInfo.Description = "A service for issuing jwt tokens"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s",appHost,appPort)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Version = "1.0"

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get secret and retrieve data
	mysqlCon,err := client.CoreV1().Secrets("default").Get(context.TODO(),"authms-mysql-conn",metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	dbName := mysqlCon.Data["db"]
	username := mysqlCon.Data["username"]
	password := mysqlCon.Data["password"]

	authMsEnv,err := client.CoreV1().ConfigMaps("default").Get(context.TODO(),"authms-env",metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	dbHost := authMsEnv.Data["host"]
	dbPort := authMsEnv.Data["port"]

	// setup dependencies
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",username,password,dbHost,dbPort,dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err!=nil{
		panic(err)
	}
    err = db.AutoMigrate(domain.RefreshToken{})
	if err != nil {
		panic(err)
	}
	tr := repositories.NewTokenRepository(db)
	us := clients.NewUserClient()
	k := services.NewKeyGenerator("keys/tym.pem")

	// an entry point to the core layer
	ts := services.NewTokenService(tr,us,k)

	// mux implementation of rest api
	service := presentation.NewMuxHttpServer(ts)
	// start service
	err = service.Start(appHost,appPort)
	if err != nil {
		log.Fatal(err)
	}
}