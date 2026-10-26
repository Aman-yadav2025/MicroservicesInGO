spin up whole application
    docker compose up -d 
tear down the main application container
    docker compose down -v
exce sql cammand in postgres running on docker
    docker exec -it video_plateform_db psql -U postgres -d auth_db