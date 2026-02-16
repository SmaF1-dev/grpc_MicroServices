package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	order "github.com/SmaF1-dev/grpc_MicroServices/api/order"
	"github.com/SmaF1-dev/grpc_MicroServices/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := order.NewOrderServiceClient(conn)

	const numOrders = 20

	results := make(chan string, numOrders)

	var wg sync.WaitGroup

	for i := 0; i < numOrders; i++ {
		wg.Add(1)
		go func(orderNum int) {
			defer wg.Done()

			orderID := fmt.Sprintf("order_%d_%d", orderNum, time.Now().UnixNano())

			req := &order.CreateOrderRequest{
				OrderId:       orderID,
				UserId:        fmt.Sprintf("user_%d", rand.Intn(10)),
				Items:         []string{"item1", "item2"},
				TotalAmount:   rand.Float64() * 100,
				PaymentMethod: "credit_card",
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			resp, err := client.CreateOrder(ctx, req)
			if err != nil {
				results <- fmt.Sprintf("Order %s failed: %v", orderID, err)
				return
			}

			if resp.Status == "confirmed" {
				results <- fmt.Sprintf("Order %s confirmed: %s", orderID, resp.Message)
			} else {
				results <- fmt.Sprintf("Order %s failed: %s", orderID, resp.Message)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
	}

	fmt.Println("All orders processed")
}
