package collector

import (
	"net/http"
	"time"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func StartServer() {
	var err error
	if HttpsEnabled {
		err = StartHttpsServer()
	} else {
		err = StartHttpServer()
	}
	if err != nil {
		logs.Logger.Errorf("Start server failed, error is: %s", err.Error())
		logs.FlushLogAndExit(1)
	}
}

func StartHttpServer() error {
	server := &http.Server{
		Addr:         CloudConf.Global.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	logs.Logger.Infof("Start server at %s", CloudConf.Global.Port)
	if err := server.ListenAndServe(); err != nil {
		logs.Logger.Errorf("Error occur when start server %s", err.Error())
		return err
	}
	return nil
}
