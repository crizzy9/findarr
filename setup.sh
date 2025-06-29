#!/bin/bash

# Findarr setup and run script

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print banner
echo -e "${GREEN}"
echo "  ______ _           _                 "
echo " |  ____(_)         | |                "
echo " | |__   _ _ __   __| | __ _ _ __ _ __ "
echo " |  __| | | '_ \ / _\` |/ _\` | '__| '__|"
echo " | |    | | | | | (_| | (_| | |  | |   "
echo " |_|    |_|_| |_|\__,_|\__,_|_|  |_|   "
echo -e "${NC}"
echo "Universal Media Metadata Resolver & Organizer"
echo "============================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

# Create config directory if it doesn't exist
if [ ! -d "config" ]; then
    echo -e "${YELLOW}Creating config directory...${NC}"
    mkdir -p config
fi

# Function to show help
show_help() {
    echo "Usage: ./setup.sh [OPTION]"
    echo "Options:"
    echo "  start       Start Findarr (default if no option provided)"
    echo "  stop        Stop Findarr"
    echo "  restart     Restart Findarr"
    echo "  logs        Show Findarr logs"
    echo "  status      Show Findarr status"
    echo "  build       Build Findarr Docker image"
    echo "  help        Show this help message"
}

# Process command line arguments
case "$1" in
    start|"")
        echo -e "${GREEN}Starting Findarr...${NC}"
        docker-compose up -d
        echo -e "${GREEN}Findarr is now running at http://localhost:8080${NC}"
        ;;
    stop)
        echo -e "${YELLOW}Stopping Findarr...${NC}"
        docker-compose down
        echo -e "${GREEN}Findarr has been stopped.${NC}"
        ;;
    restart)
        echo -e "${YELLOW}Restarting Findarr...${NC}"
        docker-compose down
        docker-compose up -d
        echo -e "${GREEN}Findarr has been restarted and is running at http://localhost:8080${NC}"
        ;;
    logs)
        echo -e "${YELLOW}Showing Findarr logs (press Ctrl+C to exit)...${NC}"
        docker-compose logs -f
        ;;
    status)
        echo -e "${YELLOW}Findarr container status:${NC}"
        docker-compose ps
        ;;
    build)
        echo -e "${YELLOW}Building Findarr Docker image...${NC}"
        docker-compose build
        echo -e "${GREEN}Findarr Docker image has been built.${NC}"
        ;;
    help)
        show_help
        ;;
    *)
        echo -e "${RED}Invalid option: $1${NC}"
        show_help
        exit 1
        ;;
esac

exit 0

