# sego
Search Engine written in Go.

This engine will index the docs from `linux-docs` folder using the TF-IDF method and serve them.

## Documentation
[Wikipedia](https://en.wikipedia.org/wiki/Tf%E2%80%93idf)

:notebook:
- For Term frequency, we use the `raw count weighting scheme`.
- For Inverse document frequency, we use the `inverse document frequency smooth weighting scheme`.

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
curl 'localhost:4000/search?query="memory%20management"'
```

## Frontend
```shell
cd ui
npm run dev
```

## Inner workings
- recorrer todos los ficheros .html
- parsearlos y construir un json que mapee cada palabra a los docs donde aparece, y cuantas veces aparece en cada doc
- servir html para hacer busquedas

## TODO
- UI in svelte
- enable debug logs
- probar libreria de fast encoding con json
- probar el cambio de representacion en disco de json a gob, protobuf y flatbuffers

## Indexed files
We will index the linux kernel documentation. We have obtained this docs from the linux repo:
```bash
git clone --depth 1 https://github.com/torvalds/linux.git
cd linux
make htmldocs
```

Now, inside `Documentation/output`, there will be all the docs in `.html` format.
We will only index the docs from `core-api` to speed up.
