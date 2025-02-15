package example

import "gameengine/src/engine/state"

type Player struct {
	stateMachine *state.StateMachine
}

func (p *Player) IsMoving() bool { return false }
func (p *Player) IsJumpPressed() bool { return false }
func (p *Player) IsGrounded() bool { return true }
func (p *Player) SetAnimation(name string) {}
func (p *Player) ApplyJumpForce() {} 