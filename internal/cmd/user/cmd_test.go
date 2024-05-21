package user

import (
	"bytes"
	"testing"

	"github.com/cosmo-workspace/cosmo/pkg/cli"
	. "github.com/cosmo-workspace/cosmo/pkg/snap"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/cobra"
)

func TestCommandUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cosmoctl user suite")
}

var _ = Describe("help", func() {
	It("should match snapshot", func() {
		cmd := &cobra.Command{}
		out := bytes.Buffer{}
		cmd.SetOut(&out)
		AddCommand(cmd, cli.NewRootOptions())
		cmd.SetArgs([]string{"user", "--help"})
		err := cmd.Execute()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out.String()).To(MatchSnapShot())
	})
})
