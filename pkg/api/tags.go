package api

import (
    "encoding/json"
	"encoding/xml"
	"net/http"
    "os/exec"
	"fmt"
)

type TagInfo struct {
	Table string 	`xml:"-"`
	Name string 	`xml:"name,attr"`
	Type string 	`xml:"type,attr"`
	Writable bool 	`xml:"writable,attr"`
	Desc[] DescInfo `xml:"desc"`
}

type DescInfo struct {
	Value string 	`xml:",chardata"`
	Lang string 	`xml:"lang,attr"`
}

func (tag TagInfo) MarshalJSON() ([]byte, error) {
	desc := make(map[string]string)
	
	for _, d := range tag.Desc {
		desc[d.Lang] = d.Value
	}
	
    return json.Marshal(map[string]interface{}{
		"writable": tag.Writable,
		"path": tag.Table + ":" + tag.Name,
		"group": tag.Table,
		"type": tag.Type,
		"description": desc,
	})
}

func Tags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	cmd := exec.Command("exiftool", "-listx")
	
	stdout, _ := cmd.StdoutPipe()
	
    if err := cmd.Start(); err != nil {
       fmt.Printf("", err)
    }
	
	decoder := xml.NewDecoder(stdout)
	encoder := json.NewEncoder(w)
	
	var inElement string
	var tableName string
	var tagInfo TagInfo
	
	writeComma := false
	
	fmt.Fprintf(w, `{"tags": [`)
	
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			
			if inElement == "table" {
				tableName = se.Attr[0].Value
			
			} else if inElement == "tag" {
				decoder.DecodeElement(&tagInfo, &se)
				tagInfo.Table = tableName
				
				if writeComma {
					fmt.Fprintf(w, ",")
				}
				
				encoder.Encode(tagInfo)	
				
				writeComma = true
			}
			
		default:
		}
	}
	
	fmt.Fprintf(w, `]}`)
	
	cmd.Wait()
}