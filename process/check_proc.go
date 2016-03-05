// check_proc.go
package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/process"
)

func getProcessList() []process.Process {
	var proc []process.Process
	pids, _ := process.Pids()
	for i := 0; i < len(pids); i++ {
		p, _ := process.NewProcess(int32(pids[i]))
		proc = append(proc, *p)
	}
	return proc
}

//ps -eo pid
func getPsProcessList() []process.Process {
	cmd := exec.Command("ps", "-eo", "pid")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	pids := strings.Split(out.String(), "\n")[1:]
	var proc []process.Process
	for i := 0; i < len(pids); i++ {
		t, _ := strconv.Atoi(strings.TrimSpace(pids[i]))
		p, err := process.NewProcess(int32(t))
		if err == nil {
			proc = append(proc, *p)
		}
	}
	return proc
}

func main() {
	fmt.Println("[**Process**]: Estraggo i processi da ps...")
	pidps := getPsProcessList()
	fmt.Println("[**Process**]: Estraggo i processi da go...")
	pidgo := getProcessList()
	if len(pidgo) > len(pidps) {
		fmt.Println("[**Process**]: Attenzione ci sono processi nascosti!")
		fmt.Println("[**Process**]: Ricerco i processi nascosti in esecuzione...")
		for pid := 0; pid < len(pidgo); pid++ {
			res := 0
			for i := 0; i < len(pidps[2:]); i++ {
				if pidgo[pid].Pid == pidps[i].Pid {
					res++
				}
			}
			if res == 0 {
				name, _ := pidgo[pid].Name()
				user, _ := pidgo[pid].Username()
				cmd, _ := pidgo[pid].Cmdline()
				fmt.Println("[**Process**] Processo nascosto: ", name, " --Pid: ", pidgo[pid].Pid, " --User: ", user, " --Comando: ", cmd)
			}
		}
	} else {
		fmt.Println("[**Process**]: Non ci sono processi nascosti!")
	}
}
