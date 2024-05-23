package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Customer struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

type CustomerRepository struct {
	customers map[string]Customer
}

func NewCustomerRepository() CustomerRepository {
	return CustomerRepository{
		customers: make(map[string]Customer),
	}
}

func (cr *CustomerRepository) Save(customer Customer) error {

	cr.customers[customer.Name] = customer

	return nil
}

func saveCustomer(store CustomerRepository, customer Customer) error {
	err := store.Save(customer)

	return err
}

type CustomerApi struct{}

func (ca CustomerApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "", "GET":
		ca.Get(w, r)
	default:
		slog.Info("Unrecognized method", "method", r.Method)
	}
}

func (ca CustomerApi) Get(w http.ResponseWriter, r *http.Request) {

	e, ok := FromExecutorContext(r.Context())

	if !ok {
		internalServerError(w)
		return
	}

	brokers := e.GetAvailableBrokers()

	slog.Info("Found some brokers", "brokers", brokers)

	data := Customer{
		Name: "TestUser",
		Age:  12,
	}
	serialized, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	okJsonResponse(w, serialized)
}
