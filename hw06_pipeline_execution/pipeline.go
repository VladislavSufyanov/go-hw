package hw06pipelineexecution

type (
	In  = <-chan any
	Out = In
	Bi  = chan any
)

type Stage func(in In) (out Out)

func wrapStage(done In, in In, s Stage) (Out, <-chan struct{}) {
	ch := make(chan any)
	ready := make(chan struct{})
	go func() {
		defer close(ch)
		inS := make(chan any)
		outS := s(inS)
		defer func() {
			close(inS)
			<-outS
		}()
		close(ready)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case inS <- v:
					stageRes, ok := <-outS
					if !ok {
						return
					}
					select {
					case <-done:
						return
					default:
					}
					select {
					case ch <- stageRes:
					case <-done:
						return
					}
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return ch, ready
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := in
	var ready <-chan struct{}
	for _, s := range stages {
		ch, ready = wrapStage(done, ch, s)
		<-ready
	}
	return ch
}
