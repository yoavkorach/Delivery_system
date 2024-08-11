package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const N = 10
const MaxDistance = 10
const WaitTime = 20

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
	fmt.Printf("\nOrder number %v delivered successfully", delivery.currOrder.OrderID)
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

func whereIsMyOrder(orderid int) string { // Function to check the status of order
	mu.Lock()
	defer mu.Unlock()

	for _, order := range orders { // Check if order is in queue
		if order.OrderID == orderid {
			return "Your order is in queue"
		}
	}
	for _, person := range unavailablePeople { // Check if order is with delivery person
		if person.currOrder.OrderID == orderid {
			return fmt.Sprintf("Your order is with delivery person and will be delivered in %d minutes", person.currOrder.Distance)
		}
	}
	return "Order not found"
}

// HTTP Handlers
func handlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // Check if request method is POST
		return
	}

	orderName := r.FormValue("order") // Get order name from request
	if orderName == "" {              // Check if order name is empty
		http.Error(w, "Order name is required", http.StatusBadRequest)
		return
	}

	placeOrder(orderName)
}

func handleWhereIsMyOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // Check if request method is GET
		return
	}

	orderIDStr := r.URL.Query().Get("orderid") // Get order ID from request
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	status := whereIsMyOrder(orderID)
	fmt.Fprintln(w, status)
}

func main() {
	for i := 0; i < N; i++ { // Create delivery people
		availablePeople = append(availablePeople, CreateDeliveryPerson(fmt.Sprintf("Person %v", i+1)))
	}

	// HTTP Routes
	http.HandleFunc("/placeOrder", handlePlaceOrder)
	http.HandleFunc("/whereIsMyOrder", handleWhereIsMyOrder)

	// Start HTTP server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
