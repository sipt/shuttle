package mock

import (
	"context"
	"net"
	"testing"
)

func TestDNSMock_GetOrAllocateIP(t *testing.T) {
	mock := NewDNSMock()

	// 测试第一次分配
	domain1 := "example.com"
	ip1 := mock.GetOrAllocateIP(domain1)

	if ip1 == nil {
		t.Fatal("分配的IP不应该为空")
	}

	if !mock.IsMockIP(ip1) {
		t.Fatalf("分配的IP %s 不在Mock网段内", ip1.String())
	}

	// 测试相同域名返回相同IP
	ip1Again := mock.GetOrAllocateIP(domain1)
	if !ip1.Equal(ip1Again) {
		t.Fatalf("相同域名应该返回相同IP，期望 %s，实际 %s", ip1.String(), ip1Again.String())
	}

	// 测试不同域名分配不同IP
	domain2 := "google.com"
	ip2 := mock.GetOrAllocateIP(domain2)

	if ip1.Equal(ip2) {
		t.Fatalf("不同域名应该分配不同IP，都分配了 %s", ip1.String())
	}
}

func TestDNSMock_GetDomainByIP(t *testing.T) {
	mock := NewDNSMock()

	domain := "example.com"
	ip := mock.GetOrAllocateIP(domain)

	// 测试反向查询
	foundDomain, exists := mock.GetDomainByIP(ip)
	if !exists {
		t.Fatal("应该能够找到对应的域名")
	}

	if foundDomain != domain {
		t.Fatalf("反向查询结果不正确，期望 %s，实际 %s", domain, foundDomain)
	}

	// 测试查询不存在的IP
	nonExistentIP := net.ParseIP("198.1.1.1")
	_, exists = mock.GetDomainByIP(nonExistentIP)
	if exists {
		t.Fatal("不应该找到不存在的IP对应的域名")
	}
}

func TestDNSMock_IsMockIP(t *testing.T) {
	mock := NewDNSMock()

	// 测试Mock IP
	mockIP := net.ParseIP("198.1.2.3")
	if !mock.IsMockIP(mockIP) {
		t.Fatalf("IP %s 应该被识别为Mock IP", mockIP.String())
	}

	// 测试非Mock IP
	realIP := net.ParseIP("8.8.8.8")
	if mock.IsMockIP(realIP) {
		t.Fatalf("IP %s 不应该被识别为Mock IP", realIP.String())
	}

	// 测试nil IP
	if mock.IsMockIP(nil) {
		t.Fatal("nil IP不应该被识别为Mock IP")
	}
}

func TestDNSMock_Clear(t *testing.T) {
	mock := NewDNSMock()

	// 分配一些IP
	mock.GetOrAllocateIP("example.com")
	mock.GetOrAllocateIP("google.com")

	mappings := mock.ListMappings()
	if len(mappings) != 2 {
		t.Fatalf("应该有2个映射，实际有 %d 个", len(mappings))
	}

	// 清空
	mock.Clear()

	mappings = mock.ListMappings()
	if len(mappings) != 0 {
		t.Fatalf("清空后应该没有映射，实际有 %d 个", len(mappings))
	}
}

func TestDNSMock_IPAllocation(t *testing.T) {
	mock := NewDNSMock()

	// 测试IP分配是否连续
	domain1 := "test1.com"
	domain2 := "test2.com"

	ip1 := mock.GetOrAllocateIP(domain1)
	ip2 := mock.GetOrAllocateIP(domain2)

	// 转换为uint32比较
	ip1Uint := ipToUint32(ip1)
	ip2Uint := ipToUint32(ip2)

	if ip2Uint != ip1Uint+1 {
		t.Fatalf("IP分配应该是连续的，ip1: %s (%d), ip2: %s (%d)",
			ip1.String(), ip1Uint, ip2.String(), ip2Uint)
	}
}

func TestNewMockHandle(t *testing.T) {
	// 创建Mock处理器
	handle := NewMockHandle(nil)

	ctx := context.Background()
	domain := "example.com"

	// 测试DNS解析
	result := handle(ctx, domain)

	if result == nil {
		t.Fatal("DNS解析结果不应该为空")
	}

	if result.Domain != domain {
		t.Fatalf("域名不匹配，期望 %s，实际 %s", domain, result.Domain)
	}

	if result.Typ != "mock" {
		t.Fatalf("类型不匹配，期望 mock，实际 %s", result.Typ)
	}

	if len(result.IP) != 1 {
		t.Fatalf("应该返回1个IP，实际返回 %d 个", len(result.IP))
	}

	// 验证IP在Mock网段内
	mockInstance := NewDNSMock()
	if !mockInstance.IsMockIP(result.CurrentIP) {
		t.Fatalf("返回的IP %s 不在Mock网段内", result.CurrentIP.String())
	}
}

func TestDNSMock_ReverseLookup(t *testing.T) {
	mock := NewDNSMock()

	domain := "example.com"
	ip := mock.GetOrAllocateIP(domain)

	// 测试反向查询
	result := mock.ReverseLookup(ip)

	if result == nil {
		t.Fatal("反向查询结果不应该为空")
	}

	if result.Domain != domain {
		t.Fatalf("反向查询域名不匹配，期望 %s，实际 %s", domain, result.Domain)
	}

	if !result.CurrentIP.Equal(ip) {
		t.Fatalf("反向查询IP不匹配，期望 %s，实际 %s", ip.String(), result.CurrentIP.String())
	}

	// 测试查询非Mock IP
	realIP := net.ParseIP("8.8.8.8")
	result = mock.ReverseLookup(realIP)
	if result != nil {
		t.Fatal("查询非Mock IP应该返回nil")
	}
}

// 辅助函数：将IP转换为uint32
func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
