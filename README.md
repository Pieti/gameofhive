# Game of Hive

Game of life, on a hexgrid, with BEES


## Requirements

* [Go](https://go.dev/)
* [Ebiten](https://ebitengine.org/en/documents/install.html)


## Run the game

Compile into WebAssembly

```
go install github.com/hajimehoshi/wasmserve@latest
wasmserve .
```

Browse to http://localhost:8080/ and you should see a game!


Alternatively it can run natively
```
go run .
```
