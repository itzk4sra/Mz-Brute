# 🚀 MarzbanBrute v0.2

**MarzbanBrute** is a high-performance, multi-threaded brute-force tool for testing login panels using a custom list of username-password combinations. It is designed for penetration testing and allows you to launch thousands of concurrent attempts to find valid credentials quickly and efficiently.

## ✨ Features

- 🔥 **Super-fast Concurrent Brute Force**: Supports up to **64,000 concurrent threads**.
- 📄 **Custom Combo Lists**: Easily load username-password pairs from a file.
- 📊 **Real-time Statistics**: Displays checks per second (CPS), valid credentials, failed attempts, and active thread count.
- 💾 **Logs Cracked Credentials**: Saves valid login details to `cracked.txt` for future use.

## ⚡ Quick Start

### 1. Clone the Repository:
```bash
git clone https://github.com/yourusername/marzbanbrute.git
cd marzbanbrute
```

### 2. Compile and Run:
```bash
go run main.go combo.txt <port>
```
combo.txt: A file containing username:password pairs for brute-force attempts.
<port>: The port of the target web panel.

### Example:
```bash
go run main.go combo.txt 8000
```
This command starts the brute force on port `8000` using combinations from `combo.txt`

## 🔧 Advanced Usage:
### Customizing Your Combo List
Prepare your combo list file in the following format:
```txt
user1:password1
user2:password2
...
```
Make sure the file is saved as `combo.txt` or another filename of your choice.
## Running in High Concurrency Mode

MarzbanBrute supports up to 64,000 concurrent threads. To efficiently use this feature, ensure your machine has the necessary resources (CPU, memory) to handle high concurrency without bottlenecks.

## Real-time Statistics
#### The program logs detailed statistics every second, showing:

📈 CPS (Checks per second): How many login attempts are made every second.
✅ Valid: The number of valid login credentials found.
❌ Fails: Number of failed attempts.
🧵 RunningThreads: The number of currently active brute-force threads.

## Valid Credentials Logging
#### Successfully cracked credentials are stored in the `cracked.txt` file, in the format:
```bash
<target-ip>:<username>:<password>
```
This allows easy reference and usage for further testing.

## Thread Management
By default, MarzbanBrute uses `MaxThreads` (set to 64,000) for handling concurrency. If you reach the limit, the program will wait for threads to complete before starting new ones. This keeps the system from being overloaded.

## 🚨 Disclaimer
This tool is intended for legal and authorized penetration testing only. Unauthorized use of this tool on systems you do not own or have explicit permission to test is illegal and unethical. Use responsibly!

## 👩‍💻 Authors

- **Mh-ProDev** – Development
- **ItzK4sra** – Development, Contributions

