package collector

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestStartHttpsServer(t *testing.T) {
	testCases := []struct {
		name    string
		patches func() *gomonkey.Patches
		expect  func(t *testing.T, err error)
	}{
		{
			"start_success",
			func() *gomonkey.Patches {
				patches := getPatches()
				patches.ApplyMethod(reflect.TypeOf(&http.Server{}), "Serve", func(srv *http.Server,
					l net.Listener) error {
					return nil
				})
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {})
				patches.ApplyFuncReturn(net.Listen, &net.TCPListener{}, nil)
				patches.ApplyMethod(reflect.TypeOf(&bufio.Reader{}), "ReadString", func(reader *bufio.Reader,
					delim byte) (string,
					error) {
					return "ca crt key password", nil
				})
				patches.ApplyFuncReturn(os.ReadFile, []byte{}, nil)
				patches.ApplyFuncReturn(os.Remove, nil)
				patches.ApplyFuncReturn(tls.X509KeyPair, tls.Certificate{}, nil)
				patches.ApplyFuncReturn(getServerKeyBytes, []byte("test"), nil)
				HttpsEnabled = true
				return patches
			},
			func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			"read_https_conf_error",
			func() *gomonkey.Patches {
				patches := getPatches()
				patches.ApplyMethod(reflect.TypeOf(&http.Server{}), "Serve", func(srv *http.Server,
					l net.Listener) error {
					return nil
				})
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {})
				patches.ApplyFuncReturn(net.Listen, &net.TCPListener{}, nil)
				patches.ApplyMethod(reflect.TypeOf(&bufio.Reader{}), "ReadString", func(reader *bufio.Reader,
					delim byte) (string,
					error) {
					return "", errors.New("read error")
				})
				patches.ApplyFuncReturn(os.ReadFile, []byte{}, nil)
				patches.ApplyFuncReturn(os.Remove, nil)
				patches.ApplyFuncReturn(tls.X509KeyPair, tls.Certificate{}, nil)
				patches.ApplyFuncReturn(getServerKeyBytes, []byte("test"), nil)
				HttpsEnabled = true
				return patches
			},
			func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			"read_https_conf_length_error",
			func() *gomonkey.Patches {
				patches := getPatches()
				patches.ApplyMethod(reflect.TypeOf(&http.Server{}), "Serve", func(srv *http.Server,
					l net.Listener) error {
					return nil
				})
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {})
				patches.ApplyFuncReturn(net.Listen, &net.TCPListener{}, nil)
				patches.ApplyMethod(reflect.TypeOf(&bufio.Reader{}), "ReadString", func(reader *bufio.Reader,
					delim byte) (string,
					error) {
					return "aaa", nil
				})
				patches.ApplyFuncReturn(os.ReadFile, []byte{}, nil)
				patches.ApplyFuncReturn(os.Remove, nil)
				patches.ApplyFuncReturn(tls.X509KeyPair, tls.Certificate{}, nil)
				patches.ApplyFuncReturn(getServerKeyBytes, []byte("test"), nil)
				HttpsEnabled = true
				return patches
			},
			func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			"start_without_security_mode",
			func() *gomonkey.Patches {
				patches := getPatches()
				patches.ApplyMethod(reflect.TypeOf(&http.Server{}), "Serve", func(srv *http.Server,
					l net.Listener) error {
					return nil
				})
				patches.ApplyFunc(logs.FlushLogAndExit, func(code int) {})
				patches.ApplyFuncReturn(net.Listen, &net.TCPListener{}, nil)
				patches.ApplyMethod(reflect.TypeOf(&bufio.Reader{}), "ReadString", func(reader *bufio.Reader,
					delim byte) (string,
					error) {
					return "aaa", nil
				})
				patches.ApplyFuncReturn(os.ReadFile, []byte{}, nil)
				patches.ApplyFuncReturn(os.Remove, nil)
				patches.ApplyFuncReturn(tls.X509KeyPair, tls.Certificate{}, nil)
				patches.ApplyFuncReturn(getServerKeyBytes, []byte("test"), nil)
				HttpsEnabled = false
				return patches
			},
			func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			patches := testCase.patches()
			defer patches.Reset()
			logs.InitLog("")
			err := StartHttpsServer()
			testCase.expect(t, err)
		})
	}
}

func TestTcpKeepAliveListener_Accept(t *testing.T) {
	testCases := []struct {
		name    string
		patches func() *gomonkey.Patches
		expect  func(t *testing.T, err error)
	}{
		{
			"accept_tpc_error",
			func() *gomonkey.Patches {
				patches := gomonkey.NewPatches()
				patches.ApplyMethod(reflect.TypeOf(&net.TCPListener{}), "AcceptTCP",
					func(listener *net.TCPListener) (*net.TCPConn, error) {
						return nil, errors.New("accept error")
					})
				return patches
			},
			func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"set_keep_alive_error",
			func() *gomonkey.Patches {
				patches := gomonkey.NewPatches()
				patches.ApplyMethod(reflect.TypeOf(&net.TCPListener{}), "AcceptTCP", func(listener *net.TCPListener) (*net.TCPConn, error) {
					return &net.TCPConn{}, nil
				})
				patches.ApplyMethod(reflect.TypeOf(&net.TCPConn{}), "SetKeepAlive", func(conn *net.TCPConn,
					keepalive bool) error {
					return errors.New("set keep alive error")
				})
				return patches
			},
			func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"set_keep_alive_period_error",
			func() *gomonkey.Patches {
				patches := gomonkey.NewPatches()
				patches.ApplyMethod(reflect.TypeOf(&net.TCPListener{}), "AcceptTCP", func(listener *net.TCPListener) (*net.
					TCPConn, error) {
					return &net.TCPConn{}, nil
				})
				patches.ApplyMethod(reflect.TypeOf(&net.TCPConn{}), "SetKeepAlive", func(conn *net.TCPConn,
					keepalive bool) error {
					return nil
				})
				patches.ApplyMethod(reflect.TypeOf(&net.TCPConn{}), "SetKeepAlivePeriod", func(conn *net.TCPConn, d time.Duration) error {
					return errors.New("set keep alive period error")
				})
				return patches
			},
			func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"success",
			func() *gomonkey.Patches {
				patches := gomonkey.NewPatches()
				patches.ApplyMethod(reflect.TypeOf(&net.TCPListener{}), "AcceptTCP",
					func(listener *net.TCPListener) (*net.TCPConn, error) {
						return &net.TCPConn{}, nil
					})
				patches.ApplyMethod(reflect.TypeOf(&net.TCPConn{}), "SetKeepAlive", func(conn *net.TCPConn, keepalive bool) error {
					return nil
				})
				patches.ApplyMethod(reflect.TypeOf(&net.TCPConn{}), "SetKeepAlivePeriod", func(conn *net.TCPConn, d time.Duration) error {
					return nil
				})
				return patches
			},
			func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			patches := testCase.patches()
			defer patches.Reset()
			ln := TcpKeepAliveListener{}
			_, err := ln.Accept()
			testCase.expect(t, err)
		})
	}
}

func TestCliExecutor(t *testing.T) {
	testCases := []struct {
		name    string
		patches func() *gomonkey.Patches
		expect  func(t *testing.T, err error)
	}{
		{
			"cmd_run_success",
			func() *gomonkey.Patches {
				patches := getPatches()
				cmd := &exec.Cmd{}
				patches.ApplyMethod(reflect.TypeOf(cmd), "Run", func(cmd *exec.Cmd) error {
					return nil
				})
				patches.ApplyFuncReturn(exec.Command, cmd)
				return patches
			},
			func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			"cmd_run_error",
			func() *gomonkey.Patches {
				patches := getPatches()
				cmd := &exec.Cmd{}
				patches.ApplyMethod(reflect.TypeOf(cmd), "Run", func(cmd *exec.Cmd) error {
					return errors.New("run cmd error")
				})
				patches.ApplyFuncReturn(exec.Command, cmd)
				return patches
			},
			func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			patches := testCase.patches()
			defer patches.Reset()
			logs.InitLog("")
			_, err := getServerKeyBytes("./server.key", "aaa")
			testCase.expect(t, err)
		})
	}
}

func TestVerifyPeerCrt(t *testing.T) {
	verifiedChains := [][]*x509.Certificate{
		{
			{
				DNSNames: []string{"aaa", "bbb"},
			},
			{
				DNSNames: []string{"ccc", "ddd"},
			},
		},
	}
	CloudConf.Global.ClientCN = "aaa,bbb,ccc,ddd"
	err := verifyPeerCert(nil, verifiedChains)
	assert.Nil(t, err)

	err = verifyPeerCert(nil, [][]*x509.Certificate{})
	assert.NotNil(t, err)
}
