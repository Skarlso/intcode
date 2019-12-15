package intcode

import (
	"fmt"
	"strconv"
	"time"
)

const (
	add = iota + 1
	multi
	input
	output
	jmp
	jmpf
	less
	eq
	adj
)

const (
	position = iota
	immediate
	relative
)

// Machine is a running environment for an intcode program.
type Machine struct {
	Position     int
	Memory       map[int]int
	Input        []int
	Name         string
	RelativeBase int
}

// NewMachine returns an initialized intcode machine.
func NewMachine(m map[int]int) *Machine {
	return &Machine{
		Input:  make([]int, 0),
		Memory: m,
	}
}

func (m Machine) String() string {
	return fmt.Sprintf("Name: %s; position: %d; memory: %+v; input: %+v",
		m.Name,
		m.Position,
		m.Memory,
		m.Input)
}

// Reset the machine to zero state.
func (m *Machine) Reset() {
	m.Position = 0
	m.RelativeBase = 0
}

// ProcessProgram will run an intcode.
func (m *Machine) ProcessProgram() (out []int, done bool) {
loop:
	for {
		opcode := m.Memory[m.Position]
		codes := fmt.Sprintf("%05d", opcode)
		fmt.Println("Codes: ", codes)
		op, _ := strconv.Atoi(codes[3:])
		modes := codes
		fmt.Printf("Op: %d, modes: %+v\n", op, modes)
		fmt.Println(m.Memory)
		time.Sleep(1 * time.Second)
		switch op {
		case add:
			arg1, arg2 := m.getParameter(1, modes), m.getParameter(2, modes)
			dest := m.getDestination(3, modes)
			m.Memory[dest] = arg1 + arg2
			m.Position += 4
		case multi:
			//args := m.getArguments(3, modes)
			//m.Memory[args[2]] = args[0] * args[1]
			arg1, arg2 := m.getParameter(1, modes), m.getParameter(2, modes)
			dest := m.getDestination(3, modes)
			m.Memory[dest] = arg1 * arg2
			m.Position += 4
		case input:
			dest := m.getDestination(3, modes)
			//fmt.Printf("Modes: %+v; Position: %d; Op: %d Args: %+v; Value At pos: %d\n", modes, m.Position, op, args, m.Memory[args[0]])
			if len(m.Input) > 0 {
				var in int
				in, m.Input = m.Input[0], m.Input[1:]
				m.Memory[dest] = in
				//fmt.Println(m.Memory)
			} else {
				m.Memory[m.Memory[dest]] = dest
			}
			m.Position += 2
		case output:
			//var oout int
			//if len(modes) > 0 {
			//	switch modes[0] {
			//	case position:
			//		oout = m.Memory[m.Memory[m.Position+1]]
			//	case immediate:
			//		oout = m.Memory[m.Position+1]
			//	case relative:
			//		pos := m.RelativeBase + m.Memory[m.Position+1]
			//		oout = m.Memory[pos]
			//	}
			//} else {
			//	oout = m.Memory[m.Memory[m.Position+1]]
			//}
			//fmt.Println(m.Memory)
			dest := m.getDestination(1, modes)
			//fmt.Printf("args: %+v\n", args)
			fmt.Printf("Modes: %+v; Position: %d; Op: %d\n", modes, m.Position, op)
			//out = append(out, args[0])
			out = append(out, dest)
			//fmt.Printf("Out of %q is: %+v\n", m.Name, out)
			m.Position += 2
		case jmp:
			arg1 := m.getParameter(1, modes)
			if arg1 != 0 {
				m.Position = arg1
			} else {
				m.Position += 3
			}
			//fmt.Printf("5 i: %d args: %+v\n", i, args)
		case jmpf:
			arg1 := m.getParameter(1, modes)
			if arg1 == 0 {
				m.Position = arg1
			} else {
				m.Position += 3
			}
			//fmt.Printf("6 i: %d args: %+v\n", i, args)
		case less:
			arg1, arg2 := m.getParameter(1, modes), m.getParameter(2, modes)
			dest := m.getDestination(3, modes)
			if arg1 < arg2 {
				m.Memory[dest] = 1
			} else {
				m.Memory[dest] = 0
			}
			//fmt.Printf("7 i: %d args: %+v\n", i, args)
			m.Position += 4
		case eq:
			arg1, arg2 := m.getParameter(1, modes), m.getParameter(2, modes)
			dest := m.getDestination(3, modes)
			if arg1 == arg2 {
				m.Memory[dest] = 1
			} else {
				m.Memory[dest] = 0
			}
			//fmt.Printf("8 i: %d args: %+v\n", i, args)
			m.Position += 4
		case adj:
			arg1 := m.getParameter(1, modes)
			m.RelativeBase += arg1
			m.Position += 2
		case 99:
			break loop
		default:
			m.Position += 4
		}
	}

	return out, true
}

func (m *Machine) getParameter(index int, mode string) (address int) {
	switch mode[3-index] {
	case '0':
		address = m.Memory[m.Position+index]
	case '1':
		address = m.Position + index
	case '2':
		address = m.RelativeBase + m.Memory[m.Position+index]
	}
	return
}

func (m *Machine) getDestination(index int, mode string) (dest int) {
	fmt.Println("Mode 3-index: ", string(mode[3-index]))
	switch mode[3-index] {
	case '0', '1':
		dest = m.Memory[m.Position+index]
		fmt.Println("dest in 0,1: ", dest)
	case '2':
		dest = m.RelativeBase + m.Memory[m.Position+index]
		fmt.Println("dest in 2: ", dest)
	}
	return
}
