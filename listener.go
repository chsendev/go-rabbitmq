package rmq

type Listener interface {
	Use(...HandlerFunc) Listener
	Listen(string, ...HandlerFunc) Listener
	Group(...HandlerFunc) Listener
	Error() error
}

type RouterGroup struct {
	*Engine
	Handlers HandlersChain
}

func (r *RouterGroup) Group(handlers ...HandlerFunc) Listener {
	return &RouterGroup{
		Handlers: r.combineHandlers(handlers),
		Engine:   r.Engine,
	}
}

func (r *RouterGroup) Listen(queue string, handlers ...HandlerFunc) Listener {
	return r.listen(queue, handlers)
}

func (r *RouterGroup) Use(handlers ...HandlerFunc) Listener {
	r.Handlers = append(r.Handlers, handlers...)
	return r
}

func (r *RouterGroup) Error() error {
	return r.Engine.Error()
}

func (r *RouterGroup) listen(queue string, handlers HandlersChain) Listener {
	handlers = r.combineHandlers(handlers)
	r.Engine.addHandlers(queue, handlers)
	return r
}

func (r *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	return append(r.Handlers, handlers...)
}
