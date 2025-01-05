package process

type PCB struct {
	pid int
}

func (p *PCB) GetPID() int {
	return p.pid
}
