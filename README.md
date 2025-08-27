# Telemetry Harbor OSS

<!-- OSS Badges -->
![License](https://img.shields.io/github/license/TelemetryHarbor/telemetry-harbor-oss.svg)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED.svg)
![Last Commit](https://img.shields.io/github/last-commit/TelemetryHarbor/telemetry-harbor-oss.svg)
![Issues](https://img.shields.io/github/issues/TelemetryHarbor/telemetry-harbor-oss.svg)
![Pull Requests](https://img.shields.io/github/issues-pr/TelemetryHarbor/telemetry-harbor-oss.svg)
![Repo Size](https://img.shields.io/github/repo-size/TelemetryHarbor/telemetry-harbor-oss.svg)
![Contributors](https://img.shields.io/github/contributors/TelemetryHarbor/telemetry-harbor-oss.svg)
<!-- Fun / Community -->
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)
![Stars](https://img.shields.io/github/stars/TelemetryHarbor/telemetry-harbor-oss.svg?style=social)
![Forks](https://img.shields.io/github/forks/TelemetryHarbor/telemetry-harbor-oss.svg?style=social)

Telemetry Harbor OSS is the open-source ingestion and visualization stack behind Telemetry Harbor. Self-host your own telemetry backend with full control over your data and infrastructure.

**_Repo Link:_** https://github.com/TelemetryHarbor/telemetry-harbor-oss



## What's Included

*   🚀 **API Ingestion Layer** - Go Fiber-based REST API with Redis-backed queue
*   ⚡ **Background Worker** - Efficient data processing and TimescaleDB insertion
*   📊 **Grafana Integration** - Pre-configured dashboards and datasource
*   🗄️ **TimescaleDB** - Optimized time-series database for telemetry data
*   🔄 **Redis Queue** - Reliable message queue management
*   🛠️ **SDK Compatible** - Works with all official Telemetry Harbor SDKs



## 🚀 Quick Start

Get your self-hosted Telemetry Harbor instance running in minutes:

```bash
git clone https://github.com/TelemetryHarbor/telemetry-harbor-oss.git
cd telemetry-harbor-oss
docker compose up -d
```

Once started:
*   **API available at** → `http://localhost:8000/api/v1`
*   **Grafana available at** → `http://localhost:3000` (default: `admin` / `StrongAdminPassword!`)


## 📋 System Requirements

*   **Docker** with Docker Compose support
*   **Minimum 2GB RAM** (4GB+ recommended for production)
*   **10GB+ disk space** (depends on data retention needs)



## 🔐 Security Configuration

:warning: Security Notice
This repository ships with default credentials for ease of testing. **Before using in production**, you must change the following in `docker-compose.yml`:

*   `POSTGRES_PASSWORD` - TimescaleDB password
*   `REDIS_PASSWORD` - Redis authentication password  
*   `API_KEY` - API authentication key for data ingestion
*   `GF_SECURITY_ADMIN_PASSWORD` - Grafana admin password

Failure to do so will leave your system vulnerable.


## 📡 API Ingestion

Replace Telemetry Harbor Cloud URLs with your own domain, omitting the harbor ID.

### Single Data Push
```
POST http://yourdomain.com/api/v1/ingest/
```

### Batch Data Push  
```
POST http://yourdomain.com/api/v1/ingest/batch
```

Both endpoints require the API key set via the `X-API-Key` header.

### Example Usage

**Single Data Point:**
```bash
curl -X POST "http://localhost:8000/api/v1/ingest/" \
-H "X-API-Key: your_api_key_here" \
-H "Content-Type: application/json" \
-d '{
  "time": "2025-01-18T19:24:00.948Z",
  "ship_id": "test_device_single",
  "cargo_id": "test_metric",
  "value": 123.45
}'
```

**Batch Data:**
```bash
curl -X POST "http://localhost:8000/api/v1/ingest/batch" \
-H "X-API-Key: your_api_key_here" \
-H "Content-Type: application/json" \
-d '[
  {
    "time": "2025-01-18T19:24:00.948Z",
    "ship_id": "batch_device",
    "cargo_id": "temperature",
    "value": 25.5
  },
  {
    "time": "2025-01-18T19:24:00.948Z", 
    "ship_id": "batch_device",
    "cargo_id": "humidity",
    "value": 60.2
  }
]'
```


## 📊 Visualization with Grafana

Grafana comes pre-configured with:
*   **Telemetry Harbor Datasource** (TimescaleDB connection)
*   **Comprehensive Telemetry Dashboard** 
*   **Ready-to-use panels** for time-series visualization

Log into Grafana at `http://localhost:3000` and start exploring your telemetry data immediately.



## 🗄️ Data Retention

By default, your telemetry data is kept for **365 days**.  
Want a different retention period? Just tweak it in [`init.sql`](https://github.com/TelemetryHarbor/telemetry-harbor-oss/blob/main/init.sql) before starting the stack.


## 🛠️ SDK Compatibility

Telemetry Harbor OSS is fully compatible with all official Telemetry Harbor SDKs:

*   [**Telemetry Harbor SDKs**](/docs/sdks)

Just replace your ingest endpoint with your OSS URL - no code changes needed!


## ☁️ OSS vs Cloud Comparison

| Feature | OSS Self-hosted | Telemetry Harbor Cloud |
|---------|----------------|-------------------------|
| **Pricing** | Free | Free/Paid |
| **Rate Limits** | No rate limits (your hardware) | Rate limits based on plan |
| **Data Retention** | 365 days (configurable) | Varies by plan |
| **Storage** | Based On Your Hardware | Unlimited |
| **Grafana Access** | Dedicated instance | Shared/dedicated |
| **Infrastructure** | Self-managed | Fully managed |
| **Harbor AI** | Not included | Available |
| **Backup & Recovery** | Not Available | All Plans |
| **High Availability** | Not Available | All Plans |
| **Scalability** | Manual scaling | Fully managed |
| **Support** | Community | Priority Email Support |
| **Updates** | Manual | Automatic |
| **Security** | Self-managed | Enterprise-grade |

### When to Choose OSS Self-hosted

**✅ Choose OSS if you:**
*   Need full control over your data and infrastructure
*   Want to customize the stack for specific needs
*   Have technical expertise to manage infrastructure
*   Require on-premises deployment for compliance

**✅ Choose Cloud if you:**
*   Want managed infrastructure and automatic updates
*   Need Harbor AI features for data insights
*   Require enterprise-grade backup and high availability
*   Prefer dedicated support and SLAs
*   Want to focus on your application, not infrastructure
*   Need guaranteed uptime and scalability



## 📜 License

This project is licensed under the **Apache License 2.0**.

*   ✅ **Free to use** (commercial + personal)
*   ✅ **Free to modify and redistribute**
*   ⚠️ **Must include attribution** (keep copyright + NOTICE file)

See [LICENSE](https://github.com/TelemetryHarbor/telemetry-harbor-oss/blob/main/LICENSE) for details.

If you use this project in your product, please credit Telemetry Harbor with a link to https://telemetryharbor.com.



## 🤝 Contributing

We welcome issues, pull requests, and feature suggestions!

*   **Open a GitHub issue** for bugs or feature requests
*   **Fork the repo** and submit PRs for improvements
*   **Join our community** discussions

**_Repo Link:_** https://github.com/TelemetryHarbor/telemetry-harbor-oss
