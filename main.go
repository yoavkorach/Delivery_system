package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const N = 10
const MaxDistance = 10
const WaitTime = 1

var ID = 1
var mu sync.Mutex

type deliveryPerson struct { // Delivery person struct
	name      string
	currOrder Order
}

func CreateDeliveryPerson(name string) deliveryPerson { // Function to create delivery person
	return deliveryPerson{name, Order{}}
}

var availablePeople []deliveryPerson
var unavailablePeople []deliveryPerson

type Order struct { // Order struct
	OrderID   int
	OrderName string
	Distance  int
}

var orders []Order

func CreateOrder(orderID int, orderName string, distance int) Order { // Function to create order
	return Order{orderID, orderName, distance}
}

func AssignOrderToPerson(order Order) { // Function to assign order to delivery person
	mu.Lock()
	defer mu.Unlock()
	if len(availablePeople) == 0 { // If no delivery person is available then place order in queue
		orders = append(orders, order)
		fmt.Println("\nNo delivery person available at the moment. Order placed in queue")
	} else { // If delivery person is available then assign order to delivery person
		availablePeople[0].currOrder = order
		// Pass the index of the delivery person to the decreaser function
		delivery := availablePeople[0]
		go decreaser(delivery)
		unavailablePeople = append(unavailablePeople, availablePeople[0])
		availablePeople = availablePeople[1:]
		fmt.Printf("\nOrder assigned to delivery person, arrival time is %v\n", order.Distance)
	}
}
func decreaser(delivery deliveryPerson) { // Function to decrease the distance of order
	for delivery.currOrder.Distance > 0 { // Decrease the distance of order
		time.Sleep(WaitTime * time.Second)
		delivery.currOrder.Distance--
	}

	mu.Lock()
	defer mu.Unlock()

	// When order is delivered
	fmt.Println("\nOrder delivered")
	delivery.currOrder = Order{}
	availablePeople = append(availablePeople, delivery)
	for index, person := range unavailablePeople { // Remove the delivery person from unavailable list
		if person.name == delivery.name {
			unavailablePeople = append(unavailablePeople[:index], unavailablePeople[index+1:]...)
			break
		}
		// If there are orders in queue then assign the next order to the delivery person
		if len(orders) > 0 {
			nextOrder := orders[0]
			orders = orders[1:]
			AssignOrderToPerson(nextOrder)
		}
	}
}

func placeOrder(order string) { // Function to place order
	distance := rand.Intn(MaxDistance) + 1 // Ensure the distance is at least 1
	ID++
	fmt.Println("\nYour order ID is: ", ID) // Assign order to delivery person
	AssignOrderToPerson(CreateOrder(ID, order, distance))
}

func whereIsMyOrder(orderid int) { // Function to check the status of order
	mu.Lock()
	defer mu.Unlock()

	for _, order := range orders { // Check if order is in queue
		if order.OrderID == orderid {
			fmt.Println("\nYour order is in queue")
			return
		}
	}
	for _, person := range unavailablePeople { // Check if order is with delivery person
		if person.currOrder.OrderID == orderid {
			fmt.Println("\nYour order is with delivery person and will be delivered in", person.currOrder.Distance, "minutes")
			return
		}
	}
	fmt.Println("\nOrder not found")
}

func main() {
	for i := 0; i < N; i++ { // Create delivery people
		availablePeople = append(availablePeople, CreateDeliveryPerson(fmt.Sprintf("Person %v", i+1)))
	}
	reader := bufio.NewReader(os.Stdin)
	for { // Take request from user

		fmt.Println("\nPlease enter your request:\n 1. Place Order\n 2. Where is my Order\n 3. Exit")
		request, _ := reader.ReadString('\n')
		request = strings.TrimSpace(request)
		switch request {
		case "1":
			var order string
			fmt.Println("\nPlease place order for delivery: ") // Take order name
			order, _ = reader.ReadString('\n')
			order = strings.TrimSpace(order)
			placeOrder(order)
		case "2":
			var orderid int
			fmt.Println("Enter your order ID: ")
			orderidStr, _ := reader.ReadString('\n')
			orderidStr = strings.TrimSpace(orderidStr)
			orderid, err := strconv.Atoi(orderidStr)
			if err != nil {
				fmt.Println("Invalid order ID")
				continue
			}
			whereIsMyOrder(orderid)
		case "3":
			return
		default:
			fmt.Println("Invalid request")
		}
	}
}
