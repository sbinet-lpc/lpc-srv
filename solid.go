package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

func solidHandler(w http.ResponseWriter, r *http.Request) error {
	conn, err := net.Dial("tcp", "clrmedaq02.in2p3.fr:10000")
	if err != nil {
		return err
	}
	defer conn.Close()

	var hdr uint32
	err = binary.Read(conn, binary.LittleEndian, &hdr)
	if err != nil {
		return err
	}

	buf := make([]byte, int(hdr))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}

	var data = make(map[string]interface{})
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return err
	}

	buf, err = json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%s\n", string(buf))
	return nil
}
