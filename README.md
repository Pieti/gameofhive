# Game of Hive

A version of [Conway's Game of life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life). On a hexgrid. With BEES!

Play it [here](https://pieti.github.io/gameofhive/)

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
