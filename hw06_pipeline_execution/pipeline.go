package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for i := 0; i < len(stages); i++ {
		select {
		case <-done:
			return in
		default:
			in = stages[i](in)
		}
	}

	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if ok == false {
					return
				}
				out <- v
			}
		}
	}()
	return out
}
