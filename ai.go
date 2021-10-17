package main

type AIHandler struct {
}

func NewAIHandler() *AIHandler {
	return &AIHandler{}
}

func (a *AIHandler) Read() Input {
	return KeyDefaultOrGraveyard
}

func (a *AIHandler) Cancel() {

}
