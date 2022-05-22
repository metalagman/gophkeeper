package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gophkeeper/pkg/logger"
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret management",
	Long:  `Choose one of the command to do with secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.CheckErr(cmd.Help())
	},
}

var secretListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List own secrets",
	Long:  `Allows you to list own secrets stored on server`,
	Run:   secretList,
}

func init() {
	rootCmd.AddCommand(secretCmd)
	secretCmd.AddCommand(secretListCmd)
}

func secretList(cmd *cobra.Command, args []string) {

}

func clientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	metadata.AppendToOutgoingContext(ctx)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
