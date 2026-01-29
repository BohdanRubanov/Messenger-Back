```bash
DOCKER
=======================================================================
docker compose down          # stop container
docker volume rm lessonproj_postgres_data  # clear volume
docker compose up -d         # start container

# Connect to PostgreSQL database 'lesson' as user 'lesson_user' inside the container
docker exec -it shop_postgres psql -U lesson_user -d lesson 
=======================================================================
```

