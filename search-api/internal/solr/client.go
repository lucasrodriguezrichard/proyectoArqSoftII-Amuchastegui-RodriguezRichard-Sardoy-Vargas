package solr

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

type Client struct {
    base   string // e.g. http://localhost:8983/solr
    core   string // e.g. reservations
    httpc  *http.Client
}

func New(base, core string) *Client {
    return &Client{
        base:  stringsTrimRight(base, "/"),
        core:  core,
        httpc: &http.Client{Timeout: 5 * time.Second},
    }
}

func stringsTrimRight(s, cutset string) string {
    for len(s) > 0 && stringsHasSuffix(s, cutset) {
        s = s[:len(s)-len(cutset)]
    }
    return s
}

func stringsHasSuffix(s, suffix string) bool {
    if len(suffix) == 0 { return false }
    if len(s) < len(suffix) { return false }
    return s[len(s)-len(suffix):] == suffix
}

// Update API payloads
type addDoc struct{
    Add struct{
        Doc any `json:"doc"`
        CommitWithin int `json:"commitWithin,omitempty"`
    } `json:"add"`
}

// Index/add a document
func (c *Client) Index(doc any) error {
    payload := addDoc{}
    payload.Add.Doc = doc
    payload.Add.CommitWithin = 1000
    return c.postJSON("/update", payload)
}

// Update is equivalent to Index in Solr (atomic/partial updates omitted for brevity)
func (c *Client) Update(doc any) error {
    return c.Index(doc)
}

// Delete by id
func (c *Client) Delete(id string) error {
    body := map[string]any{"delete": map[string]string{"id": id}, "commit": true}
    return c.postJSON("/update", body)
}

// Search with basic params
type SearchResponse struct{
    Response struct{
        NumFound int `json:"numFound"`
        Start    int `json:"start"`
        Docs     []map[string]any `json:"docs"`
    } `json:"response"`
}

func (c *Client) Search(q string, start, rows int, filters map[string]string, sort string) (*SearchResponse, error) {
    values := url.Values{}
    if q == "" { q = "*:*" }
    values.Set("q", q)
    values.Set("start", fmt.Sprintf("%d", start))
    values.Set("rows", fmt.Sprintf("%d", rows))
    for k, v := range filters {
        if v == "" { continue }
        values.Add("fq", fmt.Sprintf("%s:%s", k, v))
    }
    if sort != "" { values.Set("sort", sort) }

    u := fmt.Sprintf("%s/%s/select?%s", c.base, c.core, values.Encode())
    req, _ := http.NewRequest(http.MethodGet, u, nil)
    resp, err := c.httpc.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("solr search status %d: %s", resp.StatusCode, string(b))
    }
    var sr SearchResponse
    if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil { return nil, err }
    return &sr, nil
}

func (c *Client) postJSON(path string, payload any) error {
    b, _ := json.Marshal(payload)
    u := fmt.Sprintf("%s/%s%s", c.base, c.core, path)
    req, _ := http.NewRequest(http.MethodPost, u, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpc.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("solr update status %d: %s", resp.StatusCode, string(body))
    }
    return nil
}

