package main

import (
	"fmt"
	"maps"
)

type Evaluator struct {
	gatesInOrder []Gate
}

func NewEvaluator(gates []Gate) (*Evaluator, error) {
	n := len(gates)

	type GateMeta struct {
		in1Ready, in2Ready bool
		g                  *Gate
	}

	gatesMeta := make([]GateMeta, n)

	consumers := map[string][]*GateMeta{}

	readyGates := make([]*GateMeta, n)
	nxtReadyGateIdx := 0

	isInput := func(name string) bool {
		return name[0] == 'x' || name[0] == 'y'
	}

	for i, g := range gates {
		gm := &gatesMeta[i]

		gm.g = &gates[i]

		consumers[g.in1] = append(consumers[g.in1], gm)
		consumers[g.in2] = append(consumers[g.in2], gm)

		if isInput(g.in1) {
			gm.in1Ready = true
		}
		if isInput(g.in2) {
			gm.in2Ready = true
		}
		if gm.in1Ready && gm.in2Ready {
			readyGates[nxtReadyGateIdx] = gm
			nxtReadyGateIdx++
		}
	}

	gatesInOrder := make([]Gate, n)

	for curIdx := range n {
		if curIdx >= nxtReadyGateIdx {
			return nil, fmt.Errorf("can't determine gates order")
		}

		gate := *readyGates[curIdx].g

		gatesInOrder[curIdx] = gate

		out := gate.out
		for _, gm := range consumers[out] {
			if gm.g.in1 == out {
				gm.in1Ready = true
			}
			if gm.g.in2 == out {
				gm.in2Ready = true
			}
			if gm.in1Ready && gm.in2Ready {
				readyGates[nxtReadyGateIdx] = gm
				nxtReadyGateIdx++
			}
		}
	}

	return &Evaluator{
		gatesInOrder: gatesInOrder,
	}, nil
}

func (e *Evaluator) Evaluate(initialState GatesState) (finalState GatesState) {
	finalState = maps.Clone(initialState)

	for _, gate := range e.gatesInOrder {
		left := finalState[gate.in1]
		right := finalState[gate.in2]
		res := false

		switch gate.op {
		case AND:
			res = left && right
		case XOR:
			res = left != right
		case OR:
			res = left || right
		default:
			panic(fmt.Errorf("unknown gate op: %v", gate.op))
		}

		finalState[gate.out] = res
	}

	return finalState
}
