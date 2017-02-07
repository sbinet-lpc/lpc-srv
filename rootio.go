package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go-hep.org/x/hep/rootio"
)

func rootioHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.New("upload").Parse(rootioPage)
		if err != nil {
			return err
		}

		var data = struct {
			Token string
			Path  string
		}{
			Token: token,
			Path:  strings.Replace(r.URL.Path+"/root-file-upload", "//", "/", -1),
		}

		err = t.Execute(w, data)
		if err != nil {
			return err
		}

	case "POST":
		r.ParseMultipartForm(500 << 20)
		f, handler, err := r.FormFile("upload-file")
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		log.Printf("inspecting %q... (from remote=%v)\n", handler.Filename, r.RemoteAddr)
		out, err := inspectROOT(f, handler.Filename)
		if err != nil {
			{
				os.MkdirAll("./rootio-files", 0755)
				now := time.Now().Unix()
				o, err := os.Create(fmt.Sprintf(
					"./rootio-files/%10d-%s", now,
					handler.Filename,
				))
				if err != nil {
					log.Printf("error creating file: %v\n", err)
				} else {
					defer o.Close()
					f.Seek(0, io.SeekStart)
					io.Copy(o, f)
				}
			}
			return err
		}

		fmt.Fprintf(w, out)

	default:
		return fmt.Errorf("invalid request %q", r.Method)
	}

	return nil
}

const rootioPage = `<html>
<head>
    <title>go-hep/rootio file inspector</title>
</head>
<body>
<h2>go-hep/rootio ROOT file inspector</h2>
<form id="rootio-form" enctype="multipart/form-data" action={{.Path}} method="post">
      <input id="rootio-file" type="file" name="upload-file" />
      <input type="hidden" name="token" value="{{.Token}}"/>
      <input type="submit" value="upload" />
</form>
<script type="text/javascript">
	document.getElementById("rootio-file").onchange = function() {
		document.getElementById("rootio-form").submit();
	}
</script>
</body>
</html>
`

func inspectROOT(r rootio.Reader, fname string) (string, error) {
	f, err := rootio.NewReader(r, fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "=== inspecting file %q...\n", fname)
	fmt.Fprintf(w, "version: %v\n", f.Version())
	for _, k := range f.Keys() {
		obj, err := k.Object()
		if err != nil {
			return "", fmt.Errorf("failed to extract key %q: %v", k.Name(), err)
		}
		switch obj := obj.(type) {
		case rootio.Tree:
			tree := obj
			fmt.Fprintf(w, "%-8s %-40s %s (entries=%d)\n", k.Class(), k.Name(), k.Title(), tree.Entries())
			displayBranches(w, tree, 2)
		default:
			fmt.Fprintf(w, "%-8s %-40s %s (cycle=%d)\n", k.Class(), k.Name(), k.Title(), k.Cycle())
		}
	}

	return string(w.Bytes()), nil
}

type brancher interface {
	Branches() []rootio.Branch
}

func displayBranches(w io.Writer, bres brancher, indent int) {
	branches := bres.Branches()
	for _, b := range branches {
		fmt.Fprintf(w, "%s%-20s %-20q %v\n", strings.Repeat(" ", indent), b.Name(), b.Title(), b.Class())
		displayBranches(w, b, indent+2)
	}
}
