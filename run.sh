# Build the image
docker build -t golang-hexagonal:latest .

# Run the container
docker run -p 8080:8080 --env-file .env golang-hexagonal:latest

# Or with individual environment variables
docker run -p 8080:8080 \
  -e DATABASE_URL="your_db_url" \
  -e REDIS_URL="your_redis_url" \
  golang-hexagonal:latest