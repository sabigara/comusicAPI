# comusicAPI

## Start developing

```bash
git clone https://github.com/sabigara/comusicAPI
cd comusicAPI
cp docker-compose.dev.sample.yaml docker-compose.dev.yaml
docker-compose -f docker-compose.dev.yaml build
docker-compose -f docker-compose.dev.yaml up -d
```

```bash
docker container exec -it comusic_api sh
# in the container
make migrate-up
make dev
```

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

Then, paste the `studio_id` on `const studioId=<studio_id>` in client's `/src/components/Browser.tsx` .
