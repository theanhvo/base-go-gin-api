#!/bin/bash

# BaseAPI Multi-Worker Deployment Script
# This script automates the deployment of multiple BaseAPI workers
# Author: DevOps Team
# Usage: ./deploy-workers.sh [start|stop|restart|status]

set -euo pipefail

# Configuration variables
APP_NAME="baseapi"
APP_USER="appuser"
APP_GROUP="appgroup"
APP_DIR="/opt/baseapi"
BINARY_NAME="codebase-golang"
SERVICE_PREFIX="baseapi-worker"
WORKER_COUNT=3
LOG_DIR="/var/log/baseapi"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if script is run as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Function to create application user and directories
setup_user_and_directories() {
    print_info "Setting up application user and directories..."
    
    # Create application user if not exists
    if ! id "$APP_USER" &>/dev/null; then
        useradd --system --home-dir "$APP_DIR" --shell /bin/false --create-home "$APP_USER"
        print_success "Created application user: $APP_USER"
    else
        print_info "Application user already exists: $APP_USER"
    fi
    
    # Create directories
    mkdir -p "$APP_DIR" "$LOG_DIR"
    chown -R "$APP_USER:$APP_GROUP" "$APP_DIR" "$LOG_DIR"
    chmod 755 "$APP_DIR" "$LOG_DIR"
    
    print_success "Directories setup completed"
}

# Function to build and deploy the application
deploy_application() {
    print_info "Building and deploying application..."
    
    # Navigate to project directory
    cd "$(dirname "$0")/../.."
    
    # Build the application
    print_info "Building Go application..."
    go build -ldflags="-s -w" -o "$BINARY_NAME" main.go
    
    # Copy binary to application directory
    cp "$BINARY_NAME" "$APP_DIR/"
    chown "$APP_USER:$APP_GROUP" "$APP_DIR/$BINARY_NAME"
    chmod +x "$APP_DIR/$BINARY_NAME"
    
    # Copy configuration files if they exist
    if [[ -f ".env.production" ]]; then
        cp .env.production "$APP_DIR/.env"
        chown "$APP_USER:$APP_GROUP" "$APP_DIR/.env"
        chmod 600 "$APP_DIR/.env"
        print_success "Production environment file copied"
    fi
    
    print_success "Application deployed successfully"
}

# Function to install systemd service files
install_services() {
    print_info "Installing systemd service files..."
    
    for ((i=1; i<=WORKER_COUNT; i++)); do
        SERVICE_NAME="${SERVICE_PREFIX}${i}"
        SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
        SOURCE_FILE="deployment/systemd/${SERVICE_NAME}.service"
        
        if [[ -f "$SOURCE_FILE" ]]; then
            cp "$SOURCE_FILE" "$SERVICE_FILE"
            print_success "Installed service: $SERVICE_NAME"
        else
            print_warning "Service file not found: $SOURCE_FILE"
        fi
    done
    
    # Reload systemd
    systemctl daemon-reload
    print_success "SystemD daemon reloaded"
}

# Function to start all workers
start_workers() {
    print_info "Starting all workers..."
    
    for ((i=1; i<=WORKER_COUNT; i++)); do
        SERVICE_NAME="${SERVICE_PREFIX}${i}"
        
        # Enable and start service
        systemctl enable "$SERVICE_NAME"
        systemctl start "$SERVICE_NAME"
        
        # Check status
        if systemctl is-active --quiet "$SERVICE_NAME"; then
            print_success "Worker $i started successfully (port 808$((i-1)))"
        else
            print_error "Failed to start worker $i"
            systemctl status "$SERVICE_NAME" --no-pager
        fi
    done
    
    # Wait a bit for services to stabilize
    sleep 5
    
    # Show overall status
    show_status
}

# Function to stop all workers
stop_workers() {
    print_info "Stopping all workers..."
    
    for ((i=1; i<=WORKER_COUNT; i++)); do
        SERVICE_NAME="${SERVICE_PREFIX}${i}"
        
        if systemctl is-active --quiet "$SERVICE_NAME"; then
            systemctl stop "$SERVICE_NAME"
            print_success "Worker $i stopped"
        else
            print_info "Worker $i was not running"
        fi
    done
}

# Function to restart all workers
restart_workers() {
    print_info "Restarting all workers..."
    stop_workers
    sleep 2
    start_workers
}

# Function to show status of all workers
show_status() {
    print_info "Worker Status:"
    echo "=============================================="
    
    for ((i=1; i<=WORKER_COUNT; i++)); do
        SERVICE_NAME="${SERVICE_PREFIX}${i}"
        PORT="808$((i-1))"
        
        if systemctl is-active --quiet "$SERVICE_NAME"; then
            STATUS="${GREEN}RUNNING${NC}"
            # Check if port is responding
            if curl -s "http://localhost:$PORT/health" >/dev/null 2>&1; then
                HEALTH="${GREEN}HEALTHY${NC}"
            else
                HEALTH="${RED}UNHEALTHY${NC}"
            fi
        else
            STATUS="${RED}STOPPED${NC}"
            HEALTH="${RED}DOWN${NC}"
        fi
        
        echo -e "Worker $i (Port $PORT): $STATUS | Health: $HEALTH"
    done
    
    echo "=============================================="
    
    # Show nginx status if installed
    if systemctl list-unit-files nginx.service >/dev/null 2>&1; then
        if systemctl is-active --quiet nginx; then
            echo -e "Nginx Load Balancer: ${GREEN}RUNNING${NC}"
        else
            echo -e "Nginx Load Balancer: ${RED}STOPPED${NC}"
        fi
    fi
}

# Function to setup nginx load balancer
setup_nginx() {
    print_info "Setting up Nginx load balancer..."
    
    # Install nginx if not present
    if ! command -v nginx &> /dev/null; then
        print_info "Installing Nginx..."
        apt-get update
        apt-get install -y nginx
    fi
    
    # Copy nginx configuration
    NGINX_CONFIG_SOURCE="deployment/nginx/load-balancer.conf"
    NGINX_CONFIG_DEST="/etc/nginx/sites-available/baseapi"
    
    if [[ -f "$NGINX_CONFIG_SOURCE" ]]; then
        cp "$NGINX_CONFIG_SOURCE" "$NGINX_CONFIG_DEST"
        
        # Enable the site
        ln -sf "$NGINX_CONFIG_DEST" "/etc/nginx/sites-enabled/baseapi"
        
        # Remove default site
        rm -f "/etc/nginx/sites-enabled/default"
        
        # Test nginx configuration
        if nginx -t; then
            systemctl enable nginx
            systemctl restart nginx
            print_success "Nginx load balancer configured and started"
        else
            print_error "Nginx configuration test failed"
            return 1
        fi
    else
        print_error "Nginx configuration file not found: $NGINX_CONFIG_SOURCE"
        return 1
    fi
}

# Function to show logs
show_logs() {
    local worker_num=${1:-1}
    SERVICE_NAME="${SERVICE_PREFIX}${worker_num}"
    
    print_info "Showing logs for $SERVICE_NAME..."
    journalctl -u "$SERVICE_NAME" -f --no-pager
}

# Function to run health checks
health_check() {
    print_info "Running health checks..."
    
    local failed=0
    
    for ((i=1; i<=WORKER_COUNT; i++)); do
        PORT="808$((i-1))"
        
        if curl -s --max-time 5 "http://localhost:$PORT/health" | grep -q "healthy"; then
            print_success "Worker $i health check passed (port $PORT)"
        else
            print_error "Worker $i health check failed (port $PORT)"
            ((failed++))
        fi
    done
    
    if [[ $failed -eq 0 ]]; then
        print_success "All health checks passed!"
        return 0
    else
        print_error "$failed workers failed health check"
        return 1
    fi
}

# Function to show help
show_help() {
    echo "BaseAPI Multi-Worker Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  setup     - Setup user, directories, and deploy application"
    echo "  start     - Start all workers"
    echo "  stop      - Stop all workers"
    echo "  restart   - Restart all workers"
    echo "  status    - Show status of all workers"
    echo "  logs [N]  - Show logs for worker N (default: 1)"
    echo "  health    - Run health checks on all workers"
    echo "  nginx     - Setup nginx load balancer"
    echo "  deploy    - Full deployment (setup + install + start)"
    echo "  help      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 deploy          # Full deployment"
    echo "  $0 restart         # Restart all workers"
    echo "  $0 logs 2          # Show logs for worker 2"
    echo "  $0 health          # Check all workers health"
}

# Main execution logic
main() {
    local command=${1:-help}
    
    case $command in
        "setup")
            check_root
            setup_user_and_directories
            deploy_application
            install_services
            ;;
        "start")
            check_root
            start_workers
            ;;
        "stop")
            check_root
            stop_workers
            ;;
        "restart")
            check_root
            restart_workers
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "${2:-1}"
            ;;
        "health")
            health_check
            ;;
        "nginx")
            check_root
            setup_nginx
            ;;
        "deploy")
            check_root
            setup_user_and_directories
            deploy_application
            install_services
            setup_nginx
            start_workers
            print_success "Full deployment completed!"
            print_info "Your API is now running on multiple workers behind nginx load balancer"
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# Execute main function with all arguments
main "$@"