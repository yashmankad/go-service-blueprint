// This class provides methods to initialize a unit testcase and tear it down

package util

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

// Test object used during unit tests
type Test struct {
	// directory for test logs
	TestDir string
}

// testDir where unit test logs will be written
var testDir = flag.String("test_dir", "", "Directory for test files and logs")

// TestInit initializes the unit test object
func TestInit(testName string) (*Test, error) {
	// test logs get written to the requested <testDir>/<testcase-name-randomId>
	// if the requested testDir is empty we place the logs under $HOMEDIR/testout/<testcase-name-randomId>
	// create this directory

	randomNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256)
	*testDir = filepath.Join(*testDir, fmt.Sprintf("%s-%d", testName, randomNum))
	if _, err := os.Stat(*testDir); os.IsNotExist(err) {
		if err := os.MkdirAll(*testDir, os.ModePerm); err != nil {
			log.Fatalf("failed to create test directory, error: %v", err)
			return nil, err
		}
	}

	return &Test{TestDir: *testDir}, nil
}

// TestCleanup cleans up any test artifacts and logs (if the test was successful)
func (t *Test) TestCleanup(test *testing.T) {
	if !test.Failed() {
		if err := os.RemoveAll(t.TestDir); err != nil {
			test.Errorf("test directory cleanup failed %s: %v", t.TestDir, err)
		}
	}
}
