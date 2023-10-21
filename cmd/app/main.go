package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"imersaofullcycle/internal/app/akafka"
	"imersaofullcycle/internal/app/repository"
	"imersaofullcycle/internal/app/usecase"
	"imersaofullcycle/internal/app/web"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306/products)")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := repository.NewProductRepositoryMySql(db)

	createProductUseCase := usecase.NewCreateProductUseCase(repository)
	listProductsUseCase := usecase.NewListProductsUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUseCase, listProductsUseCase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductsHandler)

	go http.ListenAndServe(":8000", r)

	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"products"}, "host.docker.internal:9094", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			fmt.Println("error in json: ", err)
			continue
		}

		_, err = createProductUseCase.Execute(dto)
		if err != nil {
			fmt.Println("error during createproduct: ", err)
			continue
		}
	}
}
