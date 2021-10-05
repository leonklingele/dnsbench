package dnsbench

type Result struct {
	WorkItem     workItem
	Measurements Measurements
}

type ResultChan <-chan Result
