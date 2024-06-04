package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"ex0/cache"
)

type Server struct {
	Cache *cache.Cache
}

func (s *Server) Serve() {
	tmpl := template.Must(template.ParseFiles("frontend/frontend.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			http.Error(w, "ID parameter is missing", http.StatusBadRequest)
			return
		}

		order, found := s.Cache.GetOrderFromCache(orderID)
		if !found {
			http.Error(w, "Order not found in cache", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		orderTemplate := `
            <h2>Model</h2>
            <table>
                <tr><th>Order UID</th><td>{{.OrderUid}}</td></tr>
                <tr><th>Track Number</th><td>{{.TrackNumber}}</td></tr>
                <tr><th>Entry</th><td>{{.Entry}}</td></tr>
                <tr><th>Locale</th><td>{{.Locale}}</td></tr>
                <tr><th>Internal Signature</th><td>{{.InternalSignature}}</td></tr>
                <tr><th>Customer ID</th><td>{{.CustomerId}}</td></tr>
                <tr><th>Delivery Service</th><td>{{.DeliveryService}}</td></tr>
                <tr><th>Shardkey</th><td>{{.Shardkey}}</td></tr>
                <tr><th>SM ID</th><td>{{.SmId}}</td></tr>
                <tr><th>Date Created</th><td>{{.DateCreated}}</td></tr>
                <tr><th>OOF Shard</th><td>{{.OofShard}}</td></tr>
            </table>

            <h2>Delivery</h2>
            <table>
                <tr><th>Name</th><td>{{.Dev.Name}}</td></tr>
                <tr><th>Phone</th><td>{{.Dev.Phone}}</td></tr>
                <tr><th>Zip</th><td>{{.Dev.Zip}}</td></tr>
                <tr><th>City</th><td>{{.Dev.City}}</td></tr>
                <tr><th>Address</th><td>{{.Dev.Address}}</td></tr>
                <tr><th>Region</th><td>{{.Dev.Region}}</td></tr>
                <tr><th>Email</th><td>{{.Dev.Email}}</td></tr>
            </table>

            <h2>Payment</h2>
            <table>
                <tr><th>Transaction</th><td>{{.Pay.Transaction}}</td></tr>
                <tr><th>Request ID</th><td>{{.Pay.RequestId}}</td></tr>
                <tr><th>Currency</th><td>{{.Pay.Currency}}</td></tr>
                <tr><th>Provider</th><td>{{.Pay.Provider}}</td></tr>
                <tr><th>Amount</th><td>{{.Pay.Amount}}</td></tr>
                <tr><th>Payment Date</th><td>{{.Pay.PaymentDt}}</td></tr>
                <tr><th>Bank</th><td>{{.Pay.Bank}}</td></tr>
                <tr><th>Delivery Cost</th><td>{{.Pay.DeliveryCost}}</td></tr>
                <tr><th>Goods Total</th><td>{{.Pay.GoodsTotal}}</td></tr>
                <tr><th>Custom Fee</th><td>{{.Pay.CustomFee}}</td></tr>
            </table>

            <h2>Items</h2>
            {{range .It}}
            <table>
                <tr><th>Chrt ID</th><td>{{.ChrtId}}</td></tr>
                <tr><th>Track Number</th><td>{{.TrackNumber}}</td></tr>
                <tr><th>Price</th><td>{{.Price}}</td></tr>
                <tr><th>RID</th><td>{{.Rid}}</td></tr>
                <tr><th>Name</th><td>{{.Name}}</td></tr>
                <tr><th>Sale</th><td>{{.Sale}}</td></tr>
                <tr><th>Size</th><td>{{.Size}}</td></tr>
                <tr><th>Total Price</th><td>{{.TotalPrice}}</td></tr>
                <tr><th>NM ID</th><td>{{.NmId}}</td></tr>
                <tr><th>Brand</th><td>{{.Brand}}</td></tr>
                <tr><th>Status</th><td>{{.Status}}</td></tr>
            </table>
            {{end}}
        `
		tmpl, err := template.New("order").Parse(orderTemplate)
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, order)
	})

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
