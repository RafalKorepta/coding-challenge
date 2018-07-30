// Copyright [2018] [Rafał Korepta]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates zsg completion scripts",
	Long: `To load completion run

. <(portal-backend completion)

To configure your zsh shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(portal-backend completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RootCmd.GenZshCompletion(os.Stdout)
		if err != nil {
			zap.L().Error("unable to generate bash completion")
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
