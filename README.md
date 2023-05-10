# sego
Search Engine written in Go.

This engine will index the linux API documentation stored in `linux-docs` folder using the `TF-IDF` method.
It can also serve them through an API.

## Documentation
[Wikipedia](https://en.wikipedia.org/wiki/Tf%E2%80%93idf)

:notebook:
- For **Term Frequency**, we use the `raw count weighting scheme`.
- For **Inverse document Frequency**, we use the `inverse document frequency smooth weighting scheme`.

## Run
- Index files:
```shell
go run main.go -index
```

- Serve files:
```shell
go run main.go -serve
```

- Query the server:
```shell
curl 'localhost:4000/search?query=memory%20management'
```

## Frontend
```shell
cd ui
npm install
npm run dev
```

## Inner workings
- Index: parse the .html docs into a json that maps, for each document every word occurrence inside it.
- Serve: load the json file and apply `TF-IDF` algorithm to the search terms.

## TODO
- enable debug logs
- try changing representation format to a more performant one
- docker/docker-compose

## Indexed files
We will index the linux kernel documentation. We have obtained this docs from the linux repo:
```bash
git clone --depth 1 https://github.com/torvalds/linux.git
cd linux
make htmldocs
```

Now, inside `Documentation/output`, there will be all the docs in `.html` format.
