package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
	mplusBigFace    *text.GoTextFace
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
	mplusNormalFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   15,
	}
	mplusBigFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   32,
	}

}

func DrawCenteredText(screen *ebiten.Image, s string, font *text.GoTextFace, cx, cy int, color color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(cx), float64(cy))
	op.LineSpacing = font.Size * 1.5
	text.Draw(screen, s, font, op)
}
