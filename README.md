# Smart Genset Monitoring System

A comprehensive solution for remote monitoring and management of diesel generators using IoT technology.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Components](#components)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [Troubleshooting](#troubleshooting)

---

## Features

### ✅ Real-time Monitoring
- **Engine Parameters**: RPM, Oil Temp, Coolant Temp, Battery Voltage, Oil Pressure
- **Power Metrics**: Voltage, Current, Power Factor, Frequency, Active/Reactive Power
- **Fuel Level**: Ultrasonic fuel sensor with tank calibration
- **Status**: Running, Stopped, Auto, Alarm

### 🔧 Remote Control
- **Start/Stop**: Remotely control the generator
- **Mode Switch**: Toggle between Auto and Manual modes
- **ATS Control**: Manage Automatic Transfer Switch

### 🔔 Notifications
- **Telegram Bot**: Real-time alerts for events and alarms
- **Email Alerts**: Email notifications for critical events
- **WhatsApp Integration** (future)

### 📊 Analytics & Reporting
- **Historical Data**: 1-year data retention
- **Load Profiling**: Analyze power consumption patterns
- **Maintenance Tracking**: Schedule and track maintenance

### 🔒 Security
- **Token-based Authentication**: Secure API access
- **User Roles**: Admin and Operator roles
- **Encryption**: HTTPS/TLS for data transmission

---

## Architecture

```
┌─────────────────────┐      ┌───────────────────────┐
│    Genset Controller│      │   Smart Controller    │
│ (Arduino/ESP32)     │      │ (Raspberry Pi)        │
└────────┬────────────┘      └────────┬──────────────┘
         │                        │
         │ 2G/4G (MQTT)           │ 2G/4G (MQTT)
         ▼                        ▼
  ┌────────────────────┐     ┌────────────────────┐
  │  MQTT Broker       │     │  Cloud Database  │
  │  (Cloud/Local)     │     │  (PostgreSQL)      │
  └─────────┬──────────┘     └─────────┬──────────┘
            │                          │
            ▼                          ▼
  ┌────────────────────┐     ┌────────────────────┐
  │   API Server       │     │  Web Dashboard   │
  │   (Node.js)        │     │  (React + Grafana) │
  └─────────┬──────────┘     └─────────┬──────────┘
            │                          │
            ▼                          ▼
  ┌────────────────────┐     ┌────────────────────┐
  │  Notification      │     │  Telegram Bot    │
  │  Service           │     │  (Node.js)         │
  └────────────────────┘     └────────────────────┘
```

---

## Getting Started

### Prerequisites

#### Hardware
- [ ] Genset Controller (Arduino/ESP32)
- [ ] Raspberry Pi 4/5
- [ ] GSM Module (SIM800L/SIM7600)
- [ ] Ultrasonic Sensor (for fuel)
- [ ] Current Transformers (CTs)
- [ ] Voltage Sensors
- [ ] Relay Modules
- [ ] Power Supply (24V/12V DC)

#### Software
- **Controller**: Arduino IDE 1.8.19+
- **API Server**: Node.js 18+, PostgreSQL 13+
- **Dashboard**: React 18+, Grafana 9+
- **Broker**: Mosquitto MQTT Broker
- **OS**: Raspberry Pi OS (64-bit recommended)

### Installation

#### 0. Docker Setup (Recommended)
The easiest way to run the entire backend infrastructure and API is using Docker Compose:

```bash
# Build and start all services
docker compose up -d --build

# View real-time logs
docker compose logs -f api
```

#### 1. Controller Setup

**Install Libraries:**
```bash
# Install Arduino IDE
# Install libraries: PubSubClient, ArduinoJson, ESP32Servo, etc.
```

**Connect Hardware:**
- Connect sensors to ADC pins
- Connect relays to digital pins
- Connect GSM module via Serial

**Upload Firmware:**
```c++
# Open 'controller/genset_controller/genset_controller.ino'
# Select board: ESP32 Dev Module
# Upload to device
```

#### 2. Raspberry Pi Setup

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Install Mosquitto
sudo apt install mosquitto mosquitto-clients -y

# Install Grafana
sudo apt-get install -y apt-transport-https software-properties-common wget
wget -q -O - https://apt.grafana.com/gpg.key | sudo apt-key add -

# Add Grafana repository
cat << EOF | sudo tee -a /etc/apt/sources.list.d/grafana.list
deb https://apt.grafana.com stable main
EOF

# Install Grafana
sudo apt update
sudo apt install grafana -y
```

#### 3. Database Setup

```sql
-- Create database
CREATE DATABASE genset_db;
CREATE USER genset_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE genset_db TO genset_user;

-- Run migrations
cd genset-monitoring/backend
npm install
# Follow migration instructions in migration.sql
```

#### 4. API Server Setup

```bash
cd genset-monitoring/backend
npm install

# Configure environment
cp .env.example .env
# Edit .env with your database credentials

# Create database tables
# Run migrations
npm run migrate

# Start server
npm run dev
```

#### 5. Dashboard Setup

**Grafana Configuration:**
```bash
# Start Grafana
sudo systemctl start grafana
# Access at http://[IP_ADDRESS]
# Login: admin / admin
```

**Import Dashboard:**
1. Download `grafana/dashboard.json`
2. In Grafana, click **+ > Import**
3. Upload JSON file
4. Select data source (PostgreSQL)
5. Click **Import**

#### 6. Telegram Bot Setup

```bash
cd genset-monitoring/backend/telegram-bot
npm install

# Configure environment
cp .env.example .env
# Edit .env with your Telegram Bot Token and Chat ID

# Start bot
npm run dev
```

---

## Usage

### Controller Status Indicators

| LED | Color | Status |
|-----|-------|--------|
| **D1** | Green | Power ON |
| **D2** | Green | GSM Connected |
| **D3** | Yellow | Generator Running |
| **D4** | Red | Low Battery |
| **D5** | Blue | WiFi Connected |
| **D6** | Red | Critical Alarm |

### Web Interface

**Access URLs:**
- **Backend API**: `http://localhost:8080`
- **Swagger UI**: `http://localhost:8080/swagger/index.html` (Interactive API Documentation)
- **Web Dashboard**: `http://[IP_ADDRESS]/dashboard`
- **Grafana**: `http://localhost:3000`
- **Log Viewer (Dozzle)**: `http://localhost:8083`

**Login Credentials:**
- **Admin**: admin / admin
- **Operator**: operator / operator

### Telegram Commands

Send these commands to the bot to control the genset:

| Command | Description |
|---------|-------------|
| `/status` | Get current genset status |
| `/start` | Start the generator |
| `/stop` | Stop the generator |
| `/auto` | Set to Auto mode |
| `/manual` | Set to Manual mode |
| `/ats_on` | Enable ATS |
| `/ats_off` | Disable ATS |
| `/help` | Show command list |

---

## Development

### Directory Structure

```
genset-monitoring/
├── backend/                   # Go API Server (Gin, GORM)
├── controller/                # Arduino/ESP32 firmware
├── docker/                    # Infrastructure config (Mosquitto)
└── docker-compose.yml         # Full stack orchestration
```

### Useful Commands (Backend)

Run these from the `backend/` directory:

| Command | Description |
|---------|-------------|
| `make swagger` | **Generate Swagger Documentation** |
| `make dev` | Run server locally with hot-reload |
| `make docker-up` | Start infrastructure + API in Docker |
| `make test` | Run unit tests |
| `make lint` | Run golangci-lint |
