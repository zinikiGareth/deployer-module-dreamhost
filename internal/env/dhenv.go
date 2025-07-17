package env

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// https://api.dreamhost.com/?key=1A2B3C4D5E6F7G8H&cmd=dns-add_record&record=example.com&type=TXT&value=test123
const dhdnsapibase = "https://api.dreamhost.com/"

type DreamhostEnv struct {
	apiKey string
}

func InitDreamhostEnv() *DreamhostEnv {
	key := os.Getenv("DH_API_KEY")
	if key == "" {
		fmt.Printf("no DH_API_KEY env var set\n")
		return nil
	}
	return &DreamhostEnv{apiKey: key}
}

type DNSRecord struct {
	Record string
	Type   string
	Value  string
	Zone   string
}

type DHListResponse struct {
	Result string
	Data   []DNSRecord
}

type DHActionResponse struct {
	Result string
	Data   string
}

func (de *DreamhostEnv) DNSRecordsFor(zone string) ([]*DNSRecord, error) {
	log.Printf("finding records for zone %s\n", zone)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, de.assemble("dns-list_records"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req.WithContext(context.TODO()))
	if err != nil {
		log.Printf("have error: %v\n", err)
		io.Copy(os.Stderr, resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("error from dreamhost API: %d %v", resp.StatusCode, err)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := &DHListResponse{}
	json.Unmarshal(bs, &r)
	ret := []*DNSRecord{}
	for _, q := range r.Data {
		if q.Zone == zone || strings.HasSuffix(q.Record, zone) {
			// log.Printf("%s %s %s %s\n", q.Zone, q.Record, q.Type, q.Value)
			ret = append(ret, &q)
		}
	}

	return ret, nil
}

func (de *DreamhostEnv) FindRecord(ty, record string) (string, error) {
	all, err := de.DNSRecordsFor(record)
	if err != nil {
		return "", err
	}
	for _, r := range all {
		if r.Type == ty && r.Record == record {
			return r.Value, nil
		}
	}
	return "", nil
}

func (de *DreamhostEnv) InsertDNSRecord(record, ty, value string) error {
	// See if it's already there
	current, err := de.DNSRecordsFor(record)
	if err != nil {
		return err
	}
	// log.Printf("#records = %d\n", len(current))
	for _, r := range current {
		if r.Type == "CNAME" && r.Record == record {
			if value != r.Value && value+"." != r.Value && value != r.Value+"." {
				log.Printf("have %s but it has value %s, not %s\n", record, r.Value, value)
				panic("need to be able to change value")
			}
			// log.Printf("have %s %s\n", r.Record, r.Value)
			return nil
		}
	}

	client := &http.Client{}
	cmd := fmt.Sprintf("dns-add_record&record=%s&type=%s&value=%s", record, ty, value)
	req, err := http.NewRequest(http.MethodGet, de.assemble(cmd), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req.WithContext(context.TODO()))
	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("error from dreamhost API: %d %v", resp.StatusCode, err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("dreamhost API returned error status: %d", resp.StatusCode)
	}

	var r DHActionResponse
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(d, &r)
	if r.Result != "success" {
		return fmt.Errorf("error from dreamhost API: %s", r.Data)
	}

	return nil
}

func (de *DreamhostEnv) DeleteDNSRecord(record, ty string) error {
	// See if it's already there
	current, err := de.DNSRecordsFor(record)
	if err != nil {
		return err
	}

	// log.Printf("#records = %d\n", len(current))
	value := ""
	for _, r := range current {
		if r.Type == "CNAME" && r.Record == record {
			log.Printf("matched %s %s %s\n", r.Type, r.Record, r.Value)
			value = r.Value
			break
		}
	}

	if value == "" {
		return nil
	}

	client := &http.Client{}
	cmd := fmt.Sprintf("dns-remove_record&record=%s&type=%s&value=%s", record, ty, value)
	// log.Printf("issuing command: %s\n", cmd)
	req, err := http.NewRequest(http.MethodGet, de.assemble(cmd), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req.WithContext(context.TODO()))
	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("error from dreamhost API: %d %v", resp.StatusCode, err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("dreamhost API returned error status: %d", resp.StatusCode)
	}

	var r DHActionResponse
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// log.Printf("response was %s\n", d)
	json.Unmarshal(d, &r)
	if r.Result != "success" {
		return fmt.Errorf("error from dreamhost API: %s", r.Data)
	}

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func (de *DreamhostEnv) assemble(cmd string) string {
	return fmt.Sprintf("%s?key=%s&cmd=%s&format=json", dhdnsapibase, de.apiKey, cmd)
}

func (de *DreamhostEnv) AssertDNSRecord(zone, key, value string) error {
	// In AWS-land there is a "." at the end of a "zone id"; dreamhost does not want that
	key = strings.TrimSuffix(key, ".")

	// Didn't find it, so try and insert it
	return de.InsertDNSRecord(key, "CNAME", value)
}
