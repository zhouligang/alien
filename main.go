package main

//go:generate go install github.com/hajimehoshi/file2byteslice
//go:generate mkdir resources
//go:generate file2byteslice -input ../images/ship.png -output resources/ship.go -package resources -var ShipPng
//go:generate file2byteslice -input ../images/alien.png -output resources/alien.go -package resources -var AlienPng
//go:generate file2byteslice -input config.json -output resources/config.go -package resources -var ConfigJson

import (
	"games/alien"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game := alien.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	//if err := http.ListenAndServe(":8080", http.FileServer(http.Dir("."))); err != nil {
	//	log.Fatal(err)
	//}
}
