package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/goexl/db"
	"github.com/goexl/gox/rand"
	"github.com/harluo/xorm/internal/config/internal"
)

type Server struct {
	// 数据库类型
	Type db.Type `default:"sqlite3" json:"type,omitempty" validate:"required,oneof=mysql sqlite sqlite3 mssql oracle postgres postgresql"` // nolint:lll

	// 主机
	Host string `json:"host,omitempty" validate:"required,hostname|ip"`
	// 端口
	Port int `default:"3306" json:"port,omitempty" validate:"required,max=65535"`
	// 用户名
	Username string `json:"username,omitempty"`
	// 密码
	Password string `json:"password,omitempty"`
	// 连接协议
	Protocol string `default:"tcp" json:"protocol,omitempty" validate:"required,oneof=tcp udp"`

	// 连接池配置
	Connection Connection `json:"connection,omitempty"`
	// 参数配置
	Sqlite internal.Sqlite `json:"sqlite,omitempty"`

	// 安全连接
	SSL *SSL `json:"ssl,omitempty"`
}

func (s *Server) SSLParam() {
	if s.SSL == nil {
		return
	}

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("/path/to/ca.pem") // 替换为你的CA文件路径
	if err != nil {
		log.Fatal(err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("failed to append CA cert")
	}

	// 2. 创建 TLS 配置
	clientCert := make([]tls.Certificate, 0)
	// 如果需要双向认证(mTLS)，需要在这里加载 client-cert.pem 和 client-key.pem
	// cert, err := tls.LoadX509KeyPair("/path/to/client-cert.pem", "/path/to/client-key.pem")
	// clientCert = append(clientCert, cert)

	tlsConfig := &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
		// ServerName: "mysql-server-hostname", // 仅在证书主机名与连接主机名不同时需要
	}

	// 3. 将 TLS 配置注册到驱动
	mysql.RegisterTLSConfig("custom", tlsConfig)
}

func (s *Server) loadCA(content []byte) (err error) {
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(content); !ok {
		log.Fatal("failed to append CA cert")
	} else if s.SSL.Cert != "" && s.SSL.Key != "" {

	} else {
		certificates := make([]tls.Certificate, 0)
		// 如果需要双向认证(mTLS)，需要在这里加载 client-cert.pem 和 client-key.pem
		// cert, err := tls.LoadX509KeyPair("/path/to/client-cert.pem", "/path/to/client-key.pem")
		// clientCert = append(clientCert, cert)

		tlsConfig := &tls.Config{
			RootCAs:      pool,
			Certificates: certificates,
			ServerName:   s.Host,
		}

		// 3. 将 TLS 配置注册到驱动
		mysql.RegisterTLSConfig("custom", tlsConfig)
	}

}

func (s *Server) loadKey() (err error) {
	if _, se := os.Stat(s.SSL.Cert); se != nil && os.IsNotExist(se) {

	}
	cert, err := tls.LoadX509KeyPair("/path/to/client-cert.pem", "/path/to/client-key.pem")
}

func (s *Server) loadFile(content string) (path string, err error) {
	if _, se := os.Stat(content); se != nil && os.IsNotExist(se) {
		os.CreateTemp
	}
}

func (s *Server) saveTmp(content string) (path string, err error) {
	filename := rand.New().String().Build().Generate()
	path = filepath.Join(os.TempDir(), filename)
	if tmp, cte := os.CreateTemp("", filename); cte != nil {
		err = cte
	} else {
		defer s.close(tmp)

		_, err = tmp.Write([]byte(content))
	}

	return
}

func (s *Server) close(file *os.File) {
	_= file.Close()
}
