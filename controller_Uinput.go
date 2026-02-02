package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/kenshaw/evdev"
	"github.com/lunixbochs/struc"
	"golang.org/x/sys/unix"
)

type EventMap struct {
	data   []byte
	Events []evdev.Event
}

const (
	ABS_MT_POSITION_X  = 0x35
	ABS_MT_POSITION_Y  = 0x36
	ABS_MT_SLOT        = 0x2F
	ABS_MT_TRACKING_ID = 0x39
	EV_SYN             = 0x00
	EV_KEY             = 0x01
	EV_REL             = 0x02
	EV_ABS             = 0x03
	REL_X              = 0x00
	REL_Y              = 0x01
	REL_WHEEL          = 0x08
	REL_HWHEEL         = 0x06
	SYN_REPORT         = 0x00
	BTN_TOUCH          = 0x14A
)

func makeEventsMMap(size int) EventMap {
	sizeofEvent := int(unsafe.Sizeof(evdev.Event{}))
	var byteSlice = make([]byte, size)
	dataPtr := unsafe.Pointer(&byteSlice[0])
	eventSlice := unsafe.Slice((*evdev.Event)(dataPtr), size/sizeofEvent)
	// return byteSlice, eventSlice
	return EventMap{data: byteSlice, Events: eventSlice}
}

func toUInputName(name []byte) [uinputMaxNameSize]byte {
	var fixedSizeName [uinputMaxNameSize]byte
	copy(fixedSizeName[:], name)
	return fixedSizeName
}

func uInputDevToBytes(uiDev UinputUserDev) []byte {
	var buf bytes.Buffer
	_ = struc.PackWithOptions(&buf, &uiDev, &struc.Options{Order: binary.LittleEndian})
	return buf.Bytes()
}

func createDevice(f *os.File) (err error) {
	return ioctl(f.Fd(), UIDEVCREATE(), uintptr(0))
}

func create_u_input_mouse_keyboard(devName string) *os.File {
	deviceFile, err := os.OpenFile("/dev/uinput", syscall.O_WRONLY|syscall.O_NONBLOCK, 0660)
	if err != nil {
		logger.Errorf("create u_input dev error:%v", err)
		return nil
	}
	ioctl(deviceFile.Fd(), UISETEVBIT(), evSyn)
	ioctl(deviceFile.Fd(), UISETEVBIT(), evKey)
	ioctl(deviceFile.Fd(), UISETEVBIT(), evRel)
	ioctl(deviceFile.Fd(), UISETEVRELBIT(), relX)
	ioctl(deviceFile.Fd(), UISETEVRELBIT(), relY)
	ioctl(deviceFile.Fd(), UISETEVRELBIT(), relWheel)
	ioctl(deviceFile.Fd(), UISETEVRELBIT(), relHWheel)
	for i := 0x110; i < 0x117; i++ {
		ioctl(deviceFile.Fd(), UISETKEYBIT(), uintptr(i))
	}
	for i := 0; i < 256; i++ {
		ioctl(deviceFile.Fd(), UISETKEYBIT(), uintptr(i))
	}

	uiDev := UinputUserDev{
		Name: toUInputName([]byte(devName)),
		ID: InputID{
			BusType: 0,
			Vendor:  uint16(rand.Intn(0x2000)),
			Product: uint16(rand.Intn(0x2000)),
			Version: uint16(rand.Intn(0x20)),
		},
		EffectsMax: 0,
		AbsMax:     [absCnt]int32{},
		AbsMin:     [absCnt]int32{},
		AbsFuzz:    [absCnt]int32{},
		AbsFlat:    [absCnt]int32{},
	}
	deviceFile.Write(uInputDevToBytes(uiDev))
	createDevice(deviceFile)
	return deviceFile
}

type UinputMouseKeyboard struct {
	mu        sync.Mutex
	device    *os.File // 保留文件对象
	deviceFd  int
	moveEvent EventMap
	keyEvent  EventMap
}

func NewMouseKeyboard_Uinput(devname string) *UinputMouseKeyboard {
	moveEvent := makeEventsMMap(4 * 24)
	keyEvent := makeEventsMMap(2 * 24)
	moveEvent.Events[0] = evdev.Event{Type: EV_REL, Code: REL_X, Value: 0}
	moveEvent.Events[1] = evdev.Event{Type: EV_REL, Code: REL_Y, Value: 0}
	moveEvent.Events[2] = evdev.Event{Type: EV_REL, Code: REL_WHEEL, Value: 0}
	moveEvent.Events[3] = evdev.Event{Type: EV_SYN, Code: SYN_REPORT, Value: 0}
	keyEvent.Events[0] = evdev.Event{Type: EV_KEY, Code: 0, Value: 0}
	keyEvent.Events[1] = evdev.Event{Type: EV_SYN, Code: SYN_REPORT, Value: 0}
	device := create_u_input_mouse_keyboard(devname)
	runtime.KeepAlive(device)
	if device == nil {
		return nil
	}
	return &UinputMouseKeyboard{
		device:    device,
		deviceFd:  int(device.Fd()),
		moveEvent: moveEvent,
		keyEvent:  keyEvent,
	}
}

func (mk *UinputMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.moveEvent.Events[0].Value = dx
	mk.moveEvent.Events[1].Value = dy
	mk.moveEvent.Events[2].Value = Wheel
	_, err := unix.Write(mk.deviceFd, mk.moveEvent.data)
	return err
}

func (mk *UinputMouseKeyboard) MouseBtnDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.keyEvent.Events[0].Code = uint16(hid2linux_Mouse[keyCode])
	mk.keyEvent.Events[0].Value = 1
	_, err := unix.Write(mk.deviceFd, mk.keyEvent.data)
	return err
}

func (mk *UinputMouseKeyboard) MouseBtnUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.keyEvent.Events[0].Code = uint16(hid2linux_Mouse[keyCode])
	mk.keyEvent.Events[0].Value = 0
	_, err := unix.Write(mk.deviceFd, mk.keyEvent.data)
	return err
}

func (mk *UinputMouseKeyboard) KeyDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.keyEvent.Events[0].Code = uint16(hid2linux[keyCode])
	mk.keyEvent.Events[0].Value = 1
	_, err := unix.Write(mk.deviceFd, mk.keyEvent.data)
	return err
}

func (mk *UinputMouseKeyboard) KeyUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.keyEvent.Events[0].Code = uint16(hid2linux[keyCode])
	mk.keyEvent.Events[0].Value = 0
	_, err := unix.Write(mk.deviceFd, mk.keyEvent.data)
	return err
}
