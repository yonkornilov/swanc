package lib

import (
	"strings"
	"text/template"
	"time"

	"k8s.io/apimachinery/pkg/labels"
)

const (
	confPerm = 0644
	confDir  = "/srv/swanc"
	confFile = "ipsec.conf"

	// secretsPath = "/etc/ipsec.secrets"
	reloadCmd = "/usr/sbin/ipsec update"

	nodeKey = "net.beta.appscode.com/vpn"
)

var (
	nodeSelector = labels.SelectorFromSet(map[string]string{
		nodeKey: "",
	})

	funcMap = template.FuncMap{
		"replace": strings.Replace,
	}

	cfgTemplate = template.Must(template.New("cfg").Funcs(funcMap).Parse(
		`# IPSec configuration generated by https://github.com/appscode/swanc
# DO NOT EDIT!

config setup
        # strictcrlpolicy=yes
        # uniqueids = no

conn %default
        ikelifetime=60m
        keylife=20m
        rekeymargin=3m
        keyingtries=1
        mobike=no
        keyexchange=ikev2

{{ if .HostIP }}
{{ range $peer_ip := $.NodeIPs }}{{ if ne .HostIP $peer_ip }}
conn {{ replace .HostIP "." "_" -1 }}__{{ replace $peer_ip "." "_" -1 }}
        authby=secret
        left={{ .HostIP }}
        right={{ $peer_ip }}
        type=transport
        auto=start
        esp=aes128gcm16!
{{ end }}{{ end }}
{{ end }}
`))
)

type TemplateData struct {
	HostIP  string
	NodeIPs []string
}

type Options struct {
	NodeName             string
	PreferredAddressType string
	QPS                  float32
	Burst                int
	ResyncPeriod         time.Duration
	MaxNumRequeues       int
}
