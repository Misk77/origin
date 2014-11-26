package cmd

import (
	"io"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kubecmd "github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl/cmd"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func (f *OriginFactory) NewCmdBuildLogs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build-logs <build>",
		Short: "Show container logs from the build container",
		Long: `Retrieve logs from the containers where the build occured

NOTE: This command may be moved in the future.

Examples:
$ kubectl build-logs 566bed879d2d
<stream logs from container to stdout>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				usageError(cmd, "<build> are required for log")
			}

			// TODO: Wouldn't be necessary, in upstream it is builtin.
			namespace := api.NamespaceDefault
			if ns := kubecmd.GetFlagString(cmd, "namespace"); len(ns) > 0 {
				namespace = ns
			}

			c, err := f.OriginClientFunc(cmd, nil)
			request := c.Get().Namespace(namespace).Path("redirect").Path("buildLogs").Path(args[0])
			readCloser, err := request.Stream()
			if err != nil {
				glog.Fatalf("Error: %v", err)
			}
			defer readCloser.Close()
			if _, err := io.Copy(out, readCloser); err != nil {
				glog.Fatalf("Error: %v", err)
			}
		},
	}
	return cmd
}
