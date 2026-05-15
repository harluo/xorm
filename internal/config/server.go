package config

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/goexl/db"
	"github.com/goexl/exception"
	"github.com/goexl/gox/field"
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

func (s *Server) sslParam(parameters internal.Parameters) (err error) {
	if s.SSL == nil {
		return
	}

	switch s.Type {
	case db.TypeMySQL:
		err = s.loadMysql(parameters)
	case db.TypePostgres:
		err = s.loadPostgres(parameters)
	default:
		err = exception.New().Message("暂未支持安全配置").Field(field.New("type", s.Type)).Build()
	}

	return
}

func (s *Server) loadPostgres(parameters internal.Parameters) (err error) {
	if caPath, lfe := s.loadFile(s.SSL.CA); lfe != nil {
		err = lfe
	} else {
		parameters["sslmode"] = "verify-full"
		parameters["sslrootcert"] = caPath
	}

	return
}

func (s *Server) loadMysql(parameters internal.Parameters) (err error) {
	name := rand.New().String().Build().Generate()
	if config, lke := s.loadKey(); lke != nil {
		err = lke
	} else if rce := mysql.RegisterTLSConfig(name, config); rce != nil {
		err = rce
	} else {
		parameters["tls"] = name
	}

	return
}

func (s *Server) loadKey() (config *tls.Config, err error) {
	pool := x509.NewCertPool()
	if lce := s.loadCA(pool); lce != nil {
		err = lce
	} else if cert, lke := s.loadKeypair(); lke != nil {
		err = lke
	} else {
		certificates := make([]tls.Certificate, 0)
		certificates = append(certificates, cert)
		config = &tls.Config{
			RootCAs:      pool,
			Certificates: certificates,
			ServerName:   s.Host,
		}
	}

	return
}

func (s *Server) loadCA(pool *x509.CertPool) (err error) {
	if caPath, lfe := s.loadFile(s.SSL.CA); lfe != nil {
		err = lfe
	} else if caData, rfe := os.ReadFile(caPath); rfe != nil {
		err = rfe
	} else if ok := pool.AppendCertsFromPEM(caData); !ok {
		err = exception.New().Message("CA证书未被支持").Field(field.New("type", s.Type)).Build()
	}

	return
}

func (s *Server) loadKeypair() (certificate tls.Certificate, err error) {
	if certPath, ce := s.loadFile(s.SSL.Cert); ce != nil {
		err = ce
	} else if keyPath, ke := s.loadFile(s.SSL.Key); ke != nil {
		err = ke
	} else if certPath != "" && keyPath != "" {
		certificate, err = tls.LoadX509KeyPair(certPath, keyPath)
	}

	return
}

func (s *Server) loadFile(content string) (path string, err error) {
	if content == "" {
		return
	}

	if _, se := os.Stat(content); se != nil && os.IsNotExist(se) {
		path, err = s.saveTmp(content)
	} else {
		path = content
	}

	return
}

func (s *Server) saveTmp(content string) (path string, err error) {
	filename := rand.New().String().Build().Generate()
	path = filepath.Join(os.TempDir(), filename)
	if tmp, cte := os.Create(path); cte != nil {
		err = cte
	} else {
		defer s.close(tmp)

		_, err = tmp.Write([]byte(content))
	}

	return
}

func (s *Server) close(file *os.File) {
	_ = file.Close()
}
