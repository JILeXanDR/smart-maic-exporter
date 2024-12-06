package main

//
// // Configuration variables
// var (
// 	dataSourceURL = getEnv("DATA_SOURCE_URL", "http://192.168.10.55/?page=getwdata")
// 	exporterPort  = getEnvAsInt("EXPORTER_PORT", 8000)
// )
//
// // Data fetching mutex
// var mutex sync.Mutex
//
// func main() {
// 	if !(len(dataSourceURL) > 0 && (dataSourceURL[:7] == "http://" || dataSourceURL[:8] == "https://")) {
// 		log.Fatal("Invalid DATA_SOURCE_URL. Must start with 'http://' or 'https://'.")
// 	}
//
// 	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
// 		fetchAndUpdateMetrics() // Fetch and update metrics on request
// 		handler := promhttp.HandlerFor(CustomRegistry, promhttp.HandlerOpts{})
// 		handler.ServeHTTP(w, r) // Serve updated metrics
// 	})
//
// 	log.Printf("Starting Smart Maic Exporter on port %d...", exporterPort)
//
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", exporterPort), nil))
// }
//
// // fetchAndUpdateMetrics fetches data from the device and updates Prometheus metrics
// func fetchAndUpdateMetrics() {
// 	mutex.Lock()
// 	defer mutex.Unlock()
//
// 	started := time.Now()
//
// 	log.Println("Fetching data from the Smart Maic device...")
//
// 	resp, err := http.Get(dataSourceURL)
// 	if err != nil {
// 		log.Printf("Error fetching data: %v", err)
// 		SetDeviceAPIStatus(DeviceAPIStatusOffline)
// 		return
// 	}
// 	defer resp.Body.Close()
//
// 	log.Printf("Data from the Smart Maic device has been fetched in %s...", time.Since(started))
//
// 	if resp.StatusCode == http.StatusTooManyRequests {
// 		log.Println("Received 429 Too Many Requests from the API.")
// 		SetDeviceAPIStatus(DeviceAPIStatusTooManuRequests)
// 		return
// 	}
//
// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("Unexpected status code: %d", resp.StatusCode)
// 		SetDeviceAPIStatus(DeviceAPIStatusOffline)
// 		return
// 	}
//
// 	var responseData T
// 	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
// 		log.Printf("Error decoding JSON: %v", err)
// 		SetDeviceAPIStatus(DeviceAPIStatusOffline)
// 		return
// 	}
//
// 	// Update metrics
//
// 	SetMetrics(responseData)
//
// 	log.Println("Successfully updated metrics.")
// }
//
