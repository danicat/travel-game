package main

import "github.com/hajimehoshi/ebiten/v2"

type Sprite interface {
	Draw(target *ebiten.Image)
}

type Scene struct {
	sprites []Sprite
}

func (s *Scene) AddSprite(sprite Sprite) {
	s.sprites = append(s.sprites, sprite)
}

func (s *Scene) Sprites() []Sprite {
	return s.sprites
}
