package taskwarrior

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUDAs_noUDAs_existAfterwards(t *testing.T) {
	// Use a custom .taskrc for testing.
	taskrc := filepath.Join(t.TempDir(), "taskrc")
	file, err := os.Create(taskrc)
	assert.NoError(t, err, "Failed to create .taskrc for testing.")
	defer file.Close()
	os.Setenv("TASKRC", taskrc)
	udaName := "foo"
	udaLabel := "bar"

	err = CreateUDA(udaName, udaLabel)

	assert.NoError(t, err, "CreateUDA returned an error.")
	out, err := exec.Command("bash", "-c", fmt.Sprintf("cat %s", taskrc)).
		CombinedOutput()
	fmt.Println(string(out))
	err = exec.Command("bash", "-c", fmt.Sprintf("cat %s | grep uda.%s.label=%s", taskrc, udaName, udaLabel)).
		Run()
	assert.NoError(
		t,
		err,
		fmt.Sprintf(
			"'%s' has no config entry 'uda.%s.label=%s'.",
			taskrc,
			udaName,
			udaLabel,
		),
	)
	err = exec.Command("bash", "-c", fmt.Sprintf("cat %s | grep uda.%s.type=string", taskrc, udaName)).
		Run()
	assert.NoError(
		t,
		err,
		fmt.Sprintf("'%s' has no config entry 'uda.%s.type=string'.", taskrc, udaName),
	)
}
