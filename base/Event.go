package base


type EventContainer interface {
	GetEvent() EventInterface
}

type EventInterface interface {
	Trigger(string, EventContainer)
	On(string, func(EventInterface, EventContainer))
	GetData() map[string]interface{}
	SetData(map[string]interface{})
}

type Event struct {
	obverseList map[string][]func(EventInterface, EventContainer)
	data map[string]interface{}
}

func (event *Event) Trigger(name string, object EventContainer)  {
	if _, ok := event.obverseList[name]; !ok {
		return
	}

	for _, handler := range event.obverseList[name] {
		handler(event, object)
	}
}

func (event *Event) On(name string, handler func(EventInterface, EventContainer))  {
	if event.obverseList == nil {
		event.obverseList = make(map[string][]func(EventInterface, EventContainer))
	}

	event.obverseList[name] = append(event.obverseList[name], handler)
}

func (event *Event) GetData() map[string]interface{}  {
	return event.data
}

func (event *Event) SetData(data map[string]interface{}) {
	event.data = data
}