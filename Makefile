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

.PHONY: help build run test clean dev docker-dev docker-prod docker-stop deps fmt lint init gitea-setup gitea-logs ci-logs

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
	@echo "🔧 DevOps (Gitea + Woodpecker CI):"
	@echo "  gitea-setup  Generate secrets and show setup instructions"
	@echo "  gitea-logs   View Gitea service logs"
	@echo "  ci-logs      View Woodpecker CI logs (server + agent)"
	@echo ""
	@echo "🛠️  Maintenance:"
	@echo "  deps         Install/update dependencies"
	@echo "  init         Initialize project configuration"
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

# ============================================================================
# DEVOPS TARGETS (Gitea + Woodpecker CI)
# ============================================================================

# Generate secrets and show Gitea/Woodpecker setup instructions
gitea-setup:
	@echo "🔧 Gitea + Woodpecker CI Setup"
	@echo ""
	@if [ ! -f .env ]; then \
		echo "📋 Creating .env from .env.example..."; \
		cp .env.example .env; \
		AGENT_SECRET=$$(openssl rand -hex 32); \
		sed -i "s/^WOODPECKER_AGENT_SECRET=.*/WOODPECKER_AGENT_SECRET=$$AGENT_SECRET/" .env; \
		echo "✅ Generated WOODPECKER_AGENT_SECRET"; \
	else \
		echo "⚠️  .env already exists, skipping generation"; \
	fi
	@echo ""
	@echo "📖 Next steps:"
	@echo "  1. Start the stack:  make docker-dev"
	@echo "  2. Open Gitea:       https://localhost/gitea"
	@echo "  3. Create an admin account in Gitea"
	@echo "  4. In Gitea: Site Administration > Applications > Create OAuth2 App"
	@echo "     - Name: Woodpecker CI"
	@echo "     - Redirect URI: https://localhost/ci/authorize"
	@echo "  5. Copy Client ID and Secret to .env:"
	@echo "     - WOODPECKER_GITEA_CLIENT=<client-id>"
	@echo "     - WOODPECKER_GITEA_SECRET=<client-secret>"
	@echo "  6. Restart:  make docker-stop && make docker-dev"
	@echo "  7. Open Woodpecker:  https://localhost/ci"
	@echo ""
	@echo "📖 Full guide: docs/GITEA_WOODPECKER.md"

# View Gitea logs
gitea-logs:
	@echo "📋 Gitea Logs:"
	@docker logs theclub-gitea-dev --tail=50 2>/dev/null || \
	 docker logs theclub-gitea-prod --tail=50 2>/dev/null || \
	 echo "Gitea is not running"

# View Woodpecker CI logs (server + agent)
ci-logs:
	@echo "📋 Woodpecker Server Logs:"
	@docker logs theclub-woodpecker-server-dev --tail=30 2>/dev/null || \
	 docker logs theclub-woodpecker-server-prod --tail=30 2>/dev/null || \
	 echo "Woodpecker server is not running"
	@echo ""
	@echo "📋 Woodpecker Agent Logs:"
	@docker logs theclub-woodpecker-agent-dev --tail=30 2>/dev/null || \
	 docker logs theclub-woodpecker-agent-prod --tail=30 2>/dev/null || \
	 echo "Woodpecker agent is not running"

# Health check all services
health:
	@echo "🏥 Health Check:"
	@curl -k -s https://localhost/health 2>/dev/null && echo "✅ Main service healthy" || echo "❌ Main service unhealthy"
	@curl -k -s https://localhost/test 2>/dev/null >/dev/null && echo "✅ Go server healthy" || echo "❌ Go server unhealthy"
	@curl -k -s https://localhost/gitea/ 2>/dev/null >/dev/null && echo "✅ Gitea healthy" || echo "❌ Gitea unhealthy"
	@curl -k -s https://localhost/ci/healthz 2>/dev/null >/dev/null && echo "✅ Woodpecker CI healthy" || echo "❌ Woodpecker CI unhealthy"