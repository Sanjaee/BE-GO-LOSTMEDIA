#!/bin/bash

# Docker build and run script

echo "Building Docker images..."
docker-compose build

echo "Starting services..."
docker-compose up -d

echo "Waiting for services to be ready..."
sleep 10

echo "Checking service status..."
docker-compose ps

echo ""
echo "Services are running!"
echo "API: http://localhost:5000"
echo "RabbitMQ Management: http://localhost:15672"
echo "Meilisearch: http://localhost:7700"
echo ""
echo "To view logs: docker-compose logs -f app"
echo "To stop: docker-compose down"

