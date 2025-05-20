# SSH Cracker

[فارسی](README.fa.md) | English

A powerful SSH connection testing tool (SSH-Cracker) written in Go. This tool allows you to test SSH connections with multiple username/password combinations against a list of IP addresses.

## ⚠️ Disclaimer

This tool is for educational and testing purposes only. Always ensure you have proper authorization before testing any SSH connections. The author is not responsible for any misuse of this tool.

## 🚀 Features

- Multi-threaded SSH connection testing
- Support for multiple username/password combinations
- Real-time progress monitoring
- Discord webhook integration for successful connections
- Cross-platform support (Windows & Linux)
- License system for access control

## 📋 Prerequisites

- Go 1.16 or higher
- Python 3.7 or higher (for the license server)
- Git

### Go Dependencies

The project uses the following Go packages:
- `golang.org/x/crypto/ssh` - For SSH connection handling

To install dependencies:
```bash
go mod download
```

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/ssh-cracker.git
cd ssh-cracker
```

### 2. Configure the License Server

1. Edit `app.py` and set your admin password:
```python
PASSWORD = 'YOUR_ADMIN_PASSWORD'  # Change this to a secure password
```

2. Start the license server:
```bash
python app.py
```

### 3. Build the SSH Cracker

#### For Linux:
```bash
go build -o ssh-cracker ssh.go
```

#### For Windows:
```bash
go build -o ssh-cracker.exe ssh.go
```

## 🔧 Configuration

Before running the tool, you need to configure a few things:

1. Edit `ssh.go` and update the following constants:
```go
const (
    APIEndpoint = "http://your-api-endpoint.com:8000/check-license"
    WebhookURL  = "https://discord.com/api/webhooks/your-webhook-url"
)
```

2. Create your username and password lists:
   - Create a file with usernames (one per line)
   - Create a file with passwords (one per line)

## 🚀 Usage

1. Run the program:
```bash
# Linux
./ssh-cracker

# Windows
ssh-cracker.exe
```

2. When prompted:
   - Enter your license key
   - Provide the path to your username list file
   - Provide the path to your password list file
   - Enter the IP list file path
   - Set the timeout value (in seconds)
   - Set the maximum number of concurrent connections

3. The program will:
   - Create a combination file of usernames and passwords
   - Start testing SSH connections
   - Display real-time progress
   - Save successful connections to `su-goods.txt`
   - Send notifications to Discord (if configured)

## 📊 Output

The program provides real-time information including:
- Total connections checked
- Connection speed (IP/s)
- Elapsed time
- Remaining time
- Number of successful connections

Successful connections are saved in `su-goods.txt` in the format:
```
IP:PORT@USERNAME:PASSWORD
```

## 🔒 License System

The tool uses a license system to control access. To create a license:

1. Access the license server API:
```
http://your-api-endpoint.com:8000/create-lic?password=YOUR_ADMIN_PASSWORD&name=USER_NAME&expire=DAYS
```

2. Add IP addresses to the license:
```
http://your-api-endpoint.com:8000/add-ip?password=YOUR_ADMIN_PASSWORD&lic=LICENSE_KEY&ip=IP_ADDRESS
```

## 🔒 License System API Documentation

The license system provides several endpoints for managing licenses and IP addresses. All endpoints require authentication using the admin password.

### Authentication
All API endpoints require the admin password to be passed as a query parameter:
```
?password=YOUR_ADMIN_PASSWORD
```

### Endpoints

#### 1. Create License
Creates a new license for a user.

```
GET /create-lic
```

Parameters:
- `password` (required): Admin password
- `name` (required): User's name
- `expire` (required): Number of days until license expiration

Response:
```json
{
    "license": "GENERATED_LICENSE_KEY",
    "name": "USER_NAME",
    "ip_list": [],
    "expire": "YYYY-MM-DD"
}
```

#### 2. Add IP Address
Adds an IP address to a license's allowed IP list.

```
GET /add-ip
```

Parameters:
- `password` (required): Admin password
- `lic` (required): License key
- `ip` (required): IP address to add

Response:
```json
{
    "message": "IP added successfully"
}
```

#### 3. Clear IP List
Clears all IP addresses from a license's IP list.

```
GET /clear-ips
```

Parameters:
- `password` (required): Admin password
- `lic` (required): License key

Response:
```json
{
    "message": "IP list cleared successfully"
}
```

#### 4. Delete License
Deletes a license completely.

```
GET /delete-lic
```

Parameters:
- `password` (required): Admin password
- `lic` (required): License key

Response:
```json
{
    "message": "License deleted successfully"
}
```

#### 5. Get All Licenses
Retrieves information about all licenses.

```
GET /all
```

Parameters:
- `password` (required): Admin password

Response:
```json
{
    "LICENSE_KEY_1": {
        "name": "USER_NAME",
        "ip_list": ["IP1", "IP2"],
        "expire": "YYYY-MM-DD"
    },
    "LICENSE_KEY_2": {
        ...
    }
}
```

#### 6. Check License
Validates a license and checks if an IP is allowed.

```
POST /check-license
```

Request Body:
```json
{
    "lic": "LICENSE_KEY",
    "ip": "IP_ADDRESS"
}
```

Response:
```json
{
    "name": "USER_NAME",
    "ip_list": ["IP1", "IP2"],
    "expire": "YYYY-MM-DD"
}
```

### Error Codes

- `400`: Bad Request - Missing or invalid parameters
- `401`: Unauthorized - Invalid admin password
- `403`: Forbidden - IP not allowed or license expired
- `404`: Not Found - License not found
- `409`: Conflict - Name already exists or IP already in list

### Security Notes

1. Always use HTTPS in production
2. Change the default admin password
3. Regularly monitor license usage
4. Implement rate limiting for API endpoints
5. Keep the server's Python packages updated

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

- **SudoLite** - *Initial work*

## 🙏 Acknowledgments

- Thanks to all contributors
- Inspired by the need for efficient SSH connection testing 