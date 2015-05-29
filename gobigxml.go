package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	gviz "code.google.com/p/gographviz"
)

func main() {

	var wg sync.WaitGroup
	var path []string

	runtime.GOMAXPROCS(runtime.NumCPU())
	g := gviz.NewGraph()
	g.SetName("RouteGraph")
	g.SetDir(true)

	gchan := make(chan string)

	go func() {
		for {
			select {
			case s, open := <-gchan:
				if open {
					m := make(map[string]string)
					m["food"] = "cheese"
					m["hair"] = "brown"
					g.AddSubGraph("RouteGraph", s, m)
					wg.Done()
				}
			}
		}
	}()

	f, _ := os.Open("test.xml")
	d := xml.NewDecoder(f)

	for {
		t, _ := d.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			//log.Println(se)
			path = append(path, se.Name.Local)
			g.AddNode("RouteGraph", se.Name.Local, nil)
		case xml.EndElement:
			//log.Println(se)
			path = append(path[:len(path)-1], path[len(path):]...)
		case xml.CharData:
			v := strings.TrimSpace(string(se))
			if len(v) > 0 {
				log.Println(path, v)
				g.AddNode(path[len(path)-1], v, nil)
				g.AddEdge(path[len(path)-1], v, true, nil)
				//wg.Add(1)
				//gchan <- v
				//g.AddNode("RouteGraph", v, nil)
				//log.Println(v)
			}
		}
	}
	wg.Wait()
	close(gchan)
	fmt.Println(g.String())
}
