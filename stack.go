package main

import "github.com/hajimehoshi/ebiten/v2"

type Stack struct {
	cards []Card
	op    *ebiten.DrawImageOptions
}

func (s *Stack) Put(c Card) {
	s.cards = append(s.cards, c)
}

func (s *Stack) Remove() Card {
	if len(s.cards) == 0 {
		return Card{}
	}
	card := s.Top()
	s.cards = s.cards[:len(s.cards)-1]
	return card
}

func (s *Stack) Top() Card {
	if len(s.cards) == 0 {
		return Card{}
	}
	card := s.cards[len(s.cards)-1]
	return card
}

func (s *Stack) All() []Card {
	return s.cards
}

type BattleStack struct {
	Stack
	op *ebiten.DrawImageOptions
}

func NewBattleStack(op ebiten.DrawImageOptions) *BattleStack {
	op.GeoM.Translate(config.Layout.Battle.StartX, config.Layout.Battle.StartY)
	return &BattleStack{op: &op}
}

func (bs *BattleStack) Status() Status {
	return bs.Top().Effects.Status
}

func (bs *BattleStack) Draw(target *ebiten.Image) {
	if bs.Top().image != nil {
		target.DrawImage(bs.Top().image, bs.op)
	}
}

type TerrainStack struct {
	Stack
	op *ebiten.DrawImageOptions
}

func NewTerrainStack(op ebiten.DrawImageOptions) *TerrainStack {
	op.GeoM.Translate(config.Layout.Terrain.StartX, config.Layout.Terrain.StartY)
	return &TerrainStack{op: &op}
}

func (ts *TerrainStack) Terrain() string {
	return ts.Top().Effects.Terrain
}

func (ts *TerrainStack) Clear() {
	ts.cards = nil
}

func (ts *TerrainStack) Draw(target *ebiten.Image) {
	if ts.Top().image != nil {
		target.DrawImage(ts.Top().image, ts.op)
	}
}

type DefenseStack struct {
	Stack
	start ebiten.DrawImageOptions
}

func NewDefenseStack(op ebiten.DrawImageOptions) *DefenseStack {
	return &DefenseStack{start: op}
}

func (ds *DefenseStack) Immunities() []Status {
	var im []Status
	for _, c := range ds.cards {
		im = append(im, c.Effects.Immunity)
	}
	return im
}

func (ds *DefenseStack) Draw(target *ebiten.Image) {

	op := ds.start
	op.GeoM.Translate(config.Layout.Defense.StartX, config.Layout.Defense.StartY)

	for _, c := range ds.cards {
		target.DrawImage(c.image, &op)
		op.GeoM.Translate(config.Layout.Card.Width, 0)
	}
}

type Graveyard struct {
	Stack
	op *ebiten.DrawImageOptions
}

func NewGraveyard() *Graveyard {
	var g Graveyard
	g.op = &ebiten.DrawImageOptions{}
	g.op.GeoM.Scale(.10, .10)
	g.op.GeoM.Translate(config.Layout.System.StartX+config.Layout.Graveyard.StartX, config.Layout.System.StartY+config.Layout.Graveyard.StartY)
	return &g
}

func (g *Graveyard) Draw(target *ebiten.Image) {
	if g.Top().image != nil {
		target.DrawImage(g.Top().image, g.op)
	}
}

type TravelStack struct {
	Stack
	start ebiten.DrawImageOptions
}

func NewTravelStack(op ebiten.DrawImageOptions) *TravelStack {
	return &TravelStack{start: op}
}

func (ts *TravelStack) Draw(target *ebiten.Image) {
	op := ts.start
	op.GeoM.Translate(config.Layout.Travel.StartX, config.Layout.Travel.StartY)

	for _, c := range ts.cards {
		target.DrawImage(c.image, &op)
		op.GeoM.Translate(20, 0)
	}
}
