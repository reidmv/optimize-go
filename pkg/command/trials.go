/*
Copyright 2021 GramLabs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thestormforge/optimize-go/pkg/api"
	experiments "github.com/thestormforge/optimize-go/pkg/api/experiments/v1alpha1"
)

func newTrialsCommand(cfg Config) *cobra.Command {
	return &cobra.Command{
		Use:     "trials [NAME ...]",
		Aliases: []string{"trial"},

		// Trial names start with experiment names, so we can reuse the completion code
		ValidArgsFunction: validArgs(cfg, func(l *completionLister, toComplete string) (completions []string, directive cobra.ShellCompDirective) {
			directive |= cobra.ShellCompDirectiveNoFileComp
			l.forAllExperiments(func(item *experiments.ExperimentItem) {
				if strings.HasPrefix(item.Name.String(), toComplete) {
					completions = append(completions, item.Name.String())
				}
			})

			if len(completions) == 1 && completions[0] == toComplete {
				completions[0] += "-"
				directive |= cobra.ShellCompDirectiveNoSpace
			}

			return
		}),
	}
}

// NewGetTrialsCommand returns a command for getting trials.
func NewGetTrialsCommand(cfg Config, p Printer) *cobra.Command {
	var (
		selector string
		all      bool
	)

	cmd := newTrialsCommand(cfg)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx, out := cmd.Context(), cmd.OutOrStdout()
		client, err := api.NewClient(cfg.Address(), nil)
		if err != nil {
			return err
		}

		l := experiments.Lister{
			API: experiments.NewAPI(client),
		}

		result := &TrialOutput{Items: make([]TrialRow, 0, len(args))}

		q := experiments.TrialListQuery{}
		q.SetLabelSelector(parseLabelSelector(selector))
		q.SetStatus(experiments.TrialActive, experiments.TrialCompleted, experiments.TrialFailed)
		if all {
			q.AddStatus(experiments.TrialStaged)
		}

		if err := l.ForEachNamedTrial(ctx, args, q, false, result.Add); err != nil {
			return err
		}

		return p.Fprint(out, result)
	}

	cmd.Flags().StringVarP(&selector, "selector", "l", selector, "selector (label `query`) to filter on")
	cmd.Flags().BoolVarP(&all, "all", "A", all, "include all resources")

	return cmd
}

// NewDeleteTrialsCommand returns a command for deleting ("abandoning") trials.
func NewDeleteTrialsCommand(cfg Config, p Printer) *cobra.Command {
	var (
		ignoreNotFound bool
	)

	cmd := newTrialsCommand(cfg)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx, out := cmd.Context(), cmd.OutOrStdout()
		client, err := api.NewClient(cfg.Address(), nil)
		if err != nil {
			return err
		}

		l := experiments.Lister{
			API: experiments.NewAPI(client),
		}

		q := experiments.TrialListQuery{}
		q.SetStatus(experiments.TrialActive)
		return l.ForEachNamedTrial(ctx, args, q, ignoreNotFound, func(item *experiments.TrialItem) error {
			selfURL := item.Link(api.RelationSelf)
			if selfURL == "" {
				return fmt.Errorf("malformed response, missing self link")
			}

			err = l.API.AbandonRunningTrial(ctx, selfURL)
			if err != nil {
				return err
			}

			return p.Fprint(out, item)
		})
	}

	return cmd
}

// NewLabelTrialsCommand returns a command for labeling trials.
func NewLabelTrialsCommand(cfg Config, p Printer) *cobra.Command {
	cmd := newTrialsCommand(cfg)
	// TODO Should we extend validargsfn with suggestions like `baseline=true` ?
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx, out := cmd.Context(), cmd.OutOrStdout()
		client, err := api.NewClient(cfg.Address(), nil)
		if err != nil {
			return err
		}

		l := experiments.Lister{
			API: experiments.NewAPI(client),
		}

		q := experiments.TrialListQuery{}
		q.SetStatus(experiments.TrialCompleted)
		names, labels := argsToNamesAndLabels(args)
		return l.ForEachNamedTrial(ctx, names, q, false, func(item *experiments.TrialItem) error {
			labelsURL := item.Link(api.RelationLabels)
			if labelsURL == "" {
				return fmt.Errorf("malformed response, missing labels link")
			}

			err = l.API.LabelTrial(ctx, labelsURL, experiments.TrialLabels{Labels: labels})
			if err != nil {
				return err
			}

			return p.Fprint(out, item)
		})
	}

	return cmd
}
