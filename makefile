# Inicia todos os containers
up:
	docker-compose up

# Para e remove os containers e volumes (limpeza total)
down:
	docker-compose down -v

# Força a reconstrução das imagens
build:
	docker-compose build --no-cache

# Mostra os logs do container em tempo real
logs:
	docker-compose logs -f

# Limpa o cache do Docker
prune:
	docker builder prune -a -f