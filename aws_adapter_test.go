package rollerskates

import (
	"github.com/joho/godotenv"
	_ "github.com/orchestrate-io/dvr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallToApi(t *testing.T) {
	godotenv.Load()

	ids := GetInstanceIds("testing")
	assert.Equal(t, ids[0], "i-e0708cd1")
}

func TestReturnsTrueWhenConnectionsDrainedAndInstanceHasRestarted(t *testing.T) {
  godotenv.Load()
  assert.Equal(t, SkatesRestartLoadBalancerInstance("testing", "i-e0708cd1"), true)
}

func TestDeregisterInstancesFromLoadBalancer(t *testing.T) {
  godotenv.Load()
  assert.Equal(t, DeregisterInstancesFromLoadBalancer("testing", "i-e0708cd1"), true)
}

func TestWaitForInstancesToBeDeRegisteredFromLoadBalancer(t *testing.T) {
  godotenv.Load()
  DeregisterInstancesFromLoadBalancer("testing", "i-e0708cd1")
	assert.Equal(t, WaitForInstancesToBeDeregisteredFromLoadBalancer("testing", "i-e0708cd1"), true)
}

func TestRebootInstance(t *testing.T) {
  godotenv.Load()
	assert.Equal(t, RebootInstance("i-e0708cd1"), true)
}

func TestWaitForInstaceToBeRebooted(t *testing.T) {
  godotenv.Load()
  RebootInstance("i-75db0574")
  assert.Equal(t, WaitForInstaceToBeRebooted("i-75db0574"), true)
}
