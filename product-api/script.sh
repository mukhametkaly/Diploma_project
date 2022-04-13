#todo: change hosts and name after deployment
docker build -t registry.dar.kz/forte-kassa/fortekassa-shopping-cart-api .
docker tag registry.dar.kz/forte-kassa/fortekassa-shopping-cart-api:latest registry.dar.kz/forte-kassa/fortekassa-shopping-cart-api:v1.0.12
docker push registry.dar.kz/forte-kassa/fortekassa-shopping-cart-api:v1.0.12

