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
)

const (
	position = iota
	immediate
)

// Machine is a running environment for an intcode program.
type Machine struct {
	Position int
	Memory   map[int]int
	Input    []int
	Output   []int
	Name     string
}

// NewMachine returns an initialized intcode machine.
func NewMachine() Machine {
	return Machine{
		Input:  make([]int, 0),
		Output: make([]int, 0),
	}
}

// ProcessProgram will run an intcode.
func (m *Machine) ProcessProgram() (out []int, done bool) {
loop:
	for {
		opcode := m.Memory[m.Position]
		op, modes := getOpCodeAndModes(opcode)
		//fmt.Println(Memory)
		//time.Sleep(1 * time.Second)
		//fmt.Println("i, op: ", i, op)
		switch op {
		case add:
			args := getArguments(3, m.Position, modes, m.Memory)
			m.Memory[args[2]] = args[0] + args[1]
			m.Position += 4
		case multi:
			args := getArguments(3, m.Position, modes, m.Memory)
			m.Memory[args[2]] = args[0] * args[1]
			m.Position += 4
		case input:
			if len(m.Input) < 1 {
				//fmt.Printf("%q run out of input... returning\n", m.name)
				return out, false
			}
			var in int
			fmt.Printf("In for %q is: %d\n", m.Name, m.Input)
			in, m.Input = m.Input[0], m.Input[1:]
			m.Memory[m.Memory[m.Position+1]] = in
			m.Position += 2
		case output:
			var oout int
			if len(modes) > 0 {
				switch modes[0] {
				case position:
					oout = m.Memory[m.Memory[m.Position+1]]
				case immediate:
					oout = m.Memory[m.Position+1]
				}
			} else {
				oout = m.Memory[m.Memory[m.Position+1]]
			}
			out = append(out, oout)
			fmt.Printf("Out of %q is: %+v\n", m.Name, out)
			m.Position += 2
		case jmp:
			args := getArguments(2, m.Position, modes, m.Memory)
			if args[0] != 0 {
				m.Position = args[1]
			} else {
				m.Position += 3
			}
			//fmt.Printf("5 i: %d args: %+v\n", i, args)
		case jmpf:
			args := getArguments(2, m.Position, modes, m.Memory)
			if args[0] == 0 {
				m.Position = args[1]
			} else {
				m.Position += 3
			}
			//fmt.Printf("6 i: %d args: %+v\n", i, args)
		case less:
			args := getArguments(3, m.Position, modes, m.Memory)
			if args[0] < args[1] {
				m.Memory[args[2]] = 1
			} else {
				m.Memory[args[2]] = 0
			}
			//fmt.Printf("7 i: %d args: %+v\n", i, args)
			m.Position += 4
		case eq:
			args := getArguments(3, m.Position, modes, m.Memory)
			if args[0] == args[1] {
				m.Memory[args[2]] = 1
			} else {
				m.Memory[args[2]] = 0
			}
			//fmt.Printf("8 i: %d args: %+v\n", i, args)
			m.Position += 4
		case 99:
			break loop
		default:
			m.Position += 4
		}
	}

	return out, true
}

func getArguments(num, i int, modes []int, memory map[int]int) (args []int) {
	for p := 0; p < num; p++ {
		var m int
		if p >= len(modes) {
			m = 0
		} else {
			m = modes[p]
		}
		switch m {
		case position:
			// Because parameters that an instruction writes to is always in position mode.
			if p > 1 && p+1 == num {
				args = append(args, memory[i+p+1])
			} else {
				args = append(args, memory[memory[i+p+1]])
			}
		case immediate:
			args = append(args, memory[i+p+1])
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
