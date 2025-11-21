package main

// state machine for animations

type State interface {
	TickState()
	GetName() string
	ResetTime()
}

type StateMachine struct {
	CurrentState State
	StateMap     map[string]State
}

func NewStateMachine(initialState State) StateMachine {
	newMachine := StateMachine{
		CurrentState: initialState,
		StateMap:     make(map[string]State),
	}
	newMachine.AddState(initialState)
	return newMachine
}

func (sm *StateMachine) AddState(newState State) {
	sm.StateMap[newState.GetName()] = newState
}

func (sm *StateMachine) ChangeState(stateName string) {
	if sm.CurrentState.GetName() == stateName {
		return
	}
	sm.CurrentState = sm.StateMap[stateName]
	sm.CurrentState.ResetTime()

}

func (sm *StateMachine) Tick() {
	sm.CurrentState.TickState()
}
