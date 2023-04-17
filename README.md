# sego
Search Engine written in Go.

This engine will index the docs from `linux-core-api-docs` folder and serve them.

## Run
- Index files:
```shell
go run main.go -index
```

- Serve files:
```shell
go run main.go -serve
```

## Inner workings
- recorrer todos los ficheros .html
- parsearlos y construir un json que mapee cada palabra a los docs donde aparece, y cuantas veces aparece en cada doc
- servir html para hacer busquedas

## Mejoras
- cambiar la representacion en disco de json a protobuf/flatbuffers

## Indexed files
We will index a subset of the linux kernel documentation, the `core-api` to be precise. We have obtained this docs from the linux repo:
```bash
git clone --depth 1 https://github.com/torvalds/linux.git
cd linux
make htmldocs
```

Now, you will have a folder Documentation/output, with all the docs in `.html` format.
We will only index the docs from `core-api` to speed up.
