// Copyright 2018 hybris. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// simple wrapper around http.Fileserver and http.ServeFile
package serve

import (
    "os"
    "strconv"
    "strings"
    "net/http"
    "path/filepath"
    "github.com/kpango/glg"
)


func Serve(path_tainted string, ip_tainted string, port_tainted int) {
    if (os.Geteuid() == 0) {
        glg.Warn("running as root")
    }

    path, err_path := validatePath(path_tainted)
    ip, err_ip := validateAddress(ip_tainted)
    port, err_port := validatePort(port_tainted)

    listen := strings.Join([]string{ip,strconv.Itoa(port)}, ":")

    if (err_path == nil && err_ip == nil && err_port == nil) {
        fi, _ := os.Stat(path)
        switch mode := fi.Mode(); {
            case mode.IsDir():
                serveDirectory(path,listen)
                break
            case mode.IsRegular():
                serveFile(path,listen)
                break
        }
    }
}

// ServeDirectory will start a http.FileServer for the given path on the given listen address
// expects `path` to be a valid directory
func serveDirectory(path,listen string) {
    glg.Infof("serving directory %s on %s",path,listen)
    if (strings.Compare(path,"/") == 0) {
        glg.Warn("you are serving your whole filesystem right now, use with care")
    }
    panic(http.ListenAndServe(listen, loggingHandler(http.FileServer(http.Dir(path)))))
}

// custom handler to achieve access logging for http.FileServer
func loggingHandler(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        glg.Infof("Access: from: %s  method: %s  full_path: %s%s",r.RemoteAddr,r.Method,r.Host,r.URL.Path)
        h.ServeHTTP(w, r)
    })
}

// ServeFile will start http.ListenAndServe and answer request on "/" with the file from the given path
// expects `path` to be a valid file
func serveFile(path,listen string) {
    glg.Infof("serving file '%s' on '%s'",path,listen)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        glg.Infof("Access: from: %s  method: %s  full_path: %s%s",r.RemoteAddr,r.Method,r.Host,r.URL.Path)
        var header_string = "attachment; filename="
        header_string += filepath.Base(path)
        w.Header().Set("Content-Disposition", header_string)
        http.ServeFile(w, r, path)
    })
    panic(http.ListenAndServe(listen,nil))
}
