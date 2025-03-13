package main

type config struct {
	PreviousId int
	CurrentId  int
	NextId     int
}

func initConfig() config {
	return config{
		NextId:     1,
		CurrentId:  0,
		PreviousId: 0,
	}
}
