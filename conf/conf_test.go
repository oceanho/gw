package conf

import (
	"gopkg.in/yaml.v2"
	"testing"
)

var f = `
#
# Gw framework configurations.
#
# -----------------------
#  Common configurations
# -----------------------
version: 1.0
server:
  listenAddr: ":8090"
  TLS:
    enabled: false
    certificate:
      key: "tls/app.key"
      cert: "tls/app.crt"
      format: file # base64/file

# -----------------------
#  Service configurations
# -----------------------
service:
  name: "My App"
  prefix: "/api/v1"
  version: "Version 1.0"
  remarks: "User Account Platform  Services."
  serviceDiscovery:
    disabled: True
    registryCenter:
      addr: ""
      configration:
        auth:
          user:
          password:
    configuration:
      heathCheck:
        proto: http
        url: https://xx.

# --------------------------------
#  Backend Store Configuration
# --------------------------------
backend:
  db:
  - name: primary
    driver: mysql
    addr: 127.0.0.1
    port: "3306"
    user: gw
    password: gw@123
    database: gwdb
    ssl_mode: on
    ssl_cert: ap
    args:
      charset: utf8
      parseTime: True
  cache:
  - name: primary
    driver: redis
    addr: 127.0.0.1
    port: "6379"
    type: redis
    auth:
      disable: True
      user: ocean
      database: "1"
      password: oceanho

# -------------------------------
#  Service Security Configuration
# -------------------------------
security:
  crypto:
    hash:
      alg: "sha256" # md5 / sha1 / sha256
      salt: "cvMC33eY7o9YKarcUr7VCf9XLFmHXKWJ"
    protect:
      alg: "aes"  # aes / 3des / rsa
      secret: "sPJF7yVnCvgctCJJ9scdTjwqzeenKnHH"
    cert:
      privateKey: "config/etc/gw.key"
      publicKey: "config/etc/gw.pem"
      isFile: True # true/false. The key,cert can be base64 string or A file path.
  auth:
    trustSidKey: "x-gw-trust-sid" # 如果配置这项，请务必确保您发送到后端节点的 http header -> <trustSidKey> 是可信的.
    paramKey:
      passport: "email"
      secret: "password"
      verifyCode: "verifyCode"
    paramPattern:
      passport: '^\S{5,64}$'
      secret: '^\S{5,64}$'
      verifyCode: '^[0-9]{0,8}$'
    session:
      defaultStore:
        name: primary
        prefix: gw-sid
    permission:
      defaultStore:
        type: db        # db/cache
        name: primary
    cookie:
      key: "_gid"
      domain: ""  # xxx.com / xxx.net
      path: "/"
      maxAge: "{{ .settings.expirationTimeControl.session }}"
      httpOnly: False
      secure: False
    allowUrls:
    - name: auth
      urls:
      - "POST:{{ .service.prefix }}/ucp/account/login"
      - "POST:{{ .service.prefix }}/ucp/account/login"
      ips:
      - 127.0.0.1/32
      - 192.168.0.0/24
    - name: gw-generator
      urls:
      - "GET:{{ .service.prefix }}/gw/generator/create-js"
      ips: []
  authServer:
    addr: https://auth.gw-framework.com
    enableAuthServe: True # Enable(True)/Disable(False) Auth features.
    login:
      url: "{{ .service.prefix }}/gw/auth/login"
      authTypes:
        - basicAuth
        - credentials
        - accessKeySecret
      methods:
        - "Get"
        - "Post"
    logout:
      url: "{{ .service.prefix }}/gw/auth/logout"
      methods:
        - "Get"
  limit:
    pagination:
        minPageSize: "20"
        maxPageSize: "2000"

# -------------------------------
#  Service Settings Configuration
# -------------------------------
settings:
  gw:
    printRouterInfo:
      disabled: False # enable/disable Print gw router info.
      title: "Router Information"
  headerKey:
    requestIdKey: "X-Request-Id"
  timeoutControl:  # Timeout control, units is millisecond
    redis: "1000"
    database: "2000"
    mongo: "2000"
    http: "2000"
  expirationTimeControl:
    session: "7200"

# Any Your custom configuration item at here.
# More: https://github.com/oceanho/gw/master/docs/configuration#custom
custom:
  gwpro:
    user: OceanHo
    tags:
    - body
    - programmer
    mystruct:
     a: hello
     b: 18
`

var bf = `
version: 1.0
configProvider: localfs
configuration:
  path: "config/app.yaml"
`

func TestLoadBootFromBytes_ShouldBe_OK(t *testing.T) {
	bsc := NewBootConfigFromBytes("yaml", []byte(bf))
	t.Logf("%v", bsc)
}

func TestLoadConfig_ShouldBe_OK(t *testing.T) {
	bsc := NewBootConfigFromBytes("yaml", []byte(bf))
	t.Logf("%v", bsc)
}

func TestApplicationConfig_ParseCustomPathTo(t *testing.T) {
	var appcnf ApplicationConfig
	yaml.Unmarshal([]byte(f), &appcnf)
	s := struct {
		User     string   `json:"user" toml:"user" yaml:"user"`
		Tags     []string `json:"tags" toml:"tags" yaml:"tags"`
		MyStruct struct {
			A string `json:"a" toml:"a" yaml:"a"`
			B int64  `json:"b" toml:"b" yaml:"b"`
		} `json:"mystruct" toml:"mystruct" yaml:"mystruct"`
	}{}
	_ = appcnf.ParseCustomPathTo("gwpro", &s)
}
