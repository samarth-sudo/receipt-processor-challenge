# receipt-processor-challenge
# Receipt Processor API

## Overview
This is a **Receipt Processor API** built in **Go (Golang)** that processes receipts and awards points based on specific rules. The API supports two endpoints:

1. **POST `/receipts/process`** - Submits a receipt and returns a unique receipt ID.
2. **GET `/receipts/{id}/points`** - Retrieves the points awarded for a given receipt ID.

The application runs in-memory, meaning data is not persisted across restarts.

## Features
- Process receipts and generate a unique ID.
- Calculate points based on predefined rules.
- Retrieve points for a processed receipt.
- Dockerized setup for easy deployment.

---

## **Tech Stack**
- **Go (Golang)** 1.18+
- **UUID Library** for unique receipt IDs.
- **Docker & Docker Compose** for containerized deployment.

---

## **Setup & Installation**

### **1️⃣ Prerequisites**
Ensure you have the following installed:
- [Go](https://go.dev/doc/install) (1.18 or later)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
---
### **2️⃣ Clone the Repository**
```sh
git clone <repository_url>
cd receipt-processor
```
---
## **Docker Deployment**
### **1️⃣ Build and Run the Container**
```sh
docker-compose up --build
```

### **2️⃣ Stop the Container**
```sh
docker-compose down
```
---

### **3️⃣ Install Dependencies To Run Locally**
```sh
go mod tidy
```

### **4️⃣ Run the Server Locally**
```sh
go run main.go
```
The server will start on `http://localhost:8080`

---

## **API Endpoints**

### **1️⃣ Process Receipt (`POST /receipts/process`)**
**Request:**
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
    { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
  ],
  "total": "18.74"
}
```

**Response:**
```json
{ "id": "some-uuid" }
```

---

### **2️⃣ Get Points (`GET /receipts/{id}/points`)**
**Request:**
```sh
GET /receipts/some-uuid/points
```

**Response:**
```json
{ "points": 32 }
```

---

## **Rules for Calculating Points**
- 1 point for each alphanumeric character in the retailer name.
- 50 points if the total is a round dollar amount.
- 25 points if the total is a multiple of `0.25`.
- 5 points for every two items on the receipt.
- If the item description length is a multiple of 3, award `ceil(price * 0.2)` points.
- 6 points if the purchase day is an odd number.
- 10 points if the purchase time is between `2:00 PM - 4:00 PM`.

---

## **Project Structure**
```
receipt-processor/
├── main.go               # Main application logic
├── go.mod                # Go module file
├── go.sum                # Dependencies checksum
├── Dockerfile            # Docker setup
├── docker-compose.yml    # Docker Compose config
├── README.md             # Documentation
```

---

## **Troubleshooting**

### **Port Already in Use**
Error:
```sh
bind: address already in use
```
Solution:
```sh
lsof -i :8080  # Find the process using port 8080
kill -9 <PID>  # Kill the process
```
Alternatively, change the port in `docker-compose.yml`:
```yaml
ports:
  - "9090:8080"
```

### **Go Module Issues**
Run:
```sh
go mod tidy
```

---

## **Contributors**
- Samarth Singh

---



