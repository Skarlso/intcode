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
		codes := fmt.Sprintf("%05d", opcode)
		op, _ := strconv.Atoi(codes[3:])
		modes := codes
		switch op {
		case add:
			arg1, arg2, dest := m.getParameter(1, modes), m.getParameter(2, modes), m.getParameter(3, modes)
			m.Memory[dest] = m.Memory[arg1] + m.Memory[arg2]
			m.Position += 4
		case multi:
			arg1, arg2, dest := m.getParameter(1, modes), m.getParameter(2, modes), m.getParameter(3, modes)
			m.Memory[dest] = m.Memory[arg1] * m.Memory[arg2]
			m.Position += 4
		case input:
			dest := m.getParameter(1, modes)
			if len(m.Input) > 0 {
				var in int
				in, m.Input = m.Input[0], m.Input[1:]
				m.Memory[dest] = in
			} else {
				m.Memory[m.Memory[dest]] = m.Memory[dest]
			}
			m.Position += 2
		case output:
			dest := m.getParameter(1, modes)
			out = append(out, m.Memory[dest])
			m.Position += 2
		case jmp:
			arg1, arg2 := m.getParameter(1, modes), m.getParameter(2, modes)
			if m.Memory[arg1] != 0 {
				m.Position = m.Memory[arg2]
				continue
			}
			m.Position += 3
		case jmpf:
			arg1, arg2 := m.getParameter(1, modes), m.getParameter(2, modes)
			if m.Memory[arg1] == 0 {
				m.Position = m.Memory[arg2]
				continue
			}
			m.Position += 3
		case less:
			arg1, arg2, dest := m.getParameter(1, modes), m.getParameter(2, modes), m.getParameter(3, modes)
			if m.Memory[arg1] < m.Memory[arg2] {
				m.Memory[dest] = 1
			} else {
				m.Memory[dest] = 0
			}
			m.Position += 4
		case eq:
			arg1, arg2, dest := m.getParameter(1, modes), m.getParameter(2, modes), m.getParameter(3, modes)
			if m.Memory[arg1] == m.Memory[arg2] {
				m.Memory[dest] = 1
			} else {
				m.Memory[dest] = 0
			}
			m.Position += 4
		case adj:
			arg1 := m.getParameter(1, modes)
			m.RelativeBase += m.Memory[arg1]
			m.Position += 2
		case 99:
			break loop
		}
	}

	return out, true
}

func (m *Machine) getParameter(index int, mode string) (value int) {
	switch mode[3-index] {
	case '0':
		value = m.Memory[m.Position+index]
	case '1':
		value = m.Position + index
	case '2':
		value = m.RelativeBase + m.Memory[m.Position+index]
	}
	return
}
