package process

import "testing"

func TestPCB_GetPID(t *testing.T) {
	t.Run("should return the pid", func(t *testing.T) {
		pcb := &PCB{pid: 1}

		pid := pcb.GetPID()

		if pid != 1 {
			t.Errorf("expected pid to be 1, got %d", pid)
		}
	})
}
