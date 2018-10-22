// Copyright 2018 hybris. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package serve

import (
    "os"
    "net"
    "errors"
    "strings"
    "path/filepath"
    "github.com/kpango/glg"
)

// ValidatePath validates if the path exists and returns the full absolute path of the given path.
// if no path was provided (e.g. empty string), it will return the current working directory.
// exit with error if path is not valid or does not exist
func validatePath(path_tainted string) (string, error) {
    if path_tainted == "" {
        glg.Info("no directory specified, serving current working directory")
        return filepath.Abs(filepath.Dir(os.Args[0]))
    }

    if _, err := os.Stat(path_tainted); os.IsNotExist(err) {
        glg.Errorf("file or directory does not exist: %s",path_tainted)
        return "", err
    }

    if filepath.IsAbs(path_tainted) {
        return path_tainted, nil
    }

    return filepath.Abs(path_tainted)
}

// ValidateAddress will check if a provided string is a valid ipv4/6 address and return it as string.
// if will return ipv6 addresses with surrounding brackets.(Example: "::" => "[::]")
func validateAddress(address_tainted string) (string, error) {
    if address_tainted == "" {
        return "", nil
    }

    var address = net.ParseIP(address_tainted)
    if address == nil {
        glg.Errorf("%s is not a valid ipv4 or ipv6 address",address_tainted)
        return "", errors.New("not a valid ipv4 or ipv6 address")
    }

    if (address.To4() == nil) {
        return strings.Join([]string{"[","]"},address.String()), nil
    }

    return address.String(), nil
}

// ValidatePort will check if a given port is in valid port range.
// returns the input if port is valid
// exits with errror if port is out of range, or port <1024 is called without priviliges
func validatePort(port_tainted int) (int, error) {
    if(port_tainted < 0 || port_tainted > 65535) {
        glg.Errorf("port %d is invalid, please provide a port in range 0-65535",port_tainted)
        return 0, errors.New("Port is out of range, please provide a port in range 0-65535")
    }
    if(os.Geteuid() != 0 && port_tainted < 1024) {
        glg.Errorf("cannot bind to port %d without root priviliges",port_tainted)
        return 0, errors.New("Ports below 1024 can only be bind by root. You are not root.")
    }

    return port_tainted, nil
}
