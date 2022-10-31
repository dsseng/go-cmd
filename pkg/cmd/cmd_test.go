// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/siderolabs/go-cmd/pkg/cmd"
	"github.com/siderolabs/go-cmd/pkg/cmd/proc/reaper"
)

type CmdSuite struct {
	suite.Suite

	runReaper bool
}

func (suite *CmdSuite) SetupSuite() {
	if suite.runReaper {
		reaper.Run()
	}
}

func (suite *CmdSuite) TearDownSuite() {
	if suite.runReaper {
		reaper.Shutdown()
	}
}

func (suite *CmdSuite) TestRun() {
	type args struct {
		name string
		args []string
	}

	tests := []struct { //nolint:govet
		name      string
		args      args
		wantErr   bool
		errString string
	}{
		{
			"true",
			args{
				"true",
				[]string{},
			},
			false,
			"",
		},
		{
			"false",
			args{
				"false",
				[]string{},
			},
			true,
			"exit status 1: ",
		},
		{
			"false with output",
			args{
				"/bin/sh",
				[]string{
					"-c",
					"ls /not/found",
				},
			},
			true,
			"exit status 1: ls: /not/found: No such file or directory\n",
		},
		{
			"signal crash",
			args{
				"/bin/sh",
				[]string{
					"-c",
					"kill -2 $$",
				},
			},
			true,
			"signal: interrupt: ",
		},
		{
			"badexec",
			args{
				"badcommand",
				[]string{},
			},
			true,
			"exec: \"badcommand\": executable file not found in $PATH: ",
		},
	}

	for _, t := range tests {
		println(t.name)

		_, err := cmd.Run(t.args.name, t.args.args...)

		if t.wantErr {
			suite.Assert().Error(err)
			suite.Assert().Equal(t.errString, err.Error())
		} else {
			suite.Assert().NoError(err)
		}
	}

	stdout, err := cmd.RunContext(cmd.WithStdin(context.Background(), strings.NewReader("hello")), "xargs", "echo")
	suite.Assert().NoError(err)
	suite.Assert().Equal(stdout, "hello\n")
}

func TestCmdSuite(t *testing.T) {
	for _, runReaper := range []bool{true, false} {
		func(runReaper bool) {
			t.Run(fmt.Sprintf("runReaper=%v", runReaper), func(t *testing.T) { suite.Run(t, &CmdSuite{runReaper: runReaper}) })
		}(runReaper)
	}
}
