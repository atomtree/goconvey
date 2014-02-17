package system

import (
	"os"
	"os/exec"
	"runtime"
)

type Shell struct {
	coverage string
	gobin    string
}

func (self *Shell) GoTest(directory string) (output string, err error) {
	output, err = self.compileDependencies(directory)
	if err == nil {
		output, err = self.goTest(directory)
	}
	return
}

func (self *Shell) compileDependencies(directory string) (output string, err error) {
	return self.execute(directory, self.gobin, "test", "-i")
}

func (self *Shell) goTest(directory string) (output string, err error) {
	output, err = self.execute(directory, self.gobin, "test", "-v", "-timeout=-42s", self.coverage)

	if err != nil && self.coverage != "" {
		output, err = self.execute(directory, self.gobin, "test", "-v", "-timeout=-42s")
	}
	return
}

func (self *Shell) execute(directory, name string, args ...string) (output string, err error) {
	command := exec.Command(name, args...)
	command.Dir = directory
	rawOutput, err := command.CombinedOutput()
	output = string(rawOutput)
	return
}

func (self *Shell) Getenv(key string) string {
	return os.Getenv(key)
}

func (self *Shell) Setenv(key, value string) error {
	if self.Getenv(key) != value {
		return os.Setenv(key, value)
	}
	return nil
}

func NewShell(gobin string, cover bool) *Shell {
	self := new(Shell)
	self.gobin = gobin
	if cover && goVersion_1_2_orGreater() {
		self.coverage = coverageFlag
	}
	return self
}

func goVersion_1_2_orGreater() bool {
	version := runtime.Version() // 'go1.2....'
	major, minor := version[2], version[4]
	return major >= byte('1') && minor >= byte('2')
}

const coverageFlag = "-cover"
