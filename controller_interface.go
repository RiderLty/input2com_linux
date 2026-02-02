package main

type mouseKeyboard interface {
	MouseMove(dx, dy, Wheel int32) error //应当在内部实现超过127的处理
	MouseBtnDown(keyCode byte) error
	MouseBtnUp(keyCode byte) error
	KeyDown(keyCode byte) error
	KeyUp(keyCode byte) error
}
