package infra

// EventGroupProcessor determines which EventHandler should handle event received from event bus.
// Compared to EventProcessor, EventGroupProcessor allows to have multiple handlers that share the same subscriber instance.
type EventGroupProcessor struct {
	groupEventHandlers []GroupEventHandler
}

// NewEventGroupProcessorWithConfig creates a new EventGroupProcessor.
func NewEventGroupProcessor(handlers ...GroupEventHandler) *EventGroupProcessor {
	return &EventGroupProcessor{
		groupEventHandlers: handlers,
	}
}

// AddHandlersGroup adds a new list of GroupEventHandler to the EventGroupProcessor and adds it to the router.
//
// Compared to AddHandlers, AddHandlersGroup allows to have multiple handlers that share the same subscriber instance.
//
// It's allowed to have multiple handlers for the same event type in one group, but we recommend to not do that.
// Please keep in mind that those handlers will be processed within the same message.
// If first handler succeeds and the second fails, the message will be re-delivered and the first will be re-executed.
//
// Handlers group needs to be unique within the EventProcessor instance.
//
// Handler group name is used as handler's name in router.
func (p *EventGroupProcessor) AddHandlersGroup(groupName string, handlers ...GroupEventHandler) error {
	//if len(handlers) == 0 {
	//	return errors.New("no handlers provided")
	//}
	//if _, ok := p.groupEventHandlers[groupName]; ok {
	//	return fmt.Errorf("event handler group '%s' already exists", groupName)
	//}
	//
	//p.groupEventHandlers[groupName] = handlers

	return nil
}
