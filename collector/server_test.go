package collector

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestStartServer(t *testing.T) {
	successFlag := true
	testCases := []struct {
		name    string
		patches func() *gomonkey.Patches
		expect  func(t *testing.T)
	}{
		{
			"start_https_success",
			func() *gomonkey.Patches {
				HttpsEnabled = true
				patches := getPatches()
				patches.ApplyFuncReturn(StartHttpsServer, nil)
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {
					successFlag = false
				})
				return patches
			},
			func(t *testing.T) {
				assert.True(t, successFlag)
			},
		},
		{
			"start_https_failed",
			func() *gomonkey.Patches {
				HttpsEnabled = true
				patches := getPatches()
				patches.ApplyFuncReturn(StartHttpsServer, errors.New("start https server error"))
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {
					successFlag = false
				})
				return patches
			},
			func(t *testing.T) {
				assert.False(t, successFlag)
			},
		},
		{
			"start_http_success",
			func() *gomonkey.Patches {
				HttpsEnabled = false
				patches := getPatches()
				patches.ApplyMethod(reflect.TypeOf(&http.Server{}), "ListenAndServe", func(server *http.Server) error {
					return nil
				})
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {
					successFlag = false
				})
				return patches
			},
			func(t *testing.T) {
				assert.True(t, successFlag)
			},
		},
		{
			"start_http_failed",
			func() *gomonkey.Patches {
				HttpsEnabled = false
				patches := getPatches()
				patches.ApplyMethod(reflect.TypeOf(&http.Server{}), "ListenAndServe", func(server *http.Server) error {
					return errors.New("start http server error")
				})
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {
					successFlag = false
				})
				return patches
			},
			func(t *testing.T) {
				assert.False(t, successFlag)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			patches := testCase.patches()
			defer patches.Reset()
			logs.InitLog("")
			successFlag = true
			StartServer()
			testCase.expect(t)
		})
	}
}
