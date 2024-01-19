package main

import (
	"database/sql"
	"fmt"
	"github.com/op/go-logging"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func primeSum(num int) int {
	sum := 0
	for i := 2; i < num; i++ {
		if isPrime(i) {
			sum += i
		}
	}
	return sum
}

var (
	count = int64(0)
	lock  sync.Mutex

	db          *sql.DB
	redisClient *redis.Client
)

func init() {
	// Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// MySQL
	var err error
	db, err = sql.Open("mysql", "root@tcp(db)/information_schema")
	if err != nil {
		panic(err.Error())
	}
}

func add(w http.ResponseWriter) int64 {
	lock.Lock()
	defer lock.Unlock()

	if count%3000 == 0 {
		fmt.Fprintf(w, "============== Count =============\n")
		fmt.Fprintf(w, "count: %d\n", count)
		fmt.Fprintf(w, "==================================\n\n")
	}
	count++
	return count
}
func testMix(w http.ResponseWriter, req *http.Request) {
	//
	//defer redisClient.Close()
	now := time.Now().Unix()

	val, err := redisClient.Set("key", now, time.Second).Result()
	if err != nil {
		panic(err.Error())
	}

	val, err = redisClient.Get("key").Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "============== Redis =============\n")
	fmt.Fprintf(w, "value: %s\n", val)
	fmt.Fprintf(w, "==================================\n\n")

	rows, err := db.Query("SELECT TABLE_NAME FROM tables limit 3")
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "============== MySQL =============\n")
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintf(w, "table name: %s\n", tableName)
	}
	fmt.Fprintf(w, "==================================\n\n")

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true
	c := &http.Client{Transport: t}
	resp, err := c.Get("http://web")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "============== HTTP =============\n")
	fmt.Fprintf(w, "StatusCode: %d\n", resp.StatusCode)
	fmt.Fprintf(w, "==================================\n\n")

	fmt.Fprintf(w, "============== HTTP =============\n")
	fmt.Fprintf(w, "StatusCode: %d\n", 200)
	fmt.Fprintf(w, "==================================\n\n")

	//// 计算素数和,模仿耗时操作
	//query := req.URL.Query()
	//strNum := query.Get("num")
	//num, err := strconv.Atoi(strNum)
	//if err != nil {
	//	fmt.Fprintf(w, "============== Sum =============\n")
	//	fmt.Fprintf(w, "sum: %d\n", -1)
	//	fmt.Fprintf(w, "==================================\n\n")
	//
	//} else {
	//	sum := primeSum(num)
	//	fmt.Fprintf(w, "============== Sum =============\n")
	//	fmt.Fprintf(w, "sum: %d\n", sum)
	//	fmt.Fprintf(w, "==================================\n\n")
	//}
}

func SieveOfEratosthenes(n int) []int {
	prime := make([]bool, n+1)
	for i := range prime {
		prime[i] = true
	}

	for p := 2; p*p <= n; p++ {
		if prime[p] {
			for i := p * p; i <= n; i += p {
				prime[i] = false
			}
		}
	}

	primes := make([]int, 0)
	for p := 2; p <= n; p++ {
		if prime[p] {
			primes = append(primes, p)
		}
	}

	return primes
}
func testPrime(w http.ResponseWriter, req *http.Request) {
	primes := SieveOfEratosthenes(100000000)
	fmt.Fprintf(w, "============== Count =============\n")
	fmt.Fprintf(w, "count: %d\n", len(primes))
	fmt.Fprintf(w, "==================================\n\n")
}

func testMySQL(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT TABLE_NAME FROM tables limit 3")
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "============== MySQL =============\n")
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintf(w, "table name: %s\n", tableName)
	}
	fmt.Fprintf(w, "==================================\n\n")
}

func testApi(w http.ResponseWriter, req *http.Request) {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true
	c := &http.Client{Transport: t}
	resp, err := c.Get("http://web")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "============== HTTP =============\n")
	fmt.Fprintf(w, "StatusCode: %d\n", resp.StatusCode)
	fmt.Fprintf(w, "==================================\n\n")
}

func testRedis(w http.ResponseWriter, req *http.Request) {
	//
	//defer redisClient.Close()
	now := time.Now().Unix()

	val, err := redisClient.Set("key", now, time.Second).Result()
	if err != nil {
		panic(err.Error())
	}

	val, err = redisClient.Get("key").Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "============== Redis =============\n")
	fmt.Fprintf(w, "value: %s\n", val)
	fmt.Fprintf(w, "==================================\n\n")

	// HTTP
	// 客户端默认会启用连接池,会重用连接和协程,导致追踪失败.
	// 关闭 HTTP KeepAlive
	if count%1000 == 0 {

	}

	fmt.Fprintf(w, "============== HTTP =============\n")
	fmt.Fprintf(w, "StatusCode: %d\n", 200)
	fmt.Fprintf(w, "==================================\n\n")

	//// 计算素数和,模仿耗时操作
	//query := req.URL.Query()
	//strNum := query.Get("num")
	//num, err := strconv.Atoi(strNum)
	//if err != nil {
	//	fmt.Fprintf(w, "============== Sum =============\n")
	//	fmt.Fprintf(w, "sum: %d\n", -1)
	//	fmt.Fprintf(w, "==================================\n\n")
	//
	//} else {
	//	sum := primeSum(num)
	//	fmt.Fprintf(w, "============== Sum =============\n")
	//	fmt.Fprintf(w, "sum: %d\n", sum)
	//	fmt.Fprintf(w, "==================================\n\n")
	//}
}

//	func testPrimeSum(w http.ResponseWriter, req *http.Request) {
//		query := req.URL.Query()
//		strNum := query.Get("num")
//		primeSum(100000)
//	}
func testLog(w http.ResponseWriter, req *http.Request) {
	var a bool = true
	var b int = -1
	var c int8 = -2
	var d int16 = -3
	var e int32 = -4
	var f int64 = -5
	var g uint = 6
	var h uint8 = 7
	var i uint16 = 8
	var j uint32 = 9
	var k uint64 = 10
	var l uintptr = 11
	var m float32 = 12.12345678
	var n float64 = 13.12345678
	var o string = "string"
	log.Println(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o)

}

func testLog1(w http.ResponseWriter, req *http.Request) {
	var a bool = true
	var b int = -1
	var c int8 = -2
	var d int16 = -3
	var e int32 = -4
	var f int64 = -5
	var g uint = 6
	var h uint8 = 7
	var i uint16 = 8
	var j uint32 = 9
	var k uint64 = 10
	var l uintptr = 11
	var m float32 = 12.12345678
	var n float64 = 13.12345678
	var o string = "string"
	fmt.Print(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o)
	fmt.Printf("a=%v b=%v o=%v n=%v", a, b, o, n)
}

type Person struct {
	name    string
	age     int
	address string
}

func (p Person) String() string {
	return fmt.Sprintf("Name: %s, Age: %d, Address: %s", p.name, p.age, p.address)
}

func testLog2(w http.ResponseWriter, req *http.Request) {
	p := Person{"John", 30, "123 Main St."}
	log.Println(p)
}

func testLog3(w http.ResponseWriter, req *http.Request) {
	p := Person{"John", 30, "123 Main St."}
	log.Println(p.String())
}

func testLog4(w http.ResponseWriter, req *http.Request) {
	log := logging.MustGetLogger("example")

	var a bool = true
	var b int = -1
	var n float64 = 13.12345678
	var o string = "string"

	log.Infof("a=%v b=%v o=%v n=%v", a, b, o, n)
}

func main() {
	// HTTP handler
	http.HandleFunc("/nginx", testApi)
	http.HandleFunc("/mix", testMix)
	http.HandleFunc("/mysql", testMySQL)
	http.HandleFunc("/redis", testRedis)
	http.HandleFunc("/prime", testPrime)
	http.HandleFunc("/log", testLog)
	http.HandleFunc("/log1", testLog1)
	http.HandleFunc("/log2", testLog2)
	http.HandleFunc("/log3", testLog3)
	http.HandleFunc("/log4", testLog4)

	// Start HTTP server in a Goroutine
	go func() {
		fmt.Println("Listening on localhost:8080")
		http.ListenAndServe(":8080", nil)
	}()

	// TLS
	//certPath := "/app/certificate/certificate.pem"
	//keyPath := "/app/certificate/private-key.pem"

	//// Start HTTPS server in a Goroutine
	//go func() {
	//	fmt.Println("Listening on localhost:443")
	//	err := http.ListenAndServeTLS(":443", certPath, keyPath, nil)
	//	if err != nil {
	//		fmt.Println("HTTPS server error:", err)
	//	}
	//}()
	select {}
}
