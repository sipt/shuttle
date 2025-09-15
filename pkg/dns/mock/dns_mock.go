package mockdns

import (
	"context"
	"net"
	"sync"
)

const (
	// MockIPNetwork 使用 198.0.0.0/8 网段进行 mock
	MockIPNetwork = "198.0.0.0/8"
	// MockIPStart 起始IP地址 198.0.0.1
	MockIPStart = 0xC6000002 // 198.0.0.2 in uint32
	// MockIPEnd 结束IP地址 198.255.255.254
	MockIPEnd = 0xC6FFFFFE // 198.255.255.254 in uint32
)

// DNSMock 实现DNS Mock功能
type DNSMock struct {
	// domainToIP 域名到IP的映射
	domainToIP map[string]net.IP
	// ipToDomain IP到域名的反向映射
	ipToDomain map[string]string
	// nextIP 下一个可分配的IP地址
	nextIP uint32
	// mutex 保护并发访问
	mutex sync.RWMutex
}

// NewDNSMock 创建新的DNS Mock实例
func NewDNSMock() *DNSMock {
	return &DNSMock{
		domainToIP: make(map[string]net.IP),
		ipToDomain: make(map[string]string),
		nextIP:     MockIPStart,
	}
}

// allocateIP 分配一个新的IP地址
func (m *DNSMock) allocateIP() net.IP {
	if m.nextIP > MockIPEnd {
		// 如果超出范围，从头开始重新分配
		m.nextIP = MockIPStart
	}

	// 检查IP是否已被使用
	for {
		ip := uint32ToIP(m.nextIP)
		ipStr := ip.String()

		if _, exists := m.ipToDomain[ipStr]; !exists {
			// 找到未使用的IP
			m.nextIP++
			return ip
		}

		m.nextIP++
		if m.nextIP > MockIPEnd {
			m.nextIP = MockIPStart
		}
	}
}

// uint32ToIP 将uint32转换为net.IP
func uint32ToIP(ip uint32) net.IP {
	return net.IPv4(
		byte(ip>>24),
		byte(ip>>16),
		byte(ip>>8),
		byte(ip),
	)
}

// GetOrAllocateIP 获取或分配域名对应的IP地址
func (m *DNSMock) GetOrAllocateIP(domain string) net.IP {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查是否已经分配过IP
	if ip, exists := m.domainToIP[domain]; exists {
		return ip
	}

	// 分配新的IP
	ip := m.allocateIP()
	m.domainToIP[domain] = ip
	m.ipToDomain[ip.String()] = domain

	return ip
}

// GetDomainByIP 根据IP获取对应的域名
func (m *DNSMock) GetDomainByIP(ip net.IP) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	domain, exists := m.ipToDomain[ip.String()]
	return domain, exists
}

// IsMockIP 检查IP是否为Mock IP
func (m *DNSMock) IsMockIP(ip net.IP) bool {
	if ip == nil || ip.To4() == nil {
		return false
	}

	// 检查是否在 198.0.0.0/8 网段内
	return ip.To4()[0] == 198
}

// Clear 清空所有映射
func (m *DNSMock) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.domainToIP = make(map[string]net.IP)
	m.ipToDomain = make(map[string]string)
	m.nextIP = MockIPStart
}

// ListMappings 列出所有映射关系
func (m *DNSMock) ListMappings() map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]string)
	for domain, ip := range m.domainToIP {
		result[domain] = ip.String()
	}
	return result
}

// NewMockHandle 创建DNS Mock处理器
func NewMockHandle() *DNSMock {
	return NewDNSMock()
}

// HandleQuery 处理DNS查询
func (m *DNSMock) HandleQuery(ctx context.Context, domain string, qtype uint16) ([]net.IP, error) {
	return []net.IP{m.GetOrAllocateIP(domain)}, nil
}

// ReverseLookup 反向DNS查询，根据IP查找域名
func (m *DNSMock) ReverseLookup(ip net.IP) (string, bool) {
	if !m.IsMockIP(ip) {
		return "", false
	}

	domain, exists := m.GetDomainByIP(ip)
	if !exists {
		return "", false
	}

	return domain, true
}
