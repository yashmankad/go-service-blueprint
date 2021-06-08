// This class provides methods to initialize a unit testcase and tear it down

package common

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

// test object definition
type Test struct {
	// directory for test logs
	TestDir string
}

func TestInit(testName string) (*Test, error) {
	// test logs get written to $HOME/testout/<testcase-name-randomId>
	// create this directory

	// fetch the users home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to fetch user home dir, error: %v", err)
		return nil, err
	}

	randomNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256)
	testDir := filepath.Join(homeDir, "testout", fmt.Sprintf("%s-%d", testName, randomNum))
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
			log.Fatalf("failed to create test directory, error: %v", err)
			return nil, err
		}
	}

	return &Test{TestDir: testDir}, nil
}

// cleans up the test env, including logs if the test was successful
func (t *Test) Cleanup(test *testing.T) {
	if !test.Failed() {
		if err := os.RemoveAll(t.TestDir); err != nil {
			test.Errorf("test directory cleanup failed %s: %v", t.TestDir, err)
		}
	}
}
