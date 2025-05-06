package paloalto

import (
	"archive/tar"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nonsugar-go/tools/excel"
)

/* running-config.xml
 * <config>
 *  <mgt-config>
 *   <users/>
 *   <password-complexity/>
 *  </mgt-config>
 *  <shared/>
 *  <devices>
 *   <entry name="localhost.localdomain">
 *    <network>
 *     <interface>
 *      <ethernet/>
 *      <loopback/>
 *      <vlan/>
 *      <tunnel/>
 *      <aggregate-ethernet/>
 *     </interface>
 *     <vlan/>
 *     <virtual-wire/>
 *     <profiles/>
 *     <ike/>
 *     <qos/>
 *     <virtual-router>
 *      <entry>
 *       <protocol/>
 *       <interface/>
 *       <ecmp/>
 *       <routing-table>
 *        <ip>
 *         <static-route/>
 *        </ip>
 *       </routing-table>
 *      </entry>
 *     </virtual-router>
 *      <entry>
 *       <protocol/>
 *       <interface/>
 *       <routing-table>
 *        <ip>
 *         <static-route/>
 *        </ip>
 *       </routing-table>
 *       <ecmp/>
 *      </entry>
 *     </virtual-router>
 *     <tunnel/>
 *    </network>
 *    <deviceconfig>
 *     <system>
 *      <ip-address/>
 *      <netmask/>
 *      <update-server/>
 *      <update-schedule/>
 *      <timezone/>
 *      <service/>
 *      <snmp-settings/>
 *      <hostname/>
 *      <domain/>
 *      <locale/>
 *      <geo-location>
 *       </latitude>
 *       <longitude/>
 *      </geo-location>
 *      <dns-setting>
 *       <servers>
 *        <primary/>
 *        <secondary/>
 *       </servers>
 *       <default-gateway/>
 *      <ntp-servers>
 *       <primary-ntp-server>
 *        <ntp-server-address/>
 *        <authentication-type>
 *         <none/>
 *        </authentication-type>
 *       </primary-ntp-server>
 *       <secondary-ntp-server>
 *        <ntp-server-address/>
 *        <authentication-type>
 *         <none/>
 *        </authentication-type>
 *       </secondary-ntp-server>
 *      </ntp-servers>
 *      <server-verification/>
 *      <type>
 *       <static/>
 *      </type>
 *      <permitted-ip>
 *       <entry/>
 *      </permitted-ip>
 *      <device-telemetry>
 *       <device-health-performance/>
 *       <product-usage/>
 *       <threat-prevention/>
 *       <region/>
 *      </device-telemetry>
 *     </system>
 *    </deviceconfig>
 *    <vsys>
 *     <entry name="vsys1">
 *      <application/>
 *      <application-group/>
 *      <zone/>
 *      <service/>
 *      <service-group/>
 *      <rulebae/>
 *      <import/>
 *      <address/>
 *      <address-group/>
 *      <tag/>
 *      <log-settings/>
 *      <profiles/>
 *      <profile-group/>
 *      <external-list/>
 *      <certificate/>
 *      <device-object/>
 *     </entry>
 *    </vsys>
 *   </entry>
 *  </devices>
 * </config>
 */

// Config is root element
type Config struct {
	XMLName       xml.Name        `xml:"config"`
	Version       string          `xml:"version,attr"`
	DetailVersion string          `xml:"detail-version,attr"`
	Users         []Users         `xml:"mgt-config>users>entry"`
	Ethernet      []Ethernet      `xml:"devices>entry>network>interface>ethernet>entry"`
	VirtualRouter []VirtualRouter `xml:"devices>entry>network>virtual-router>entry"`
	Vsys          []Vsys          `xml:"devices>entry>vsys>entry"`
}

// Users is mgt-config>users>entry
type Users struct {
	Name         string        `xml:"name,attr"`
	Superuser    string        `xml:"permissions>role-based>superuser"`
	Devicereader *Devicereader `xml:"permissions>role-based>devicereader"`
}

type Devicereader struct {
	XMLName xml.Name `xml:"devicereader"`
}

func (d *Devicereader) String() string {
	if d != nil {
		return "yes"
	}
	return ""
}

// Ethernet is devices>entry>network>interface>ethernet
type Ethernet struct {
	Name                       string        `xml:"name,attr"`
	AggregateGroup             string        `xml:"aggregate-group"`
	PortPriority               string        `xml:"lacp>port-priority"`
	LinkState                  string        `xml:"link-state"`
	IP                         []EthernetIP  `xml:"layer3>ip>entry"`
	InterfaceManagementProfile string        `xml:"layer3>interface-management-profile"`
	NetflowProfile             string        `xml:"layer3>netflow-profile"`
	LLDPEnable                 string        `xml:"layer3>lldp>enable"`
	HA                         *EthernetHA   `xml:"ha"`
	Comment                    string        `xml:"comment"`
	EthernetUnits              EthernetUnits `xml:"units>entry"`
}

// EthernetIP is ip>entry
type EthernetIP struct {
	Name string `xml:"name,attr"`
}

func (e EthernetIP) String() string {
	return e.Name
}

// HA is ha
type EthernetHA struct {
	XMLName xml.Name `xml:"ha"`
}

func (e *EthernetHA) String() string {
	if e != nil {
		return "HA"
	}
	return ""
}

// EthernetUnits is units
type EthernetUnits struct {
	Name                       string       `xml:"name,attr"`
	IP                         []EthernetIP `xml:"ip>entry"`
	InterfaceManagementProfile string       `xml:"interface-management-profile"`
	Tag                        string       `xml:"tag"`
	Comment                    string       `xml:"comment"`
}

// VirtualRouter is devices>entry>network>virtual-router>entry
type VirtualRouter struct {
	Name        string        `xml:"name,attr"`
	Interface   []string      `xml:"interface>member"`
	StaticRoute []StaticRoute `xml:"routing-table>ip>static-route>entry"`
}

// StaticRoute is routing-table>ip>static-route>entry
type StaticRoute struct {
	// TODO: datail modify
	Name        string `xml:"name,attr"`
	Nexthop     string `xml:"nexthop>ip-address"`
	Bfd         string `xml:"bfd>profile"`
	Interface   string `xml:"interface"`
	Metric      string `xml:"metric"`
	Destination string `xml:"destination"`
}

// Vsys is devices>entry>vsys>entry
type Vsys struct {
	Name             string             `xml:"name,attr"`
	Zone             []Zone             `xml:"zone>entry"`
	Tag              []Tag              `xml:"tag>entry"`
	Address          []Address          `xml:"address>entry"`
	AddressGroup     []AddressGroup     `xml:"address-group>entry"`
	ApplicationGroup []ApplicationGroup `xml:"application-group>entry"`
	Service          []Service          `xml:"service>entry"`
	ServiceGroup     []ServiceGroup     `xml:"service-group>entry"`
	Security         []Security         `xml:"rulebase>security>rules>entry"`
}

// Zone is tag>entry
type Zone struct {
	Name        string   `xml:"name,attr"`
	Layer3      []string `xml:"network>layer3>member"`
	Description string   `xml:"description"` // TODO: Check
}

// Tag is tag>entry
type Tag struct {
	Name     string `xml:"name,attr"`
	Color    string `xml:"color"`
	Comments string `xml:"comments"`
}

// Address is address>entry
type Address struct {
	Name        string   `xml:"name,attr"`
	IPNetmask   string   `xml:"ip-netmask"`
	FQDN        string   `xml:"fqdn"`
	Tag         []string `xml:"tag>member"`
	Description string   `xml:"description"`
}

// AddressGroup is address-group>entry
type AddressGroup struct {
	Name        string   `xml:"name,attr"`
	Static      []string `xml:"static>member"`
	Tag         []string `xml:"tag>member"`
	Description string   `xml:"description"`
}

// ApplicationGroup is application-group>entry
type ApplicationGroup struct {
	Name        string   `xml:"name,attr"`
	Member      []string `xml:"members>member"`
	Tag         []string `xml:"tag>member"`  // TODO: Check
	Description string   `xml:"description"` // TODO: Check
}

// Service is service>entry
type Service struct {
	Name        string `xml:"name,attr"`
	TCP         TCP    `xml:"protocol>tcp"`
	UDP         UDP    `xml:"protocol>udp"`
	Description string `xml:"description"`
}

// TCP is tcp>port
type TCP struct {
	Port string `xml:"port"`
}

// UDP is udp>port
type UDP struct {
	Port string `xml:"port"`
}

// ServiceGroup is service-group>entry
type ServiceGroup struct {
	Name   string   `xml:"name,attr"`
	Static []string `xml:"members>member"`
	Tag    []string `xml:"tag>member"`
}

// Security is rulebase>security>rules>entry
type Security struct {
	Name        string   `xml:"name,attr"`
	From        []string `xml:"from>member"`
	To          []string `xml:"to>member"`
	Source      []string `xml:"source>member"`
	Destination []string `xml:"destination>member"`
	Application []string `xml:"application>member"`
	Service     []string `xml:"service>member"`
	Action      string   `xml:"action"`
	Description string   `xml:"description"`
}

func parseConfig(inFile string) (*Config, *Vsys, error) {
	var data []byte
	var err error

	if strings.HasSuffix(inFile, ".tar.gz") || strings.HasSuffix(inFile, ".tgz") {
		file, err := os.Open(inFile)
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, nil, err
		}
		defer gzipReader.Close()
		tarReader := tar.NewReader(gzipReader)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, nil, err
			}

			if header.Name == "./running-config.xml" {
				data, err = io.ReadAll(tarReader)
				if err != nil {
					return nil, nil, err
				}
				break
			}
		}
	} else {
		data, err = os.ReadFile(inFile)
		if err != nil {
			return nil, nil, err
		}
	}
	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, nil, err
	}
	for _, entry := range config.Vsys {
		if entry.Name == "vsys1" {
			return &config, &entry, nil
		}
	}
	return nil, nil, fmt.Errorf("vsys1 not found")
}

// outputUsers is <zone> output process.
func outputUsers(xl *excel.Excel, config *Config) error {
	sheet := "ユーザー"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputUsers: %w", err)
	}
	xl.SetActiveSheet()
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"パスワード", 20}, {"Superuser", 12}, {"Devicereader", 12},
	}); err != nil {
		return fmt.Errorf("outputUsers: %w", err)
	}
	entries := config.Users
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		err := xl.SetRow(&[]any{
			i + 1, e.Name, "<REDACTED>", e.Superuser, e.Devicereader})
		if err != nil {
			return fmt.Errorf("outputUsers: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputUsers: %w", err)
	}
	return nil
}

// outputEthernet is <ethernet> output process.
func outputEthernet(xl *excel.Excel, config *Config) error {
	sheet := "イーサネット"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputEthernet: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 16}, {"集約グループ", 10}, {"ポート優先度", 10},
		{"リンク状態", 6}, {"IPアドレス", 18}, {"管理プロファイル", 10},
		{"Netflowプロファイル", 10}, {"LLDP", 8}, {"HA", 4}, {"内容", 60},
	}); err != nil {
		return fmt.Errorf("outputEthernet: %w", err)
	}
	entries := config.Ethernet
	sort.SliceStable(entries, func(i, j int) bool {
		is := strings.Split(entries[i].Name, "/")
		js := strings.Split(entries[j].Name, "/")
		if is[0] == js[0] {
			is1Int, _ := strconv.Atoi(is[1])
			js1Int, _ := strconv.Atoi(js[1])
			return is1Int < js1Int
		}
		return is[0] < is[1]
	})
	for i, e := range entries {
		err := xl.SetRow(&[]any{i + 1, e.Name, e.AggregateGroup, e.PortPriority, e.LinkState,
			e.IP, e.InterfaceManagementProfile, e.NetflowProfile,
			e.LLDPEnable, e.HA, e.Comment})
		if err != nil {
			return fmt.Errorf("outputEthernet: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputEthernet: %w", err)
	}
	return nil
}

// outputZone() is <zone> output process.
func outputZone(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "ゾーン"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputZone: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"タイプ", 10}, {"インターフェイス", 20}, {"内容", 60},
	}); err != nil {
		return fmt.Errorf("outputZone: %w", err)
	}
	entries := vsys1.Zone
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		typ := ""
		var member []string
		if e.Layer3 != nil {
			typ = "Layer3"
			member = e.Layer3
		}
		err := xl.SetRow(&[]any{i + 1, e.Name, typ, member, e.Description})
		if err != nil {
			return fmt.Errorf("outputZone: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputZone: %w", err)
	}
	return nil
}

// outputVirtualRouterInterface() is <interface> output process.
func outputVirtualRouterInterface(xl *excel.Excel, config *Config) error {
	sheet := "VRインターフェイス"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputVirtualRouterInterface: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"Virtual Router", 20}, {"インターフェイス", 20},
	}); err != nil {
		return fmt.Errorf("outputVirtualRouterInterface: %w", err)
	}
	entries := config.VirtualRouter
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	r := 0
	for _, e := range entries {
		// TODO: sort of Interface
		for _, member := range e.Interface {
			r++
			err := xl.SetRow(&[]any{r, e.Name, member})
			if err != nil {
				return fmt.Errorf("outputVirtualRouterInterface: %w", err)
			}
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputVirtualRouterInterface: %w", err)
	}
	return nil
}

// outputVirtualRouterStaticRoute() is <static-route> output process.
func outputVirtualRouterStaticRoute(xl *excel.Excel, config *Config) error {
	sheet := "VRスタティックルート"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputVirtualRouterStaticRoute: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"Virtual Router", 20}, {"名前", 20}, {"宛先", 20},
		{"インターフェイス", 20}, {"Nexthop", 14}, {"メトリック", 4}, {"Bfd", 10},
	}); err != nil {
		return fmt.Errorf("outputVirtualRouterStaticRoute: %w", err)
	}
	entries := config.VirtualRouter
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	r := 0
	for _, e := range entries {
		// TODO: sort of StaticRoute
		for _, e2 := range e.StaticRoute {
			r++
			err := xl.SetRow(&[]any{r, e.Name, e2.Name, e2.Destination, e2.Interface, e2.Nexthop,
				e2.Metric, e2.Bfd})
			if err != nil {
				return fmt.Errorf("outputVirtualRouterStaticRoute: %w", err)
			}
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputVirtualRouterStaticRoute: %w", err)
	}
	return nil
}

// outputTag() is <tag> output process.
func outputTag(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "タグ"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputTag: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"色", 20}, {"コメント", 60},
	}); err != nil {
		return fmt.Errorf("outputTag: %w", err)
	}
	entries := vsys1.Tag
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		color := e.Color
		if color_name, ok := colorMap[color]; ok {
			color = color_name
		}
		err := xl.SetRow(&[]any{
			i + 1, e.Name, color, e.Comments})
		if err != nil {
			return fmt.Errorf("outputUsers: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputTag: %w", err)
	}
	return nil
}

// outputAddress() is <address> output process.
func outputAddress(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "アドレス"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputAddress: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"アドレス", 20}, {"タグ", 12}, {"内容", 60},
	}); err != nil {
		return fmt.Errorf("outputAddress: %w", err)
	}
	entries := vsys1.Address
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		content := ""
		if e.IPNetmask != "" {
			content = e.IPNetmask
		} else if e.FQDN != "" {
			content = e.FQDN
		}
		err := xl.SetRow(&[]any{
			i + 1, e.Name, content, e.Tag, e.Description})
		if err != nil {
			return fmt.Errorf("outputAddress: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputAddress: %w", err)
	}
	return nil
}

// outputAddressGroup() is <address-group> output process.
func outputAddressGroup(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "アドレスグループ"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputAddressGroup: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"メンバー", 60}, {"タグ", 12}, {"内容", 60},
	}); err != nil {
		return fmt.Errorf("outputAddressGroup: %w", err)
	}
	entries := vsys1.AddressGroup
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		err := xl.SetRow(&[]any{
			i + 1, e.Name, e.Static, e.Tag, e.Description})
		if err != nil {
			return fmt.Errorf("outputAddressGroup: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputAddressGroup: %w", err)
	}
	return nil
}

// outputApplicationGroup() is <application-group> output process.
func outputApplicationGroup(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "アプリケーショングループ"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputApplicationGroup: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"メンバー", 60}, {"タグ", 12}, {"内容", 60},
	}); err != nil {
		return fmt.Errorf("outputApplicationGroup: %w", err)
	}
	entries := vsys1.ApplicationGroup
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		err := xl.SetRow(&[]any{
			i + 1, e.Name, e.Member, e.Tag, e.Description})
		if err != nil {
			return fmt.Errorf("outputApplicationGroup: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputApplicationGroup: %w", err)
	}
	return nil
}

// outputService() is <service> output process.
func outputService(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "サービス"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputService: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 30}, {"プロトコル", 10},
		{"宛先ポート", 20}, {"内容", 60},
	}); err != nil {
		return fmt.Errorf("outputService: %w", err)
	}
	entries := vsys1.Service
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		protocol := ""
		port := ""
		if e.TCP.Port != "" {
			protocol = "TCP"
			port = e.TCP.Port
		}
		if e.UDP.Port != "" {
			if e.TCP.Port != "" {
				log.Error("service object has tcp port and udp port")
			}
			protocol = "UDP"
			port = e.UDP.Port
		}
		err := xl.SetRow(&[]any{
			i + 1, e.Name, protocol, port, e.Description})
		if err != nil {
			return fmt.Errorf("outputService: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputService: %w", err)
	}
	return nil
}

// outputServiceGroup() is <service-group> output process.
func outputServiceGroup(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "サービスグループ"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputServiceGroup: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"メンバー", 60}, {"タグ", 12},
	}); err != nil {
		return fmt.Errorf("outputServiceGroup: %w", err)
	}
	entries := vsys1.ServiceGroup
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	for i, e := range entries {
		err := xl.SetRow(&[]any{
			i + 1, e.Name, e.Static, e.Tag})
		if err != nil {
			return fmt.Errorf("outputServiceGroup: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputServiceGroup: %w", err)
	}
	return nil
}

// outputSecurity() is <security> output process.
func outputSecurity(xl *excel.Excel, vsys1 *Vsys) error {
	sheet := "セキュリティ"
	if err := xl.NewSheet(sheet); err != nil {
		return fmt.Errorf("outputSecurity: %w", err)
	}
	if err := xl.SetHeader([]excel.Header{
		{"#", 4}, {"名前", 20}, {"送信元ゾーン", 10},
		{"宛先ゾーン", 10}, {"送信元", 30}, {"宛先", 30},
		{"アプリケーション", 30}, {"サービス", 30}, {"アクション", 10},
		{"内容", 60},
	}); err != nil {
		return fmt.Errorf("outputSecurity: %w", err)
	}
	entries := vsys1.Security
	for i, e := range entries {
		err := xl.SetRow(&[]any{
			i + 1, e.Name, e.From, e.To, e.Source, e.Destination, e.Application,
			e.Service, e.Action, e.Description})
		if err != nil {
			return fmt.Errorf("outputSecurity: %w", err)
		}
	}
	if err := xl.AddTable(sheet); err != nil {
		return fmt.Errorf("outputSecurity: %w", err)
	}
	return nil
}

// writeExcel outputs parameter sheets to Excel.
func writeExcel(outFile string, config *Config, vsys1 *Vsys) error {
	xl, err := excel.New(outFile)
	if err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}
	defer func() {
		if err := xl.Close(); err != nil {
			log.Errorf("WriteExcel: %v", err)
		}
	}()

	// <users>
	if err := outputUsers(xl, config); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <ethernet>
	if err := outputEthernet(xl, config); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <zone>
	if err := outputZone(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <interface>
	if err := outputVirtualRouterInterface(xl, config); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <static-route>
	if err := outputVirtualRouterStaticRoute(xl, config); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <tag>
	if err := outputTag(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <address>
	if err := outputAddress(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <address-group>
	if err := outputAddressGroup(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <application-group>
	if err := outputApplicationGroup(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <service>
	if err := outputService(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <service-group>
	if err := outputServiceGroup(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	// <security>
	if err := outputSecurity(xl, vsys1); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}

	if err := xl.SaveAndClose(); err != nil {
		return fmt.Errorf("WriteExcel: %w", err)
	}
	return nil
}

// ConvertPAConfig converts a PaloAlto config to a parameter sheet
func ConvertPAConfig(inFile, outFile string) error {
	log.Infof("input file: %s\n", inFile)
	log.Infof("output file: %s\n", outFile)
	config, vsys1, err := parseConfig(inFile)
	if err != nil {
		return fmt.Errorf("ConvertPAConfig: %w", err)
	}
	log.Infof("The sw-version of the analyzed config is %s (detail: %s)",
		config.Version, config.DetailVersion)
	_, err = os.Stat(outFile)
	if !os.IsNotExist(err) {
		log.Infof("delete the file: %s\n", outFile)
		err := os.Remove(outFile)
		if err != nil {
			return fmt.Errorf("ConvertPAConfig: %w", err)
		}
	}
	if err = writeExcel(outFile, config, vsys1); err != nil {
		return fmt.Errorf("ConvertPAConfig: %w", err)
	}
	log.Infof("out put the excel file: %s\n", outFile)
	return nil
}
