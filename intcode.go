package intcode

import (
	"fmt"
	"strconv"
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
		op, modes := getOpCodeAndModes(opcode)
		//fmt.Printf("Op: %d, modes: %+v\n", op, modes)
		//fmt.Println(Memory)
		//time.Sleep(1 * time.Second)
		//fmt.Println("i, op: ", i, op)
		switch op {
		case add:
			args := m.getArguments(3, modes)
			m.Memory[args[2]] = args[0] + args[1]
			m.Position += 4
		case multi:
			args := m.getArguments(3, modes)
			m.Memory[args[2]] = args[0] * args[1]
			m.Position += 4
		case input:
			//fmt.Println(m.Memory)
			args := m.getArguments(1, modes)
			//fmt.Printf("Modes: %+v; Position: %d; Op: %d Args: %+v; Value At pos: %d\n", modes, m.Position, op, args, m.Memory[args[0]])
			if len(m.Input) > 0 {
				var in int
				in, m.Input = m.Input[0], m.Input[1:]
				m.Memory[m.Memory[args[0]]] = in
			} else {
				m.Memory[m.Memory[args[0]]] = m.Memory[args[0]]
			}
			//fmt.Println(m.Memory)
			m.Position += 2
			//fmt.Printf("In for %q is: %d\n", m.Name, m.Input)
			//fmt.Println(args)
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
			args := m.getArguments(1, modes)
			//fmt.Printf("args: %+v\n", args)
			//fmt.Printf("Modes: %+v; Position: %d; Op: %d Args: %+v\n", modes, m.Position, op, args)
			out = append(out, args[0])
			//out = append(out, oout)
			//fmt.Printf("Out of %q is: %+v\n", m.Name, out)
			m.Position += 2
		case jmp:
			args := m.getArguments(2, modes)
			if args[0] != 0 {
				m.Position = args[1]
			} else {
				m.Position += 3
			}
			//fmt.Printf("5 i: %d args: %+v\n", i, args)
		case jmpf:
			args := m.getArguments(2, modes)
			if args[0] == 0 {
				m.Position = args[1]
			} else {
				m.Position += 3
			}
			//fmt.Printf("6 i: %d args: %+v\n", i, args)
		case less:
			args := m.getArguments(3, modes)
			if args[0] < args[1] {
				m.Memory[args[2]] = 1
			} else {
				m.Memory[args[2]] = 0
			}
			//fmt.Printf("7 i: %d args: %+v\n", i, args)
			m.Position += 4
		case eq:
			args := m.getArguments(3, modes)
			if args[0] == args[1] {
				m.Memory[args[2]] = 1
			} else {
				m.Memory[args[2]] = 0
			}
			//fmt.Printf("8 i: %d args: %+v\n", i, args)
			m.Position += 4
		case adj:
			args := m.getArguments(1, modes)
			m.RelativeBase += args[0]
			m.Position += 2
		case 99:
			break loop
		default:
			m.Position += 4
		}
	}

	return out, true
}

func (m *Machine) getArguments(num int, modes []int) (args []int) {
	for p := 0; p < num; p++ {
		var mode int
		if p >= len(modes) {
			mode = 0
		} else {
			mode = modes[p]
		}
		switch mode {
		case position:
			// Because parameters that an instruction writes to is always in position mode.
			if p > 1 && p+1 == num {
				args = append(args, m.Memory[m.Position+p+1])
			} else {
				args = append(args, m.Memory[m.Memory[m.Position+p+1]])
			}
		case immediate:
			args = append(args, m.Memory[m.Position+p+1])
		case relative:
			if p > 1 && p+1 == num {
				//fmt.Println("Relative Base: ", m.RelativeBase)
				args = append(args, m.Memory[m.Position+p+1])
			} else {
				pos := m.RelativeBase + m.Memory[m.Position+p+1]
				args = append(args, m.Memory[pos])
			}
		}
	}
	return
}

func getOpCodeAndModes(opcode int) (o int, modes []int) {
	sop := strconv.Itoa(opcode)
	l := len(sop)
	if len(sop) == 1 {
		o, _ = strconv.Atoi(sop)
		return o, nil
	}
	o, _ = strconv.Atoi(sop[l-2:])
	smodes := sop[:l-2]
	for i := len(smodes) - 1; i >= 0; i-- {
		m, _ := strconv.Atoi(string(smodes[i]))
		modes = append(modes, m)
	}
	return
}
