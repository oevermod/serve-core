package serve

import (
    "testing"
    "fmt"
)

func TestValidatePath(t *testing.T) {
    invalidPaths := []string{"://","/'.'/"}
    validPaths := []string{".","./","/",""}

    for _, path := range invalidPaths {
        _, err := validatePath(path)
        if err == nil {
            t.Errorf("Path '%s' is supposed to be invalid.", path)
        }
    }

    for _, path := range validPaths {
        _, err := validatePath(path)
        if err != nil {
            fmt.Println(err)
            t.Errorf("Path '%s' is supposed to be valid.", path)
        }
    }
}

func TestValidateAddress(t *testing.T) {
    invalidIPs := []string{"string",":::::1::::","1.2.3.4.5","300.200.100.000","localhost"}
    validIPs := []string{"::","0.0.0.0","","::1"}

    for _, ip := range invalidIPs {
        _, err := validateAddress(ip)
        if err == nil {
            t.Errorf("IP '%s' is supposed to be invalid.", ip)
        }
    }

    for _, ip := range validIPs {
        _, err := validateAddress(ip)
        if err != nil {
            t.Errorf("IP '%s' is supposed to be valid.", ip)
        }
    }
}

func TestValidatePort(t *testing.T) {
    // test valid ports
    for i := 1024; i < 65536; i++ {
        _, err := validatePort(i)
        if err != nil {
            t.Errorf("Port %d should be valid",i)
        }
    }

    // test admin ports
    for j := 1; j < 1024; j++ {
        _, err := validatePort(j)
        if err == nil {
            t.Errorf("Port %d should not be valid",j)
        }
    }

    // test negative ports
    for k := 0; k > -65536; k-- {
        _, err := validatePort(k)
        if err == nil {
            t.Errorf("Port %d should not be valid",k)
        }
    }

    // test ports above 65536
    for l := 65536; l < 131072; l++ {
        _, err := validatePort(l)
        if err == nil {
            t.Errorf("Port %d should not be valid",l)
        }
    }
}
