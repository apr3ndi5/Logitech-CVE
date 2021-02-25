package main

import (
	//"fmt"
	"syscall"
	"unsafe"

	"github.com/hillu/go-ntdll"
)

type MouseIo struct {
	button int8
	x      int8
	y      int8
	wheel  int8
	unk1   int8
}

var Imput ntdll.Handle

func main() {
	Imput = DeviceInitialize("\\??\\ROOT#SYSTEM#0002#{1abc05c0-c378-41b9-9cef-df1aba82b015}")

	if Imput == 0 {
		Imput = DeviceInitialize("\\??\\ROOT#SYSTEM#0001#{1abc05c0-c378-41b9-9cef-df1aba82b015}")
	}

	//fmt.Println(Imput)
	MouseMove(0, -100, 0, 0)
}

func CallMouse(buffer *MouseIo) {
	if Imput == ntdll.Handle(syscall.InvalidHandle) {
		//panic(Imput) //fmt.Println("[ ] Handle Invalida")
	}
	const sz = int(unsafe.Sizeof(MouseIo{}))
	var asByteSlice []byte = (*(*[sz]byte)(unsafe.Pointer(buffer)))[:]
	err := syscall.DeviceIoControl(syscall.Handle(Imput), 0x2a2010, &asByteSlice[0], uint32(sz), nil, uint32(0), nil, nil)
	if err != nil {
		//panic(err) //fmt.Printf("[Device Io Control] %v\n", err)
	}
}

func DeviceInitialize(device string) (ret ntdll.Handle) {
	attr := ntdll.NewObjectAttributes(device, 0, 0, nil)
	//fmt.Print(attr)
	var io ntdll.IoStatusBlock
	status := ntdll.NtCreateFile(&ret, syscall.GENERIC_WRITE|syscall.SYNCHRONIZE, attr, &io, nil, syscall.FILE_ATTRIBUTE_NORMAL, 0, 3, ntdll.FILE_NON_DIRECTORY_FILE|ntdll.FILE_SYNCHRONOUS_IO_NONALERT, nil, 0)
	if status != 0 {
		//panic(status) //fmt.Printf("[Create File] %v\n", status)
	}
	return
}

func MouseMove(button int8, x int8, y int8, wheel int8) {
	var io MouseIo
	io.unk1 = 0
	io.button = button
	io.x = x
	io.y = y
	io.wheel = wheel
	CallMouse(&io)

}

func MouseClose() {
	if Imput != 0 {
		syscall.CloseHandle(syscall.Handle(Imput))
		Imput = 0
	}
}
