package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func createBanner() {
	banner := `


	████████╗ ██████╗██████╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗      ██████╗  ██████╗ 
	╚══██╔══╝██╔════╝██╔══██╗    ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗    ██╔════╝ ██╔═══██╗
	   ██║   ██║     ██████╔╝    ███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝    ██║  ███╗██║   ██║
	   ██║   ██║     ██╔═══╝     ╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗    ██║   ██║██║   ██║
	   ██║   ╚██████╗██║         ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║    ╚██████╔╝╚██████╔╝
	   ╚═╝    ╚═════╝╚═╝         ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝     ╚═════╝  ╚═════╝ 
																												   
	----------------------------------------------------------------------------------------------------------------


	██╗██████╗ ███████╗ █████╗ ███╗   ██╗    ████████╗ █████╗ ███╗   ██╗███████╗██████╗ 
	██║██╔══██╗██╔════╝██╔══██╗████╗  ██║    ╚══██╔══╝██╔══██╗████╗  ██║██╔════╝██╔══██╗
	██║██████╔╝█████╗  ███████║██╔██╗ ██║       ██║   ███████║██╔██╗ ██║█████╗  ██████╔╝
	██║██╔══██╗██╔══╝  ██╔══██║██║╚██╗██║       ██║   ██╔══██║██║╚██╗██║██╔══╝  ██╔══██╗
	██║██║  ██║██║     ██║  ██║██║ ╚████║       ██║   ██║  ██║██║ ╚████║███████╗██║  ██║
	╚═╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝╚═╝  ╚═══╝       ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
																						
                                                                                    
                                                                                                                                
                                                                            
`
	fmt.Println(banner)

}

func portScan(ip string, port int, wg *sync.WaitGroup, results chan int) {
	defer wg.Done()

	network := "tcp"
	address := fmt.Sprintf("%s:%d", ip, port)
	connection, err := net.DialTimeout(network, address, 2*time.Second)
	if err != nil {
		return
	}

	fmt.Printf("Açık Port Bulundu: %d\n", port)
	results <- port
	connection.Close()
}

func main() {
	createBanner()
	var targetIP string

	fmt.Print("Hedef IP adresini girin: ")
	fmt.Scan(&targetIP)

	if net.ParseIP(targetIP) == nil {
		fmt.Println("Geçersiz IP adresi. Program sonlandırılıyor.")
		return
	}

	var maxParallelScans int
	flag.IntVar(&maxParallelScans, "p", 100, "Paralel tarama sayısı")

	var timeout int
	flag.IntVar(&timeout, "t", 2, "Bağlantı zaman aşımı (saniye)")

	flag.Parse()

	var wg sync.WaitGroup
	results := make(chan int)

	fmt.Println("Port taraması başlatılıyor...")

	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		go func(p int) {
			portScan(targetIP, p, &wg, results)

			if p%maxParallelScans == 0 {
				fmt.Printf("Tarama devam ediyor... (%d / %d)\n", p, 65535)
				time.Sleep(time.Duration(timeout) * time.Second)
			}
		}(port)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	openPorts := []int{}
	for port := range results {
		openPorts = append(openPorts, port)
	}

	if len(openPorts) > 0 {
		fmt.Println("\nAçık Portlar:")
		for _, port := range openPorts {
			fmt.Println(port)
		}
	} else {
		fmt.Println("\nHedefte açık port bulunamadı.")
	}
}
