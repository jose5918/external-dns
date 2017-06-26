/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"fmt"
	"os"
	"testing"

	"strconv"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/kubernetes-incubator/external-dns/endpoint"
	"github.com/kubernetes-incubator/external-dns/plan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDnsimpleZonesService struct{}

func (m *mockDnsimpleZonesService) ListZones(accountID string, options *dnsimple.ZoneListOptions) (*dnsimple.ZonesResponse, error) {
	firstZone := dnsimple.Zone{
		ID:        1,
		AccountID: 12345,
		Name:      "example.com",
	}
	secondZone := dnsimple.Zone{
		ID:        2,
		AccountID: 54321,
		Name:      "example-beta.com",
	}
	zones := []dnsimple.Zone{firstZone, secondZone}
	listResponse := dnsimple.ZonesResponse{
		Response: dnsimple.Response{},
		Data:     zones,
	}
	return &listResponse, nil
}
func (m *mockDnsimpleZonesService) GetZone(accountID string, zoneName string) (*dnsimple.ZoneResponse, error) {
	return &dnsimple.ZoneResponse{}, nil
}

func (m *mockDnsimpleZonesService) ListRecords(accountID string, zoneID string, options *dnsimple.ZoneRecordListOptions) (*dnsimple.ZoneRecordsResponse, error) {
	firstRecord := dnsimple.ZoneRecord{
		ID:       2,
		ZoneID:   "example.com",
		ParentID: 0,
		Name:     "example",
		Content:  "ns1.dnsimple.com",
		TTL:      3600,
		Priority: 0,
		Type:     "SOA",
	}
	secondRecord := dnsimple.ZoneRecord{
		ID:       1,
		ZoneID:   "example.com",
		ParentID: 0,
		Name:     "example-beta",
		Content:  "127.0.0.1",
		TTL:      3600,
		Priority: 0,
		Type:     "A",
	}
	records := []dnsimple.ZoneRecord{firstRecord, secondRecord}
	listResponse := dnsimple.ZoneRecordsResponse{
		Response: dnsimple.Response{},
		Data:     records,
	}
	return &listResponse, nil
}

func (m *mockDnsimpleZonesService) CreateRecord(accountID string, zoneID string, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesService) DeleteRecord(accountID string, zoneID string, recordID int) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesService) UpdateRecord(accountID string, zoneID string, recordID int, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

type mockDnsimpleZonesServiceRecordsFail struct{}

func (m *mockDnsimpleZonesServiceRecordsFail) ListZones(accountID string, options *dnsimple.ZoneListOptions) (*dnsimple.ZonesResponse, error) {
	firstZone := dnsimple.Zone{
		ID:        1,
		AccountID: 12345,
		Name:      "example.com",
	}
	secondZone := dnsimple.Zone{
		ID:        2,
		AccountID: 54321,
		Name:      "example.com",
	}
	zones := []dnsimple.Zone{firstZone, secondZone}
	listResponse := dnsimple.ZonesResponse{
		Response: dnsimple.Response{},
		Data:     zones,
	}
	return &listResponse, nil
}
func (m *mockDnsimpleZonesServiceRecordsFail) GetZone(accountID string, zoneName string) (*dnsimple.ZoneResponse, error) {
	return &dnsimple.ZoneResponse{}, nil
}

func (m *mockDnsimpleZonesServiceRecordsFail) ListRecords(accountID string, zoneID string, options *dnsimple.ZoneRecordListOptions) (*dnsimple.ZoneRecordsResponse, error) {
	return &dnsimple.ZoneRecordsResponse{}, fmt.Errorf("Failed to list records")
}

func (m *mockDnsimpleZonesServiceRecordsFail) CreateRecord(accountID string, zoneID string, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesServiceRecordsFail) DeleteRecord(accountID string, zoneID string, recordID int) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesServiceRecordsFail) UpdateRecord(accountID string, zoneID string, recordID int, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

type mockDnsimpleZonesServiceCreateFail struct{}

func (m *mockDnsimpleZonesServiceCreateFail) ListZones(accountID string, options *dnsimple.ZoneListOptions) (*dnsimple.ZonesResponse, error) {
	firstZone := dnsimple.Zone{
		ID:        1,
		AccountID: 12345,
		Name:      "example.com",
	}
	secondZone := dnsimple.Zone{
		ID:        2,
		AccountID: 54321,
		Name:      "example.com",
	}
	zones := []dnsimple.Zone{firstZone, secondZone}
	listResponse := dnsimple.ZonesResponse{
		Response: dnsimple.Response{},
		Data:     zones,
	}
	return &listResponse, nil
}
func (m *mockDnsimpleZonesServiceCreateFail) GetZone(accountID string, zoneName string) (*dnsimple.ZoneResponse, error) {
	return &dnsimple.ZoneResponse{}, nil
}

func (m *mockDnsimpleZonesServiceCreateFail) ListRecords(accountID string, zoneID string, options *dnsimple.ZoneRecordListOptions) (*dnsimple.ZoneRecordsResponse, error) {
	firstRecord := dnsimple.ZoneRecord{
		ID:       2,
		ZoneID:   "example.com",
		ParentID: 0,
		Name:     "",
		Content:  "ns1.dnsimple.com",
		TTL:      3600,
		Priority: 0,
		Type:     "SOA",
	}
	secondRecord := dnsimple.ZoneRecord{
		ID:       1,
		ZoneID:   "example-alpha.com",
		ParentID: 0,
		Name:     "",
		Content:  "127.0.0.1",
		TTL:      3600,
		Priority: 0,
		Type:     "A",
	}
	records := []dnsimple.ZoneRecord{firstRecord, secondRecord}
	listResponse := dnsimple.ZoneRecordsResponse{
		Response: dnsimple.Response{},
		Data:     records,
	}
	return &listResponse, nil
}

func (m *mockDnsimpleZonesServiceCreateFail) CreateRecord(accountID string, zoneID string, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, fmt.Errorf("Failed to create record")
}

func (m *mockDnsimpleZonesServiceCreateFail) DeleteRecord(accountID string, zoneID string, recordID int) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesServiceCreateFail) UpdateRecord(accountID string, zoneID string, recordID int, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

type mockDnsimpleZonesServiceDeleteFail struct{}

func (m *mockDnsimpleZonesServiceDeleteFail) ListZones(accountID string, options *dnsimple.ZoneListOptions) (*dnsimple.ZonesResponse, error) {
	firstZone := dnsimple.Zone{
		ID:        1,
		AccountID: 12345,
		Name:      "example.com",
	}
	secondZone := dnsimple.Zone{
		ID:        2,
		AccountID: 54321,
		Name:      "example.com",
	}
	zones := []dnsimple.Zone{firstZone, secondZone}
	listResponse := dnsimple.ZonesResponse{
		Response: dnsimple.Response{},
		Data:     zones,
	}
	return &listResponse, nil
}
func (m *mockDnsimpleZonesServiceDeleteFail) GetZone(accountID string, zoneName string) (*dnsimple.ZoneResponse, error) {
	return &dnsimple.ZoneResponse{}, nil
}

func (m *mockDnsimpleZonesServiceDeleteFail) ListRecords(accountID string, zoneID string, options *dnsimple.ZoneRecordListOptions) (*dnsimple.ZoneRecordsResponse, error) {
	firstRecord := dnsimple.ZoneRecord{
		ID:       2,
		ZoneID:   "example.com",
		ParentID: 0,
		Name:     "",
		Content:  "ns1.dnsimple.com",
		TTL:      3600,
		Priority: 0,
		Type:     "SOA",
	}
	secondRecord := dnsimple.ZoneRecord{
		ID:       1,
		ZoneID:   "example-alpha.com",
		ParentID: 0,
		Name:     "",
		Content:  "127.0.0.1",
		TTL:      3600,
		Priority: 0,
		Type:     "A",
	}
	records := []dnsimple.ZoneRecord{firstRecord, secondRecord}
	listResponse := dnsimple.ZoneRecordsResponse{
		Response: dnsimple.Response{},
		Data:     records,
	}
	return &listResponse, nil
}

func (m *mockDnsimpleZonesServiceDeleteFail) CreateRecord(accountID string, zoneID string, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesServiceDeleteFail) DeleteRecord(accountID string, zoneID string, recordID int) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, fmt.Errorf("Failed to delete record")
}

func (m *mockDnsimpleZonesServiceDeleteFail) UpdateRecord(accountID string, zoneID string, recordID int, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

type mockDnsimpleZonesServiceUpdateFail struct{}

func (m *mockDnsimpleZonesServiceUpdateFail) ListZones(accountID string, options *dnsimple.ZoneListOptions) (*dnsimple.ZonesResponse, error) {
	firstZone := dnsimple.Zone{
		ID:        1,
		AccountID: 12345,
		Name:      "example.com",
	}
	secondZone := dnsimple.Zone{
		ID:        2,
		AccountID: 54321,
		Name:      "example.com",
	}
	zones := []dnsimple.Zone{firstZone, secondZone}
	listResponse := dnsimple.ZonesResponse{
		Response: dnsimple.Response{},
		Data:     zones,
	}
	return &listResponse, nil
}
func (m *mockDnsimpleZonesServiceUpdateFail) GetZone(accountID string, zoneName string) (*dnsimple.ZoneResponse, error) {
	return &dnsimple.ZoneResponse{}, nil
}

func (m *mockDnsimpleZonesServiceUpdateFail) ListRecords(accountID string, zoneID string, options *dnsimple.ZoneRecordListOptions) (*dnsimple.ZoneRecordsResponse, error) {
	firstRecord := dnsimple.ZoneRecord{
		ID:       2,
		ZoneID:   "example.com",
		ParentID: 0,
		Name:     "",
		Content:  "ns1.dnsimple.com",
		TTL:      3600,
		Priority: 0,
		Type:     "SOA",
	}
	secondRecord := dnsimple.ZoneRecord{
		ID:       1,
		ZoneID:   "example-alpha.com",
		ParentID: 0,
		Name:     "",
		Content:  "127.0.0.1",
		TTL:      3600,
		Priority: 0,
		Type:     "A",
	}
	records := []dnsimple.ZoneRecord{firstRecord, secondRecord}
	listResponse := dnsimple.ZoneRecordsResponse{
		Response: dnsimple.Response{},
		Data:     records,
	}
	return &listResponse, nil
}

func (m *mockDnsimpleZonesServiceUpdateFail) CreateRecord(accountID string, zoneID string, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesServiceUpdateFail) DeleteRecord(accountID string, zoneID string, recordID int) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, nil
}

func (m *mockDnsimpleZonesServiceUpdateFail) UpdateRecord(accountID string, zoneID string, recordID int, recordAttributes dnsimple.ZoneRecord) (*dnsimple.ZoneRecordResponse, error) {
	return &dnsimple.ZoneRecordResponse{}, fmt.Errorf("Failed to update record")
}

func TestNewDnsimpleProvider(t *testing.T) {
	os.Setenv("DNSIMPLE_OAUTH", "xxxxxxxxxxxxxxxxxxxxxxxxxx")
	_, err := NewDnsimpleProvider("example.com", true)
	if err == nil {
		t.Errorf("Expected to fail new provider on bad token")
	}
	os.Unsetenv("DNSIMPLE_OAUTH")
	if err == nil {
		t.Errorf("Expected to fail new provider on empty token")
	}
}

// func TestDnsimple(t *testing.T) {
// 	mockDNS := mocks.DnsimpleZoneServiceInterface{}
// 	testDnsimpleProviderZones
// 	t.Run("Zones", testDnsimpleProviderZones)
// }

// func testDnsimpleProviderZones(t *testing.T, mock.DnsimpleZoneServiceInterface) {
// 	firstZone := dnsimple.Zone{
// 		ID:        1,
// 		AccountID: 12345,
// 		Name:      "example.com",
// 	}
// 	secondZone := dnsimple.Zone{
// 		ID:        2,
// 		AccountID: 54321,
// 		Name:      "example-beta.com",
// 	}
// 	zones := []dnsimple.Zone{firstZone, secondZone}
// 	listResponse := dnsimple.ZonesResponse{
// 		Response: dnsimple.Response{},
// 		Data:     zones,
// 	}
// 	mockDNS := &mocks.DnsimpleZoneServiceInterface{}
// 	mockDNS.On("ListZones", "1", &dnsimple.ZoneListOptions{}).Return(&listResponse, nil)
// 	mockDNS.On("ListZones", "2", &dnsimple.ZoneListOptions{}).Return(nil, fmt.Errorf("Account ID not found"))
// 	provider := dnsimpleProvider{
// 		client:    mockDNS,
// 		accountID: "1",
// 	}
// 	result, err := provider.Zones()
// 	assert.Nil(t, err)
// 	validateDnsimpleZones(t, result, zones)

// 	provider.accountID = "2"
// 	result, err = provider.Zones()
// 	assert.NotNil(t, err)
// }
func TestForReal(t *testing.T) {
	os.Setenv("DNSIMPLE_OAUTH", "")
	provider, err := NewDnsimpleProvider("kube-external-test.io", false)
	if err != nil {
		t.Errorf("Expected to fail new provider on bad token: %v", err)
	}
	os.Setenv("DNSIMPLE_SANDBOX", "true")
	changes := &plan.Changes{}
	changes.Delete = []*endpoint.Endpoint{{DNSName: "example-beta.kube-external-test.io", Target: "target"}}
	err = provider.ApplyChanges(changes)
	if err != nil {
		t.Errorf("Failed to apply changes: %v", err)
	}

}
func TestDnsimpleProvider_Records(t *testing.T) {
	provider := &dnsimpleProvider{
		client: &mockDnsimpleZonesService{},
	}
	_, err := provider.Records()
	if err != nil {
		t.Errorf("Couldn't get records: %v", err)
	}
	provider.client = &mockDnsimpleZonesServiceRecordsFail{}
	_, err = provider.Records()
	if err == nil {
		t.Errorf("Expected records failure")
	}
}

func TestDnsimpleProvider_CreateRecords(t *testing.T) {
	provider := &dnsimpleProvider{
		client: &mockDnsimpleZonesService{},
	}
	endpoints := []*endpoint.Endpoint{
		{DNSName: "new.example.com", Target: "target"},
	}
	err := provider.CreateRecords(endpoints)
	if err != nil {
		t.Errorf("Could not create record: %v", err)
	}
	endpoints = []*endpoint.Endpoint{}
	err = provider.CreateRecords(endpoints)
	if err != nil {
		t.Errorf("Failed empty case: %v", err)
	}
	endpoints = []*endpoint.Endpoint{
		{DNSName: "new.example.com", Target: "target"},
	}
	provider.client = &mockDnsimpleZonesServiceCreateFail{}
	err = provider.CreateRecords(endpoints)
	if err == nil {
		t.Errorf("Expected create records failure")
	}
}

func TestDnsimpleProvider_DeleteRecords(t *testing.T) {
	provider := &dnsimpleProvider{
		client: &mockDnsimpleZonesService{},
	}
	provider.dryRun = false
	endpoints := []*endpoint.Endpoint{
		{DNSName: "example-beta.example.com", Target: "127.0.0.1"},
	}
	err := provider.DeleteRecords(endpoints)
	if err != nil {
		t.Errorf("Could not delete record: %v", err)
	}
	endpoints = []*endpoint.Endpoint{
		{DNSName: "new", Target: "target"},
	}
	provider.client = &mockDnsimpleZonesServiceDeleteFail{}
	err = provider.DeleteRecords(endpoints)
	if err == nil {
		t.Errorf("Expected delete records failure")
	}
	provider.client = &mockDnsimpleZonesService{}
	endpoints = []*endpoint.Endpoint{}
	err = provider.CreateRecords(endpoints)
	if err != nil {
		t.Errorf("Failed empty case: %v", err)
	}
}

func TestDnsimpleProvider_UpdateRecords(t *testing.T) {
	provider := &dnsimpleProvider{
		client: &mockDnsimpleZonesService{},
	}
	endpoints := []*endpoint.Endpoint{
		{DNSName: "example.example.com", Target: "127.0.0.2"},
	}
	err := provider.UpdateRecords(endpoints)
	if err != nil {
		t.Errorf("Failed to update records: %v", err)
	}
	provider.client = &mockDnsimpleZonesServiceUpdateFail{}
	err = provider.UpdateRecords(endpoints)
	if err == nil {
		t.Errorf("Expected update records failure")
	}
	provider.client = &mockDnsimpleZonesService{}
	endpoints = []*endpoint.Endpoint{}
	err = provider.CreateRecords(endpoints)
	if err != nil {
		t.Errorf("Failed empty case: %v", err)
	}
}

func TestDnsimpleProvider_ApplyChanges(t *testing.T) {
	changes := &plan.Changes{}
	provider := &dnsimpleProvider{
		client: &mockDnsimpleZonesService{},
	}
	changes.Create = []*endpoint.Endpoint{{DNSName: "example-beta.example.com", Target: "target"}}
	changes.Delete = []*endpoint.Endpoint{{DNSName: "example-beta.example.com", Target: "127.0.0.1"}}
	changes.UpdateNew = []*endpoint.Endpoint{{DNSName: "example.example.com", Target: "target-new"}}
	err := provider.ApplyChanges(changes)
	if err != nil {
		t.Errorf("Failed to apply changes: %v", err)
	}
}

func TestDnsimpleSuitableZone(t *testing.T) {
	provider := &dnsimpleProvider{
		client: &mockDnsimpleZonesService{},
	}
	zones, _ := provider.Zones()
	zone := dnsimpleSuitableZone("example-beta.example.com", zones)
	assert.Equal(t, zone.Name, "example.com")
}

func validateDnsimpleZones(t *testing.T, zones map[string]dnsimple.Zone, expected []dnsimple.Zone) {
	require.Len(t, zones, len(expected))

	for _, e := range expected {
		assert.Equal(t, zones[strconv.Itoa(e.ID)].Name, e.Name)
	}
}
