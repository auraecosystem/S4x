docker build -t aura-ecosystem:fadk-s4x .
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t aura-ecosystem:fadk-s4x \
  --build-arg VERSION=V.FADK \
  --push .
mkdir -p config/s4xv scripts src
touch Dockerfile README.md .gitignore .env

docker buildx create --name aura-builder --use
docker buildx inspect --bootstrap
