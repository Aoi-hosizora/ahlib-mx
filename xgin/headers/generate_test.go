package headers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// _ = generate
	err := generate()
	if err != nil {
		log.Fatalln(err)
	}
}

func generate() error {
	stdList, nstdList, err := getList()
	if err != nil {
		return err
	}

	sb := bytes.Buffer{}
	sb.WriteString(`package headers

// Headers are referred from https://github.com/go-http-utils/headers and https://en.wikipedia.org/wiki/List_of_HTTP_header_fields.

// Standard header fields.
const (
` + strings.Join(stdList, "\n") + `
)

// Common non-standard header fields.
const (` + strings.Join(nstdList, "\n") + `
)
`)
	err = formatAndWrite(sb.Bytes(), "_generate/headers.go")
	if err != nil {
		return fmt.Errorf("code formatAndWrite: %w", err)
	}
	return nil
}

func formatAndWrite(bs []byte, filename string) error {
	bs, err := format.Source(bs)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(filename), 0644)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func getList() (stdList []string, nstdList []string, err error) {
	url := "https://en.wikipedia.org/wiki/List_of_HTTP_header_fields"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("client: Do: %w", err)
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	headerLists := make([][]string, 4)
	re1 := regexp.MustCompile(`<h3.*?>.+?</h3.*?>\n?<table.*?>[\s\S]+?</table>`)
	re2 := regexp.MustCompile(`<tr.*?>\n*<td.*?>\n*([\s\S]+?)\n*</td>\n*<td.*?>`)
	matches1 := re1.FindAll(bs, -1)
	if len(matches1) != 4 {
		return nil, nil, errors.New("len(matches1) != 4")
	}
	for i, match := range matches1 {
		matches := re2.FindAllSubmatch(match, -1)
		for _, m := range matches {
			line := string(m[1])
			if strings.HasPrefix(line, "<a ") {
				line = regexp.MustCompile(`<a.*?>([\s\S]+?)</a>`).FindStringSubmatch(line)[1]
			}
			line = regexp.MustCompile(`<sup.+?</sup>`).ReplaceAllString(line, "")
			line = regexp.MustCompile(`<span.*?class="anchor.+?</span>`).ReplaceAllString(line, "")
			line = strings.ReplaceAll(strings.ReplaceAll(line, `<span class="nowrap">`, ""), `</span>`, "")
			line = strings.ReplaceAll(strings.ReplaceAll(line, `<p>`, ""), `</p>`, "")
			line = strings.ReplaceAll(strings.ReplaceAll(line, `<br />`, ","), "\n", ",")
			for _, item := range strings.Split(line, ",") {
				item = strings.TrimSpace(item)
				if item != "" {
					headerLists[i] = append(headerLists[i], item)
				}
			}
		}
	}
	headerLists[1] = append(headerLists[1], "X-Real-IP")
	headerLists[3] = append(headerLists[3], "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset")

	stdHeaderMap := make(map[string][]string, len(headerLists[0])+len(headerLists[2]))
	nstdHeaderMap := make(map[string][]string, len(headerLists[1])+len(headerLists[3]))
	for i, list := range headerLists {
		for _, item := range list {
			switch i {
			case 0: // std req
				stdHeaderMap[item] = append(stdHeaderMap[item], "requests")
			case 1: // nstd req
				nstdHeaderMap[item] = append(nstdHeaderMap[item], "requests")
			case 2: // std resp
				stdHeaderMap[item] = append(stdHeaderMap[item], "responses")
			case 3: // nstd resp
				nstdHeaderMap[item] = append(nstdHeaderMap[item], "responses")
			}
		}
	}

	stdList = make([]string, 0, len(stdHeaderMap))
	nstdList = make([]string, 0, len(nstdHeaderMap))
	for item, comments := range stdHeaderMap {
		item = fmt.Sprintf(`%s = "%s" // Used in %s`, xstring.PascalCase(item), item, strings.Join(comments, ", "))
		stdList = append(stdList, item)
	}
	for item, comments := range nstdHeaderMap {
		item = fmt.Sprintf(`%s = "%s" // Used in %s`, xstring.PascalCase(item), item, strings.Join(comments, ", "))
		nstdList = append(nstdList, item)
	}
	sort.Strings(stdList)
	sort.Strings(nstdList)

	return stdList, nstdList, nil
}
