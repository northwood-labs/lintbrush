// Copyright 2024, Northwood Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hashicorp/go-multierror"
	clihelpers "github.com/northwood-labs/cli-helpers"
	"github.com/northwood-labs/lintbrush/lintbrush"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
)

var (
	fOutput  string
	fSarif   bool
	fVerbose bool

	errs *multierror.Error
	now  = time.Now()

	// Standard logger
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "lintbrush",
		Short: `A linter which focuses on low-level file hygeine.`,
		Long: clihelpers.LongHelpText(`
		lintbrush

		A linter which focuses on the low-level things that are common across many
		repositories. It favors proactively fixing things first, then falling back to
		providing an error when something cannot be fixed automatically.

		Leaves the language-specific linting to more appropriate tools. Instead, this
		focuses on general file hygeine, security, and other things that are common
		across many repositories.

		--------------------------------------------------------------------------------

		👀 Check that executable files have shebangs`),
		Run: func(cmd *cobra.Command, args []string) {
			if fVerbose {
				logger.SetReportCaller(true)
			}

			if len(args) == 0 {
				logger.Fatal("No files to lint")
			}

			// Create a new Sarif report object.
			report, err := sarif.New(sarif.Version210)
			if err != nil {
				logger.Fatal(err.Error())
			}

			run := sarif.NewRunWithInformationURI("lintbrush", "https://github.com/northwood-labs/lintbrush")
			run.Tool.Driver.WithSemanticVersion("dev") // @version

			run.WithAutomationDetails(
				sarif.NewRunAutomationDetails().
					WithDescriptionText("This scan was run with Lintbrush vX.X.X on NOW.").
					WithID("Lintbrush run for repository/branch/DATE"). // @version
					WithCorrelationGUID(xid.New().String()),
			)

			//   "originalUriBaseIds": {
			//     "REPOROOT": {
			//       "description": {
			//         "text": "The directory into which the repo was cloned."
			//       },
			//       "properties": {
			//         "comment": "The SARIF producer has chosen not to specify a URI for REPOROOT. See §3.14.14, NOTE 1, for an explanation."
			//       }
			//     }
			//   },
			//   "artifacts": [
			//     {
			//       "location": {
			//         "uri": "sarif-tutorials/samples/Introduction/simple-example.js"
			//       },
			//       "length": 3444,
			//       "sourceLanguage": "javascript",
			//       "hashes": {
			//         "sha-256": "b13ce2678a8807ba0765ab94a0ecd394f869bc81"
			//       }
			//     }
			//   ],

			err = lintbrush.CheckExecutablesHaveShebangs(run, args)
			errs = multierror.Append(errs, err)

			// Print all errors.
			for i := range errs.Errors {
				logger.Errorf("%v", errs.Errors[i])
			}

			// 	run.AddRule(r.RuleID).
			// 		WithName("").
			// 		WithDescription(r.Description).
			// 		// WithFullDescription().
			// 		WithHelpURI(r.Link).
			// 		WithTextHelp("")

			// 	// add the location as a unique artifact
			// 	run.AddDistinctArtifact(r.Location.Filename)

			// 	// add each of the results with the details of where the issue occurred
			// 	run.CreateResultForRule(r.RuleID).
			// 		WithLevel(strings.ToLower(r.Severity)).
			// 		WithMessage(sarif.NewTextMessage(r.Description)).
			// 		AddLocation(
			// 			sarif.NewLocationWithPhysicalLocation(
			// 				sarif.NewPhysicalLocation().
			// 					WithArtifactLocation(
			// 						sarif.NewSimpleArtifactLocation(r.Location.Filename),
			// 					).WithRegion(
			// 					sarif.NewSimpleRegion(r.Location.StartLine, r.Location.EndLine),
			// 				),
			// 			),
			// 		)
			// }

			// add the run to the report
			report.AddRun(run)

			// print the report to stdout
			_ = report.PrettyWrite(os.Stdout)
			fmt.Fprint(os.Stdout, "\n")
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&fOutput, "output", "o", "", "Write results to a file instead of stdout")
	rootCmd.PersistentFlags().BoolVarP(&fSarif, "sarif", "s", false, "Return results in SARIF format")
	rootCmd.PersistentFlags().BoolVarP(&fVerbose, "verbose", "v", false, "Print verbose output")
}

func returnFirstEnv(envs ...string) (string, bool) {
	for i := range envs {
		if os.Getenv(envs[i]) != "" {
			return os.Getenv(envs[i]), true
		}
	}

	return "", false
}

// versionControlDetails := sarif.NewVersionControlDetails()
// versionControlDetails.WithAsOfTimeUTC(&now)

// if branchName, ok := returnFirstEnv("GITHUB_REF_NAME", "GITHUB_HEAD_REF", "GITHUB_SHA"); ok {
// 	versionControlDetails.WithBranch(branchName)
// }

// if repoName, ok := returnFirstEnv("GITHUB_REPOSITORY"); ok {
// 	versionControlDetails.WithRepositoryURI(os.Getenv("GITHUB_SERVER_URL") + "/" + repoName)
// }

// run.AddVersionControlProvenance(versionControlDetails)
