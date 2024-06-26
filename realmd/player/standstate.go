package player

type StandState uint8

const (
	StateStand          StandState = 0
	StateSit            StandState = 1
	StateSitChair       StandState = 2
	StateSleep          StandState = 3
	StateSitLowChair    StandState = 4
	StateSitMediumChair StandState = 5
	StateSitHighChair   StandState = 6
	StateDead           StandState = 7
	StateKneel          StandState = 8
	StateSubmerged      StandState = 9
)
