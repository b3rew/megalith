package main 

import "fmt"

func Process(server Server, servIndex int) {

	logcurrent := RequestLog{}
	ShouldDeleteLog(server.ID)
	LoadLog(server.ID, &logcurrent)

	for _, endpointCheck := range server.Endpoints {
		reqframe := Req(server, endpointCheck)
		apiid := fmt.Sprintf(urlformat, endpointCheck.Method, endpointCheck.Path)
		logcurrent.Requests = append(logcurrent.Requests, Request{Code: reqframe, Owner: apiid})
	}

	SaveLog(server.ID, &logcurrent)
	for endIndex, endpointCheck := range server.Endpoints {
		success := 0
		failed := 0
		apiid := fmt.Sprintf(urlformat, endpointCheck.Method, endpointCheck.Path)
		for _, reqcap := range logcurrent.Requests {
			if reqcap.Owner == apiid {
				if reqcap.Code < 300 {
					success++
				} else {
					failed++
				}
			}

		}
		GL.Lock.Lock()
		Config.Servers[servIndex].Endpoints[endIndex].Uptime = float64(success) / float64(success+failed)
		GL.Lock.Unlock()
	}
	success := 0
	failed := 0
	for _, reqcap := range logcurrent.Requests {
		if reqcap.Code < 300 {
			success++
		} else {
			failed++
		}
	}
	GL.Lock.Lock()
	Config.Servers[servIndex].Uptime = float64(success) / float64(success+failed)
	GL.Lock.Unlock()
	go Notify(Config.Servers[servIndex], Config.Contacts, Config.Mail)
	go SaveConfig(&Config)

}