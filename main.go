package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Println("Error loading app.env")
	}

	hostname, _ := os.Hostname()
	processor := getProcessorInfo()
	processList := getProcessList()
	userList := getUserList()
	osName := getOSName()
	osVersion := getOSVersion()

	// Create map
	sysInfo := map[string]interface{}{
		"hostname":  hostname,
		"processor": processor,
		"processes": processList,
		"users":     userList,
		"os_name":   osName,
		"os_ver":    osVersion,
	}

	// Convert to JSON
	jsonStr, _ := json.Marshal(sysInfo)

	// Create HTTP POST
	req, err := http.NewRequest("POST", os.Getenv("API_URL")+"/data", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status: %s", resp.Status)
	}

	fmt.Println("System information sent to API successfully.")
}

func getProcessorInfo() string {
	cmd := exec.Command("cat", "/proc/cpuinfo")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}

func getProcessList() string {
	cmd := exec.Command("ps", "-eo", "pid,ppid,user,%cpu,%mem,cmd")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}

func getUserList() string {
	cmd := exec.Command("w", "-h")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}

func getOSName() string {
	cmd := exec.Command("uname", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}

func getOSVersion() string {
	cmd := exec.Command("uname", "-r")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}
