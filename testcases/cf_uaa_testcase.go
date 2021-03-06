package testcases

import (
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CfUaaTestCase struct {
	testUser     string
	testPassword string
	name         string
}

func NewCfUaaTestCase() *CfUaaTestCase {
	randomString := RandomStringNumber()
	testUser := "uaa-test-user-" + randomString
	testPassword := "uaa-test-password-" + randomString
	return &CfUaaTestCase{testUser: testUser, testPassword: testPassword, name: "cf-uaa"}
}

func login(config Config, username, password string) {
	RunCommandSuccessfully("cf api --skip-ssl-validation", config.DeploymentToBackup.ApiUrl)
	RunCommandSuccessfully("cf auth", username, password)
}

func (tc *CfUaaTestCase) Name() string {
	return tc.name
}

func (tc *CfUaaTestCase) BeforeBackup(config Config) {
	By("we create a user and can login")
	login(config, config.DeploymentToBackup.AdminUsername, config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf create-user ", tc.testUser, tc.testPassword)
	login(config, tc.testUser, tc.testPassword)
	RunCommandSuccessfully("cf logout")
}

func (tc *CfUaaTestCase) AfterBackup(config Config) {
	By("we delete the user and verify")
	login(config, config.DeploymentToBackup.AdminUsername, config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf delete-user ", tc.testUser, "-f")
	RunCommandSuccessfully("cf logout")
	result := RunCommand("cf auth", tc.testUser, "password")
	Expect(result.ExitCode()).To(Equal(1))
}

func (tc *CfUaaTestCase) AfterRestore(config Config) {
	By("we can login again")
	login(config, tc.testUser, tc.testPassword)
}

func (tc *CfUaaTestCase) Cleanup(config Config) {
	By("We delete the user")
	login(config, config.DeploymentToBackup.AdminUsername, config.DeploymentToBackup.AdminPassword)
	RunCommandSuccessfully("cf delete-user ", tc.testUser, "-f")
}
