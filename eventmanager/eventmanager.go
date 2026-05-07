package eventmanager

type EventManager struct {
	handlers map[string][]chan SquelchEventData
}

func (ev *EventManager) AddHandler(e string, ch chan SquelchEventData) {
	if ev.handlers == nil {
		ev.handlers = make(map[string][]chan SquelchEventData)
	}

	ev.handlers[e] = append(ev.handlers[e], ch)
}

func (ev *EventManager) DelHandler(e string, ch chan SquelchEventData) {
	if _, ok := ev.handlers[e]; ok {
		for i := range ev.handlers[e] {
			if ev.handlers[e][i] == ch {
				ev.handlers[e] = append(ev.handlers[e][:i], ev.handlers[e][i+1:]...)
				break
			}
		}
	}
}

func (ev *EventManager) Emit(e string, data SquelchEventData) {
	if ev.handlers[e] != nil {
		for _, handler := range ev.handlers[e] {
			go func(handler chan SquelchEventData) {
				handler <- data
			}(handler)
		}
	}
}
