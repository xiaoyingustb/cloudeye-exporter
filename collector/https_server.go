package collector

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

var (
	caPath            string
	serverCertPath    string
	serverKeyPath     string
	serverKeyPassword string
)

func StartHttpsServer() error {
	caBytes, certBytes, keyBytes, loadErr := loadHttpsCerts()
	if loadErr != nil {
		logs.Logger.Errorf("Load https certificates failed, error is: %s", loadErr.Error())
		return loadErr
	}
	startErr := listenAndServeTLS(caBytes, certBytes, keyBytes)
	if startErr != nil {
		logs.Logger.Errorf("Listen and serve tls failed, error is: %s", startErr.Error())
		return startErr
	}
	return nil
}

func loadHttpsCerts() ([]byte, []byte, []byte, error) {
	// 命令行读取CA证书路径，https证书路径，https私钥路径和私钥密码
	readHttpsParametersFromCmd()
	// 保证读取完毕后，明文证书+私钥被删除
	defer cleanHttpsConfFiles()
	// 通过路径读取证书+私钥内容
	return loadHttpsCertsAndKey()
}

func readHttpsParametersFromCmd() {
	var err error
	fmt.Print("Please input CA certification path, https server certificate path, " +
		"https server private key path and private key password split with spaces.\n" +
		"eg: {example_ca_path example_server_cert_path example_server_key_path example_private_key_password}\n")
	reader := bufio.NewReader(os.Stdin)
	httpsParameters, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error occurred when scan https parameters: %s", err.Error())
		return
	}
	httpsParameterArray := strings.Split(httpsParameters, " ")
	if len(httpsParameterArray) < 4 {
		fmt.Printf("Read https parameters error, you should enter the https parameters in the following mode:\n" +
			"{example_ca_path example_server_cert_path example_server_key_path example_private_key_password}\n")
		return
	}
	caPath = httpsParameterArray[0]
	serverCertPath = httpsParameterArray[1]
	serverKeyPath = httpsParameterArray[2]
	serverKeyPassword = httpsParameterArray[3]
}

func loadHttpsCertsAndKey() ([]byte, []byte, []byte, error) {
	var err error
	caPath, err = NormalizePath(caPath)
	if err != nil {
		logs.Logger.Error("The ca certificate file path is abnormal and error is:", err.Error())
		return nil, nil, nil, err
	}
	caBytes, err := os.ReadFile(caPath)
	if err != nil {
		logs.Logger.Error("Read ca certificate file failed and error is:", err.Error())
		return nil, nil, nil, err
	}

	serverCertPath, err = NormalizePath(serverCertPath)
	if err != nil {
		logs.Logger.Error("The https server certificate file path is abnormal and error is:", err.Error())
		return nil, nil, nil, err
	}
	serverCertBytes, err := os.ReadFile(serverCertPath)
	if err != nil {
		logs.Logger.Error("Read https server certificate file failed and error is:", err.Error())
		return nil, nil, nil, err
	}

	serverKeyPath, err = NormalizePath(serverKeyPath)
	if err != nil {
		logs.Logger.Error("The https server private key file path is abnormal and error is:", err.Error())
		return nil, nil, nil, err
	}

	serverKeyBytes, err := getServerKeyBytes(serverKeyPath, serverKeyPassword)
	if err != nil {
		logs.Logger.Error("Open ssl get https server private key failed, and error is: ", err.Error())
		return nil, nil, nil, err
	}
	return caBytes, serverCertBytes, serverKeyBytes, nil
}

func cleanHttpsConfFiles() {
	err := os.Remove(caPath)
	if err != nil {
		logs.Logger.Errorf("Remove plaintext ca file error: %s", err.Error())
	}

	err = os.Remove(serverCertPath)
	if err != nil {
		logs.Logger.Errorf("Remove plaintext server crt file error: %s", err.Error())
	}
	err = os.Remove(serverKeyPath)
	if err != nil {
		logs.Logger.Errorf("Remove plaintext server private key file error: %s", err.Error())
	}
}

func getServerKeyBytes(serverKeyPath, serverKeyPassword string) ([]byte, error) {
	cmd := exec.Command("openssl", "rsa", "-in", serverKeyPath, "-passin", "stdin")
	cmd.Stdin = strings.NewReader(strings.Replace(serverKeyPassword, "'", "'\\''", -1))
	var (
		stdout, stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		logs.Logger.Error("Execute cmd failed and error is: ", err.Error())
		return nil, err
	}

	return stdout.Bytes(), nil
}

func listenAndServeTLS(caCert, serverCert, serverKey []byte) error {
	certificate, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	cfg := &tls.Config{
		ClientCAs:          pool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{certificate},
		MinVersion:         tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		VerifyPeerCertificate: verifyPeerCert,
	}
	server := &http.Server{
		Addr:         CloudConf.Global.Port,
		TLSConfig:    cfg,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServeTLS("", "")
}

func verifyPeerCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if len(verifiedChains) == 0 || len(verifiedChains[0]) == 0 {
		logs.Logger.Error("No verified chains")
		return fmt.Errorf("no verified chains")
	}
	cnArray := strings.Split(CloudConf.Global.ClientCN, ",")
	for i := range verifiedChains {
		for j := range verifiedChains[i] {
			for _, dnsName := range verifiedChains[i][j].DNSNames {
				if !strSliceContains(cnArray, dnsName) {
					logs.Logger.Errorf("Invalid dnsName current is %s", dnsName)
					return fmt.Errorf("invalid dnsName current is %s", dnsName)
				}
			}
		}
	}
	return nil
}

type TcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		logs.Logger.Errorf("Tcp listener accept tcp error: %s", err.Error())
		return
	}
	err = tc.SetKeepAlive(true)
	if err != nil {
		logs.Logger.Errorf("Tcp connector set keep alive error: %s", err.Error())
		return
	}
	err = tc.SetKeepAlivePeriod(3 * time.Minute)
	if err != nil {
		logs.Logger.Errorf("Tcp connector set keep alive period error: %s", err.Error())
		return
	}
	return tc, nil
}
