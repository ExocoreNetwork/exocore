package task_test

import (
	"github.com/ExocoreNetwork/exocore/testutil"
	tasktype "github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"
	"testing"

	"github.com/ExocoreNetwork/exocore/precompiles/avsTask"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

var s *TaskPrecompileTestSuite

type TaskPrecompileTestSuite struct {
	testutil.BaseTestSuite
	precompile *task.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(TaskPrecompileTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Precompile Suite")
}

func (s *TaskPrecompileTestSuite) SetupTest() {
	s.DoSetupTest()
	precompile, err := task.NewPrecompile(s.App.AuthzKeeper, s.App.TaskKeeper, tasktype.AvsKeeper{})
	s.Require().NoError(err)
	s.precompile = precompile
}
