package env

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// https://api.dreamhost.com/?key=1A2B3C4D5E6F7G8H&cmd=dns-add_record&record=example.com&type=TXT&value=test123
const dhdnsapibase = "https://api.dreamhost.com/"

type DreamhostEnv struct {
	apiKey string
}

func InitDreamhostEnv() *DreamhostEnv {
	key := os.Getenv("DH_API_KEY")
	if key == "" {
		fmt.Printf("no DH_API_KEY env var set")
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

type DHResponse struct {
	Result string
	Data   []DNSRecord
}

func (de *DreamhostEnv) DNSRecordsFor(zone string) ([]*DNSRecord, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, de.assemble("dns-list_records"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req.WithContext(context.TODO()))
	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("error from dreamhost API: %d %v", resp.StatusCode, err)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := &DHResponse{}
	json.Unmarshal(bs, &r)
	log.Printf("%s len = %d\n", r.Result, len(r.Data))
	ret := []*DNSRecord{}
	for _, q := range r.Data {
		if q.Zone == zone {
			log.Printf("%s %s %s %s\n", q.Zone, q.Record, q.Type, q.Value)
			ret = append(ret, &q)
		}
	}

	return ret, nil
}

func (de *DreamhostEnv) assemble(cmd string) string {
	return fmt.Sprintf("%s?key=%s&cmd=%s&format=json", dhdnsapibase, de.apiKey, cmd)
}
