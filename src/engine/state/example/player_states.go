package example

import (
	"gameengine/src/engine/state"
)

// プレイヤーの状態
type PlayerState struct {
	*state.BaseState
	player *Player
}

// アイドル状態
type IdleState struct {
	PlayerState
}

func NewIdleState(player *Player) *IdleState {
	return &IdleState{
		PlayerState: PlayerState{
			BaseState: state.NewBaseState("idle"),
			player:    player,
		},
	}
}

func (s *IdleState) OnEnter(data interface{}) {
	s.player.SetAnimation("idle")
}

func (s *IdleState) OnUpdate(dt float64) error {
	// 移動入力があれば走り状態へ
	if s.player.IsMoving() {
		return s.player.stateMachine.ChangeState("run", nil)
	}
	// ジャンプ入力があればジャンプ状態へ
	if s.player.IsJumpPressed() {
		return s.player.stateMachine.ChangeState("jump", nil)
	}
	return nil
}

// 走り状態
type RunState struct {
	PlayerState
}

func NewRunState(player *Player) *RunState {
	return &RunState{
		PlayerState: PlayerState{
			BaseState: state.NewBaseState("run"),
			player:    player,
		},
	}
}

func (s *RunState) OnEnter(data interface{}) {
	s.player.SetAnimation("run")
}

func (s *RunState) OnUpdate(dt float64) error {
	// 移動入力がなければアイドル状態へ
	if !s.player.IsMoving() {
		return s.player.stateMachine.ChangeState("idle", nil)
	}
	// ジャンプ入力があればジャンプ状態へ
	if s.player.IsJumpPressed() {
		return s.player.stateMachine.ChangeState("jump", nil)
	}
	return nil
}

// ジャンプ状態
type JumpState struct {
	PlayerState
}

func NewJumpState(player *Player) *JumpState {
	return &JumpState{
		PlayerState: PlayerState{
			BaseState: state.NewBaseState("jump"),
			player:    player,
		},
	}
}

func (s *JumpState) OnEnter(data interface{}) {
	s.player.SetAnimation("jump")
	s.player.ApplyJumpForce()
}

func (s *JumpState) OnUpdate(dt float64) error {
	// 着地したらアイドル状態へ
	if s.player.IsGrounded() {
		return s.player.stateMachine.ChangeState("idle", nil)
	}
	return nil
} 