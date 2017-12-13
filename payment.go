package main

import ("fmt"
        "log"
        "net/http"
        "encoding/json"
        "bytes"
        "math/rand"
        "github.com/gorilla/mux"
        "gopkg.in/mgo.v2"

        "flag"
        "github.com/opentracing/opentracing-go"
        "sourcegraph.com/sourcegraph/appdash"
        appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
      )

var (
	port           = flag.Int("port", 5005, "Example app port.")
	appdashPort    = flag.Int("appdash.port", 8700, "Run appdash locally on this port.")
)



const SERVER = "paymentdb:27017"
const DBNAME = "paymentDb"
const DOCNAME = "payment"

type Order struct {
  OrderId int
}

type Payment struct {
  PaymentId int `bson:"_id"`
  Order_Id int
}


func main() {

  flag.Parse()
  var tracer opentracing.Tracer
  addr := startAppdashServer(*appdashPort)
  tracer = appdashot.NewTracer(appdash.NewRemoteCollector(addr))
  opentracing.InitGlobalTracer(tracer)
  addr = fmt.Sprintf(":%d", *port)

  router := mux.NewRouter()
  router.HandleFunc("/payment", payment).Methods("POST")
  log.Fatal(http.ListenAndServe(":5005", router))
}

func payment(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Request started"))
  sp := opentracing.StartSpan("POST /payment") // Start a new root span.
  defer sp.Finish()


  session, err := mgo.Dial(SERVER)
  if err != nil {
    fmt.Println("Failed to establish connection to Mongo server:", err)
  }
  defer session.Close()
  c := session.DB(DBNAME).C(DOCNAME)
  OrderId := Order{}
  err = json.NewDecoder(r.Body).Decode(&OrderId)
  if err != nil{
    panic(err)
  }
  orderId := OrderId.OrderId
  for {
    paymentId := rand.Intn(1000)
    err = c.Insert(&Payment{PaymentId: paymentId, Order_Id: orderId})
    if err != nil {
   		log.Fatal("Unable to insert", err)
      continue
   	}
    break
  }
  orderIdJson, err := json.Marshal(OrderId)
  if err != nil{
    panic(err)
  }
  url := "http://orders:5004/update-order-status"
  resp, err := http.Post(url, "application/json", bytes.NewBuffer(orderIdJson))
  if resp.StatusCode == 200 {
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("HTTP status code returned!"))
  }
  w.Write([]byte("... done!"))
}
