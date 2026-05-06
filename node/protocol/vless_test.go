package protocol

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestVlessEncodeDecode 测试 VLESS 编解码完整性
func TestVlessEncodeDecode(t *testing.T) {
	original := VLESS{
		Name:   "测试节点-VLESS",
		Uuid:   "12345678-1234-1234-1234-123456789abc",
		Server: "example.com",
		Port:   443,
		Query: VLESSQuery{
			Security:   "tls",
			Encryption: "none",
			Type:       "ws",
			Host:       "cdn.example.com",
			Path:       "/vless",
			Sni:        "sni.example.com",
			Fp:         "chrome",
			Alpn:       []string{"h2", "http/1.1"},
		},
	}

	// 编码
	encoded := EncodeVLESSURL(original)
	if !strings.HasPrefix(encoded, "vless://") {
		t.Errorf("编码后应以 vless:// 开头, 实际: %s", encoded)
	}

	// 解码
	decoded, err := DecodeVLESSURL(encoded)
	if err != nil {
		t.Fatalf("解码失败: %v", err)
	}

	// 验证关键字段
	assertEqualString(t, "Server", original.Server, decoded.Server)
	assertEqualIntInterface(t, "Port", original.Port, decoded.Port)
	assertEqualString(t, "Uuid", original.Uuid, decoded.Uuid)
	assertEqualString(t, "Name", original.Name, decoded.Name)
	assertEqualString(t, "Query.Type", original.Query.Type, decoded.Query.Type)
	assertEqualString(t, "Query.Sni", original.Query.Sni, decoded.Query.Sni)
	assertEqualString(t, "Query.Path", original.Query.Path, decoded.Query.Path)

	t.Logf("✓ VLESS 编解码测试通过，名称: %s", decoded.Name)
}

// TestVlessNameModification 测试 VLESS 名称修改
func TestVlessNameModification(t *testing.T) {
	original := VLESS{
		Name:   "原始名称",
		Uuid:   "12345678-1234-1234-1234-123456789abc",
		Server: "example.com",
		Port:   443,
		Query: VLESSQuery{
			Security: "tls",
			Type:     "tcp",
		},
	}

	newName := "新名称-VLESS-测试"
	encoded := EncodeVLESSURL(original)
	decoded, _ := DecodeVLESSURL(encoded)
	decoded.Name = newName
	reEncoded := EncodeVLESSURL(decoded)
	final, _ := DecodeVLESSURL(reEncoded)

	assertEqualString(t, "修改后名称", newName, final.Name)
	assertEqualString(t, "服务器(不变)", original.Server, final.Server)
	assertEqualString(t, "UUID(不变)", original.Uuid, final.Uuid)
	assertEqualIntInterface(t, "端口(不变)", original.Port, final.Port)

	t.Logf("✓ VLESS 名称修改测试通过: %s -> %s", original.Name, final.Name)
}

// TestVlessSpecialCharacters 测试 VLESS 特殊字符
func TestVlessSpecialCharacters(t *testing.T) {
	specialNames := []string{
		"节点 with spaces",
		"节点-with-dashes",
		"节点_with_underscores",
		"节点中文测试",
		"Node🚀Emoji",
		"Node (parentheses)",
	}

	for _, name := range specialNames {
		t.Run(name, func(t *testing.T) {
			original := VLESS{
				Name:   name,
				Uuid:   "12345678-1234-1234-1234-123456789abc",
				Server: "example.com",
				Port:   443,
				Query: VLESSQuery{
					Security: "tls",
					Type:     "tcp",
				},
			}

			encoded := EncodeVLESSURL(original)
			decoded, err := DecodeVLESSURL(encoded)
			if err != nil {
				t.Fatalf("解码失败: %v", err)
			}

			assertEqualString(t, "特殊字符名称", name, decoded.Name)
			t.Logf("✓ 特殊字符测试通过: %s", name)
		})
	}
}

// TestVlessV2rayFormat 测试 v2ray 格式 VLESS 链接解析（明文URL，非base64）
func TestVlessV2rayFormat(t *testing.T) {
	// 典型的v2ray格式VLESS链接
	testCases := []struct {
		name     string
		url      string
		expected VLESSQuery
	}{
		{
			name: "WebSocket传输层",
			url:  "vless://12345678-1234-1234-1234-123456789abc@example.com:443?encryption=none&security=tls&type=ws&host=cdn.example.com&path=%2Fvless&sni=example.com&fp=chrome#测试节点",
			expected: VLESSQuery{
				Security:   "tls",
				Encryption: "none",
				Type:       "ws",
				Host:       "cdn.example.com",
				Path:       "/vless",
				Sni:        "example.com",
				Fp:         "chrome",
			},
		},
		{
			name: "Reality配置",
			url:  "vless://12345678-1234-1234-1234-123456789abc@example.com:443?encryption=none&security=reality&type=tcp&flow=xtls-rprx-vision&pbk=testpublickey&sid=testshortid&sni=example.com&fp=chrome#Reality节点",
			expected: VLESSQuery{
				Security: "reality",
				Type:     "tcp",
				Flow:     "xtls-rprx-vision",
				Pbk:      "testpublickey",
				Sid:      "testshortid",
				Sni:      "example.com",
				Fp:       "chrome",
			},
		},
		{
			name: "gRPC传输层",
			url:  "vless://12345678-1234-1234-1234-123456789abc@example.com:443?encryption=none&security=tls&type=grpc&serviceName=mygrpc&mode=gun#gRPC节点",
			expected: VLESSQuery{
				Security:    "tls",
				Type:        "grpc",
				ServiceName: "mygrpc",
				Mode:        "gun",
			},
		},
		{
			name: "H2传输层",
			url:  "vless://12345678-1234-1234-1234-123456789abc@example.com:443?encryption=none&security=tls&type=h2&host=example.com&path=%2Fh2path#H2节点",
			expected: VLESSQuery{
				Security: "tls",
				Type:     "h2",
				Host:     "example.com",
				Path:     "/h2path",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := DecodeVLESSURL(tc.url)
			if err != nil {
				t.Fatalf("解码失败: %v", err)
			}

			assertEqualString(t, "Security", tc.expected.Security, decoded.Query.Security)
			assertEqualString(t, "Type", tc.expected.Type, decoded.Query.Type)
			if tc.expected.Host != "" {
				assertEqualString(t, "Host", tc.expected.Host, decoded.Query.Host)
			}
			if tc.expected.Path != "" {
				assertEqualString(t, "Path", tc.expected.Path, decoded.Query.Path)
			}
			if tc.expected.Flow != "" {
				assertEqualString(t, "Flow", tc.expected.Flow, decoded.Query.Flow)
			}
			if tc.expected.Pbk != "" {
				assertEqualString(t, "Pbk", tc.expected.Pbk, decoded.Query.Pbk)
			}
			if tc.expected.ServiceName != "" {
				assertEqualString(t, "ServiceName", tc.expected.ServiceName, decoded.Query.ServiceName)
			}

			t.Logf("✓ %s 测试通过", tc.name)
		})
	}
}

// TestVlessPacketEncoding 测试 packet-encoding 参数
func TestVlessPacketEncoding(t *testing.T) {
	url := "vless://12345678-1234-1234-1234-123456789abc@example.com:443?encryption=none&security=tls&type=tcp&packetEncoding=xudp#xudp节点"
	decoded, err := DecodeVLESSURL(url)
	if err != nil {
		t.Fatalf("解码失败: %v", err)
	}

	assertEqualString(t, "PacketEncoding", "xudp", decoded.Query.PacketEncoding)
	t.Logf("✓ packet-encoding 测试通过")
}

func TestVlessTopLevelECHRoundTrip(t *testing.T) {
	original := VLESS{
		Name:   "ECH节点",
		Uuid:   "12345678-1234-1234-1234-123456789abc",
		Server: "example.com",
		Port:   443,
		Query: VLESSQuery{
			Security:   "tls",
			Encryption: "none",
			Type:       "ws",
			Host:       "cdn.example.com",
			Path:       "/vless",
			Sni:        "example.com",
			Ech:        "encryptedsni.com+https://dns.alidns.com/dns-query",
		},
	}

	encoded := EncodeVLESSURL(original)
	assertContains(t, "EncodedECH", encoded, "ech=encryptedsni.com%2Bhttps%3A%2F%2Fdns.alidns.com%2Fdns-query")

	decoded, err := DecodeVLESSURL(encoded)
	if err != nil {
		t.Fatalf("解码失败: %v", err)
	}

	assertEqualString(t, "ECH", original.Query.Ech, decoded.Query.Ech)
}

func TestVlessXHTTPURLMapping(t *testing.T) {
	extra := map[string]interface{}{
		"headers": map[string]interface{}{
			"User-Agent": "curl/8.0",
		},
		"noGRPCHeader":  true,
		"xPaddingBytes": "10-20",
		"downloadSettings": map[string]interface{}{
			"path":              "/download",
			"host":              "dl.example.com",
			"tls":               true,
			"server":            "dl-backend.example.com",
			"port":              float64(8443),
			"clientFingerprint": "chrome",
		},
	}
	extraBytes, err := json.Marshal(extra)
	if err != nil {
		t.Fatalf("extra 编码失败: %v", err)
	}

	original := VLESS{
		Name:   "XHTTP节点",
		Uuid:   "12345678-1234-1234-1234-123456789abc",
		Server: "example.com",
		Port:   443,
		Query: VLESSQuery{
			Security:   "tls",
			Encryption: "none",
			Type:       "xhttp",
			Host:       "cdn.example.com",
			Path:       "/xhttp",
			Mode:       "stream-up",
			Sni:        "example.com",
			Extra:      string(extraBytes),
		},
	}

	encoded := EncodeVLESSURL(original)
	assertContains(t, "EncodedType", encoded, "type=xhttp")
	assertContains(t, "EncodedMode", encoded, "mode=stream-up")
	assertContains(t, "EncodedExtra", encoded, "extra=")

	decoded, err := DecodeVLESSURL(encoded)
	if err != nil {
		t.Fatalf("解码失败: %v", err)
	}

	assertEqualString(t, "Type", "xhttp", decoded.Query.Type)
	assertEqualString(t, "Host", original.Query.Host, decoded.Query.Host)
	assertEqualString(t, "Path", original.Query.Path, decoded.Query.Path)
	assertEqualString(t, "Mode", original.Query.Mode, decoded.Query.Mode)

	decodedExtra := parseVLESSXHTTPExtra(decoded.Query.Extra)
	if decodedExtra == nil {
		t.Fatal("decoded extra 不应为空")
	}
	assertEqualString(t, "DecodedHeader", "curl/8.0", decodedExtra["headers"].(map[string]interface{})["User-Agent"].(string))
	assertEqualString(t, "DecodedPadding", "10-20", decodedExtra["x-padding-bytes"].(string))
	assertEqualString(t, "DecodedDownloadPath", "/download", decodedExtra["download-settings"].(map[string]interface{})["path"].(string))
	assertEqualString(t, "DecodedDownloadFingerprint", "chrome", decodedExtra["download-settings"].(map[string]interface{})["client-fingerprint"].(string))
}

func TestConvertProxyToVlessXHTTP(t *testing.T) {
	proxy := Proxy{
		Name:    "XHTTP节点",
		Type:    "vless",
		Server:  "example.com",
		Port:    443,
		Uuid:    "12345678-1234-1234-1234-123456789abc",
		Network: "xhttp",
		Tls:     true,
		XHTTP_opts: map[string]interface{}{
			"path": "/xhttp",
			"host": "cdn.example.com",
			"mode": "packet-up",
			"headers": map[string]interface{}{
				"User-Agent": "curl/8.0",
			},
			"no-grpc-header": true,
			"download-settings": map[string]interface{}{
				"path":               "/download",
				"client-fingerprint": "chrome",
			},
		},
	}

	vless := ConvertProxyToVless(proxy)
	assertEqualString(t, "Type", "xhttp", vless.Query.Type)
	assertEqualString(t, "Host", "cdn.example.com", vless.Query.Host)
	assertEqualString(t, "Path", "/xhttp", vless.Query.Path)
	assertEqualString(t, "Mode", "packet-up", vless.Query.Mode)

	extra := parseVLESSXHTTPExtra(vless.Query.Extra)
	if extra == nil {
		t.Fatal("extra 不应为空")
	}
	var rawExtra map[string]interface{}
	if err := json.Unmarshal([]byte(vless.Query.Extra), &rawExtra); err != nil {
		t.Fatalf("extra 解析失败: %v", err)
	}
	assertEqualString(t, "ExtraHeader", "curl/8.0", rawExtra["headers"].(map[string]interface{})["User-Agent"].(string))
	assertEqualString(t, "ExtraDownloadPath", "/download", rawExtra["downloadSettings"].(map[string]interface{})["path"].(string))
	assertEqualString(t, "ExtraDownloadFingerprint", "chrome", rawExtra["downloadSettings"].(map[string]interface{})["clientFingerprint"].(string))

	encoded := EncodeVLESSURL(vless)
	assertContains(t, "EncodedType", encoded, "type=xhttp")
}

func TestConvertProxyToVlessPreservesTopLevelECH(t *testing.T) {
	proxy := Proxy{
		Name:       "ECH节点",
		Type:       "vless",
		Server:     "example.com",
		Port:       443,
		Uuid:       "12345678-1234-1234-1234-123456789abc",
		Network:    "ws",
		Tls:        true,
		Servername: "example.com",
		ECH_opts: map[string]interface{}{
			"enable": true,
			"config": "BASE64_ECH_CONFIG",
		},
	}

	vless := ConvertProxyToVless(proxy)
	assertEqualString(t, "ECH", "BASE64_ECH_CONFIG", vless.Query.Ech)

	encoded := EncodeVLESSURL(vless)
	assertContains(t, "EncodedECH", encoded, "ech=BASE64_ECH_CONFIG")
}

func TestVlessTopLevelECHMapsToECHOpts(t *testing.T) {
	original := VLESS{
		Name:   "ECH双路径节点",
		Uuid:   "12345678-1234-1234-1234-123456789abc",
		Server: "example.com",
		Port:   443,
		Query: VLESSQuery{
			Security:   "tls",
			Encryption: "none",
			Type:       "xhttp",
			Host:       "cdn.example.com",
			Path:       "/xhttp",
			Sni:        "example.com",
			Ech:        "BASE64_ECH_CONFIG",
			Extra:      `{"downloadSettings":{"echOpts":{"config":"base64-ech","queryServerName":"dns.example.com"}}}`,
		},
	}

	proxy, err := buildVLESSProxy(Urls{Url: EncodeVLESSURL(original)}, OutputConfig{})
	if err != nil {
		t.Fatalf("buildVLESSProxy 失败: %v", err)
	}

	assertEqualString(t, "TopLevelECHConfig", "BASE64_ECH_CONFIG", proxy.ECH_opts["config"].(string))
	assertEqualBool(t, "TopLevelECHEnable", true, proxy.ECH_opts["enable"].(bool))
	downloadSettings, ok := proxy.XHTTP_opts["download-settings"].(map[string]interface{})
	if !ok {
		t.Fatal("download-settings 不应为空")
	}
	xhttpECHOpts, ok := downloadSettings["ech-opts"].(map[string]interface{})
	if !ok {
		t.Fatal("ech-opts 不应为空")
	}
	assertEqualString(t, "NestedECHConfig", "base64-ech", xhttpECHOpts["config"].(string))
	assertEqualString(t, "NestedECHQueryServerName", "dns.example.com", xhttpECHOpts["query-server-name"].(string))

	restored := ConvertProxyToVless(proxy)
	assertEqualString(t, "RestoredTopLevelECH", original.Query.Ech, restored.Query.Ech)
	assertContains(t, "RestoredExtra", restored.Query.Extra, "\"echOpts\"")
	assertContains(t, "RestoredExtraQueryServerName", restored.Query.Extra, "\"queryServerName\":\"dns.example.com\"")
}

func TestVlessDNSStyleECHUsesBestEffortECHOpts(t *testing.T) {
	original := VLESS{
		Name:   "ECH-DNS节点",
		Uuid:   "12345678-1234-1234-1234-123456789abc",
		Server: "example.com",
		Port:   443,
		Query: VLESSQuery{
			Security: "tls",
			Type:     "ws",
			Sni:      "example.com",
			Ech:      "encryptedsni.com+https://dns.alidns.com/dns-query",
		},
	}

	proxy, err := buildVLESSProxy(Urls{Url: EncodeVLESSURL(original)}, OutputConfig{})
	if err != nil {
		t.Fatalf("buildVLESSProxy 失败: %v", err)
	}

	assertEqualBool(t, "ECHEnable", true, proxy.ECH_opts["enable"].(bool))
	assertEqualString(t, "ECHQueryServerName", "encryptedsni.com", proxy.ECH_opts["query-server-name"].(string))
	if _, exists := proxy.ECH_opts["config"]; exists {
		t.Fatalf("DNS 风格 ech 不应被错误映射为 config: %#v", proxy.ECH_opts)
	}

	restored := ConvertProxyToVless(proxy)
	assertEqualString(t, "RestoredTopLevelECH", "", restored.Query.Ech)
}

func TestLinkToProxy_VLESSXHTTPSkipCertFollowsSubscriptionConfig(t *testing.T) {
	vless := VLESS{
		Name:   "测试节点-VLESS-XHTTP-SkipCert",
		Uuid:   "12345678-1234-1234-1234-123456789abc",
		Server: "example.com",
		Port:   443,
		Query: VLESSQuery{
			Security:   "tls",
			Encryption: "none",
			Type:       "xhttp",
			Host:       "cdn.example.com",
			Path:       "/xhttp",
			Mode:       "stream-one",
			Extra:      `{"downloadSettings":{"path":"/download"}}`,
		},
	}

	proxy, err := buildVLESSProxy(Urls{Url: EncodeVLESSURL(vless)}, OutputConfig{Cert: true})
	if err != nil {
		t.Fatalf("buildVLESSProxy 失败: %v", err)
	}

	assertEqualString(t, "Network", "xhttp", proxy.Network)
	assertEqualBool(t, "SkipCertVerify", true, proxy.Skip_cert_verify)
	downloadSettings, ok := proxy.XHTTP_opts["download-settings"].(map[string]interface{})
	if !ok {
		t.Fatal("download-settings 不应为空")
	}
	assertEqualBool(t, "DownloadSkipCertVerify", true, downloadSettings["skip-cert-verify"].(bool))
}
