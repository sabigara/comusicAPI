# comusicAPI

## Start developing

### Set up files

```bash
git clone https://github.com/sabigara/comusicAPI
cd comusicAPI
cp docker-compose.dev.sample.yaml docker-compose.dev.yaml
```

Then, configure environment variables below in docker-compose.dev.yaml following instruction comments.

* GOOGLE_APPLICATION_CREDENTIALS
* SENDGRID_API_KEY

### Create docker containers

```bash
docker-compose -f docker-compose.dev.yaml build
docker-compose -f docker-compose.dev.yaml up -d
```

But `make dev` command specified in `docker-compose.dev.yaml` fails for the first time, because DB tables are not yet created. 

So next process is necessary.

### Migrate DB

```bash
docker container exec -it comusic_api sh # login to the container
make migrate-up # create tables
make dev # restart server process
```

### Supply initial data

```bash
curl --request POST \
  --url http://localhost:1323/studios \
  --header 'content-type: application/json' \
  --data '{
        "name": "Studio One"
}'

# Response:
# {
#  "id": "e4bbc729-5034-4f3b-985b-e15391670ce4",
#  "createdAt": "2020-02-29T07:18:24.5240884Z",
#  "updatedAt": "2020-02-29T07:18:24.5241021Z",
#  "ownerId": "4148e7cc-a5f0-4fb4-9392-ee82f0e324d1",
#  "name": "Studio One"
# }
```

Use `id` in the response as `studio_id` for following request:

```bash
curl --request POST \
  --url 'http://localhost:1323/songs?studio_id=<id_from_previous_requet_response>' \
  --header 'content-type: application/json' \
  --data '{
        "name": "Song 1"
}'
```

### Modify source

Then, paste the `studio_id` on `const studioId=<studio_id>` in client's `/src/components/Browser.tsx` .
