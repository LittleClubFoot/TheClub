# ============================================================================
# The Club Dashboard - Build and Development Automation
# ============================================================================
#
# This Makefile provides convenient commands for building, testing, and
# deploying The Club dashboard. It supports both local development and
# containerized deployment workflows.
#
# Quick Start:
#   make help        - Show available commands
#   make docker-dev  - Start development environment
#   make docker-prod - Deploy production environment
#
# ============================================================================

.PHONY: help build run test clean dev docker-dev docker-prod docker-stop deps fmt lint init

# Default target - show help
help:
	@echo "🏴‍☠️ The Club Dashboard - Available Commands"
	@echo ""
	@echo "📦 Build & Run:"
	@echo "  build        Build the Go application binary"
	@echo "  run          Build and run the application locally"
	@echo "  test         Run all tests"
	@echo "  clean        Clean build artifacts and temporary files"
	@echo ""
	@echo "🔧 Development:"
	@echo "  dev          Start Go server with live reload (Air)"
	@echo "  docker-dev   Start full development stack with Docker"
	@echo "  fmt          Format Go code"
	@echo "  lint         Run code linting (requires golangci-lint)"
	@echo ""
	@echo "🚀 Deployment:"
	@echo "  docker-prod  Deploy production environment"
	@echo "  docker-stop  Stop all Docker services"
	@echo ""
	@echo "🛠️  Maintenance:"
	@echo "  deps         Install/update dependencies"
	@echo "  init         Initialize project configuration"
	@echo "  backup       Run configuration backup"
	@echo ""
	@echo "🌟 Enhancement Stacks (FOSS Community Tools):"
	@echo "  stack-monitoring      Deploy Uptime Kuma monitoring"
	@echo "  stack-security        Deploy WireGuard + Authelia SSO"
	@echo "  stack-observability   Deploy Prometheus + Grafana"
	@echo "  stack-automation      Deploy Kopia + Watchtower + Dockge"
	@echo "  stack-all             Deploy all enhancement stacks"
	@echo "  stack-stop            Stop all enhancement stacks"
	@echo ""
	@echo "📖 Documentation: See docs/ directory for detailed guides"

# ============================================================================
# BUILD TARGETS
# ============================================================================

# Build the Go application binary
build:
	@echo "🔨 Building The Club Go server..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/theclub-server ./cmd/webserver
	@echo "✅ Build complete: bin/theclub-server"

# Build and run the application locally
run: build
	@echo "🚀 Starting The Club server..."
	./bin/theclub-server

# Run all tests
test:
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@if [ -f coverage.out ]; then \
		echo "📊 Test coverage:"; \
		go tool cover -func=coverage.out | tail -1; \
	fi

# Clean build artifacts and temporary files
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/ tmp/ coverage.out
	go clean -cache -testcache -modcache
	@echo "✅ Clean complete"

# ============================================================================
# DEVELOPMENT TARGETS
# ============================================================================

# Start Go server with live reload using Air
dev:
	@echo "🔥 Starting development server with live reload..."
	@command -v air >/dev/null 2>&1 || { \
		echo "❌ Air not found. Installing..."; \
		go install github.com/air-verse/air@latest; \
	}
	@echo "🌐 Server will be available at: http://localhost:8080/test"
	air

# Start full development environment with Docker
docker-dev:
	@echo "🐳 Starting development environment..."
	@command -v docker >/dev/null 2>&1 || { \
		echo "❌ Docker not found. Please install Docker first."; \
		exit 1; \
	}
	docker compose -f docker-compose.dev.yml up --build
	@echo "🌐 Dashboard available at: https://localhost"

# Format Go code
fmt:
	@echo "🎨 Formatting Go code..."
	go fmt ./...
	@echo "✅ Code formatting complete"

# Run code linting
lint:
	@echo "🔍 Running code linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "✅ Linting complete"; \
	else \
		echo "⚠️  golangci-lint not installed."; \
		echo "   Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# ============================================================================
# DEPLOYMENT TARGETS
# ============================================================================

# Deploy production environment
docker-prod:
	@echo "🚀 Deploying production environment..."
	@command -v docker >/dev/null 2>&1 || { \
		echo "❌ Docker not found. Please install Docker first."; \
		exit 1; \
	}
	docker compose -f docker-compose.prod.yml up -d --build
	@echo "✅ Production deployment complete"
	@echo "🌐 Dashboard available at: https://yourdomain.com"

# Stop all Docker services
docker-stop:
	@echo "🛑 Stopping Docker services..."
	@docker compose -f docker-compose.dev.yml down 2>/dev/null || true
	@docker compose -f docker-compose.prod.yml down 2>/dev/null || true
	@echo "✅ All services stopped"

# ============================================================================
# MAINTENANCE TARGETS
# ============================================================================

# Install and update dependencies
deps:
	@echo "📦 Installing/updating dependencies..."
	go mod tidy
	go mod download
	@echo "🔥 Installing Air for live reload..."
	go install github.com/air-verse/air@latest
	@echo "✅ Dependencies updated"

# Initialize project configuration
init:
	@echo "🔧 Initializing project configuration..."
	@mkdir -p logs assets/icons
	@echo "✅ Created runtime directories"
	@echo ""
	@echo "📖 Configuration Setup:"
	@echo "   • Homer config: config/homer.yml (edit this file directly)"
	@echo "   • Changes are automatically picked up by Docker"
	@echo ""
	@echo "📖 Next steps:"
	@echo "   1. Edit Caddyfile with your domain"
	@echo "   2. Edit config/homer.yml with your services"
	@echo "   3. Run 'make docker-dev' to start development"

# ============================================================================
# UTILITY TARGETS
# ============================================================================

# Show project structure
structure:
	@echo "📁 Project Structure:"
	@if command -v tree >/dev/null 2>&1; then \
		tree -I 'node_modules|.git|tmp|bin|coverage.out' -a; \
	else \
		find . -type f -not -path './.git/*' -not -path './tmp/*' -not -path './bin/*' | sort; \
	fi

# Show Docker service status
status:
	@echo "🐳 Docker Service Status:"
	@docker compose -f docker-compose.dev.yml ps 2>/dev/null || echo "Development stack not running"
	@docker compose -f docker-compose.prod.yml ps 2>/dev/null || echo "Production stack not running"

# View logs from Docker services
logs:
	@echo "📋 Recent Docker logs:"
	@docker compose -f docker-compose.dev.yml logs --tail=50 2>/dev/null || \
	 docker compose -f docker-compose.prod.yml logs --tail=50 2>/dev/null || \
	 echo "No Docker services running"

# Health check all services
health:
	@echo "🏥 Health Check:"
	@curl -k -s https://localhost/health 2>/dev/null && echo "✅ Main service healthy" || echo "❌ Main service unhealthy"
	@curl -k -s https://localhost/test 2>/dev/null >/dev/null && echo "✅ Go server healthy" || echo "❌ Go server unhealthy"

# ============================================================================
# ENHANCEMENT STACKS (Phase 1-4)
# ============================================================================

# Deploy monitoring stack (Uptime Kuma)
stack-monitoring:
	@echo "📊 Deploying monitoring stack..."
	docker compose -f docker-compose.monitoring.yml up -d
	@echo "✅ Monitoring stack deployed"
	@echo "   Access Uptime Kuma at: https://localhost/status"

# Deploy security stack (WireGuard + Authelia)
stack-security:
	@echo "🔐 Deploying security stack..."
	@if [ ! -f ./authelia/configuration.yml ]; then \
		echo "⚠️  Warning: Authelia not configured. See authelia/configuration.yml"; \
	fi
	docker compose -f docker-compose.security.yml up -d
	@echo "✅ Security stack deployed"
	@echo "   Authelia: https://localhost/auth"
	@echo "   WireGuard configs: docker exec wireguard cat /config/peer1/peer1.conf"

# Deploy observability stack (Prometheus + Grafana)
stack-observability:
	@echo "📈 Deploying observability stack..."
	docker compose -f docker-compose.observability.yml up -d
	@echo "✅ Observability stack deployed"
	@echo "   Grafana: https://localhost/grafana (admin/changeme)"

# Deploy automation stack (Kopia + Watchtower + Dockge)
stack-automation:
	@echo "🤖 Deploying automation stack..."
	@mkdir -p /opt/stacks
	docker compose -f docker-compose.automation.yml up -d
	@echo "✅ Automation stack deployed"
	@echo "   Dockge: https://localhost/dockge"
	@echo "   Kopia: https://localhost/backup"

# Deploy all enhancement stacks
stack-all: stack-monitoring stack-security stack-observability stack-automation
	@echo ""
	@echo "🎉 All enhancement stacks deployed!"
	@echo ""
	@echo "📖 Access URLs:"
	@echo "   Dashboard:  https://localhost/"
	@echo "   Monitoring: https://localhost/status"
	@echo "   Metrics:    https://localhost/grafana"
	@echo "   Auth:       https://localhost/auth"
	@echo "   Management: https://localhost/dockge"
	@echo "   Backup:     https://localhost/backup"

# Stop all enhancement stacks
stack-stop:
	@echo "🛑 Stopping enhancement stacks..."
	@docker compose -f docker-compose.monitoring.yml down 2>/dev/null || true
	@docker compose -f docker-compose.security.yml down 2>/dev/null || true
	@docker compose -f docker-compose.observability.yml down 2>/dev/null || true
	@docker compose -f docker-compose.automation.yml down 2>/dev/null || true
	@echo "✅ All enhancement stacks stopped"

# Run configuration backup
backup:
	@echo "💾 Running configuration backup..."
	@./scripts/config-backup.sh
	@echo "✅ Backup complete"

# Update Makefile help
.PHONY: stack-monitoring stack-security stack-observability stack-automation stack-all stack-stop backup