# 🚀 Deployment Guide

## Table of Contents
1. [Docker Deployment](#docker-deployment)
2. [Cloud Linux Environment Variables](#cloud-linux-environment-variables)
3. [Production Deployment](#production-deployment)
4. [Environment Configuration](#environment-configuration)
5. [Monitoring & Logging](#monitoring--logging)

## Docker Deployment

### 1. Local Development with Docker Compose

```bash
# Setup environment
make setup

# Edit .env file with your configuration
vim .env

# Start all services (PostgreSQL, Redis, App)
make docker-compose-up

# View logs
make logs

# Stop services
make docker-compose-down
```

### 2. Production Docker Deployment

```bash
# Build optimized image
make docker-build

# Run single container
make docker-run

# Or with nginx proxy
make docker-compose-nginx
```

## Cloud Linux Environment Variables

### ✅ **Có hỗ trợ set environment variables trên cloud Linux!**

Trên các cloud Linux platforms, bạn có nhiều cách để set environment variables thay cho file `.env`:

### 1. **System Environment Variables**
```bash
# Set permanent environment variables
echo 'export DB_HOST=your-db-host' >> ~/.bashrc
echo 'export DB_PASSWORD=your-password' >> ~/.bashrc
source ~/.bashrc

# Or system-wide
sudo vim /etc/environment
# Add: DB_HOST=your-db-host
```

### 2. **Systemd Service (Recommended for Production)**
```bash
# Create systemd service file
sudo vim /etc/systemd/system/codebase-golang.service
```

```ini
[Unit]
Description=CodeBase Golang Application
After=network.target

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/codebase-golang
ExecStart=/opt/codebase-golang/codebase-golang
Restart=always
RestartSec=5

# Environment Variables
Environment=DB_HOST=your-db-host
Environment=DB_PORT=5432
Environment=DB_USER=postgres
Environment=DB_PASSWORD=your-secure-password
Environment=DB_NAME=codebase_db
Environment=REDIS_HOST=your-redis-host
Environment=REDIS_PORT=6379
Environment=SERVER_PORT=8080
Environment=JWT_SECRET=your-super-secret-jwt-key
Environment=ENVIRONMENT=production
Environment=LOG_LEVEL=info
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable codebase-golang
sudo systemctl start codebase-golang
sudo systemctl status codebase-golang
```

### 3. **Docker with Environment Variables**
```bash
# Run with environment variables
docker run -d \
  --name codebase-golang \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  -e REDIS_HOST=your-redis-host \
  codebase-golang:latest

# Or with environment file
docker run -d \
  --name codebase-golang \
  -p 8080:8080 \
  --env-file .env \
  codebase-golang:latest
```

### 4. **Kubernetes ConfigMap & Secrets**
```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "codebase_db"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
  SERVER_PORT: "8080"
  ENVIRONMENT: "production"
  LOG_LEVEL: "info"

---
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  DB_PASSWORD: "your-secure-password"
  JWT_SECRET: "your-super-secret-jwt-key"
```

### 5. **Cloud Platform Specific**

#### **AWS ECS/Fargate**
```json
{
  "taskDefinition": {
    "containerDefinitions": [
      {
        "environment": [
          {"name": "DB_HOST", "value": "your-rds-endpoint"},
          {"name": "REDIS_HOST", "value": "your-elasticache-endpoint"}
        ],
        "secrets": [
          {"name": "DB_PASSWORD", "valueFrom": "arn:aws:secretsmanager:..."}
        ]
      }
    ]
  }
}
```

#### **Google Cloud Run**
```bash
# Deploy with environment variables
gcloud run deploy codebase-golang \
  --image gcr.io/PROJECT_ID/codebase-golang \
  --set-env-vars="DB_HOST=your-db-host,REDIS_HOST=your-redis-host" \
  --set-secrets="DB_PASSWORD=db-password:latest"
```

#### **Azure Container Instances**
```bash
az container create \
  --resource-group myResourceGroup \
  --name codebase-golang \
  --image your-registry/codebase-golang:latest \
  --environment-variables 'DB_HOST'='your-db-host' 'REDIS_HOST'='your-redis-host' \
  --secure-environment-variables 'DB_PASSWORD'='your-password'
```

## Production Deployment

### 1. **Prepare for Production**
```bash
# Prepare optimized build
make deploy-prep

# Push to registry
make docker-push
```

### 2. **Security Best Practices**
```bash
# Run security scan
make security

# Use non-root user (already configured in Dockerfile)
# Set proper file permissions
chmod 600 .env
chown appuser:appgroup /opt/codebase-golang

# Use secrets management
# - AWS Secrets Manager
# - HashiCorp Vault
# - Kubernetes Secrets
```

### 3. **Database Setup**
```bash
# Run migrations
make migrate-up

# Set proper database permissions
# Create read-only user for analytics
# Setup backup strategies
```

### 4. **Monitoring Setup**
```bash
# Add monitoring endpoints
curl http://localhost:8080/health
curl http://localhost:8080/metrics

# Setup log aggregation
# - ELK Stack
# - Fluentd
# - Cloud native logging
```

## Environment Configuration Priority

Application reads configuration in this order:
1. **System Environment Variables** (highest priority)
2. **Docker environment variables**
3. **Systemd service environment**
4. **`.env` file** (lowest priority)

## Sample Production Environment Variables

```bash
# Database
DB_HOST=prod-postgres.example.com
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=super-secure-password-123
DB_NAME=codebase_prod

# Cache
REDIS_HOST=prod-redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=redis-secure-password

# Application
SERVER_PORT=8080
ENVIRONMENT=production
JWT_SECRET=your-256-bit-secret-key-for-jwt-tokens
LOG_LEVEL=info
GIN_MODE=release

# Security
ENABLE_HTTPS=true
SSL_CERT_PATH=/etc/ssl/certs/app.crt
SSL_KEY_PATH=/etc/ssl/private/app.key

# Performance
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
REDIS_POOL_SIZE=10
```

## Troubleshooting

### 1. **Check Environment Variables**
```bash
# View current environment
env | grep -E "(DB_|REDIS_|SERVER_)"

# Test with specific config
DB_HOST=localhost go run main.go
```

### 2. **Docker Issues**
```bash
# Check container environment
docker exec -it codebase-golang env

# View logs
docker logs codebase-golang

# Debug container
docker exec -it codebase-golang sh
```

### 3. **Systemd Service Issues**
```bash
# Check service status
sudo systemctl status codebase-golang

# View logs
sudo journalctl -u codebase-golang -f

# Restart service
sudo systemctl restart codebase-golang
```

## Scaling Considerations

1. **Horizontal Scaling**: Use load balancer with multiple instances
2. **Database Connection Pooling**: Configure `DB_MAX_OPEN_CONNS`
3. **Redis Clustering**: For high availability
4. **CDN**: For static assets
5. **Monitoring**: Set up alerts and dashboards