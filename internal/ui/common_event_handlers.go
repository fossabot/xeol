package ui

import (
	"fmt"
	"io"

	"github.com/gookit/color"
	"github.com/wagoodman/go-partybus"

	xeolEventParsers "github.com/xeol-io/xeol/xeol/event/parsers"
	policyTypes "github.com/xeol-io/xeol/xeol/policy/types"
)

func handleNotaryPolicyEvaluationMessage(event partybus.Event, reportOutput io.Writer) error {
	// show the report to stdout
	nt, err := xeolEventParsers.ParseNotaryPolicyEvaluationMessage(event)
	if err != nil {
		return fmt.Errorf("bad %s event: %w", event.Type, err)
	}

	var message string
	if nt.Action == policyTypes.PolicyActionDeny {
		message = color.Red.Sprintf("[%s][%s] Policy Violation: image '%s' is not signed by a trusted party.\n", nt.Action, nt.Type, nt.ImageReference)
	} else {
		if nt.FailDate != "" {
			message = color.Yellow.Sprintf("[%s][%s] Policy Violation: image '%s' is not signed by a trusted party. This policy will fail builds starting on %s.\n", nt.Action, nt.Type, nt.ImageReference, nt.FailDate)
		} else {
			message = color.Yellow.Sprintf("[%s][%s] Policy Violation: image '%s' is not signed by a trusted party.\n", nt.Action, nt.Type, nt.ImageReference)
		}
	}
	if _, err := reportOutput.Write([]byte(message)); err != nil {
		return fmt.Errorf("unable to show policy evaluation message: %w", err)
	}
	return nil
}

func handleEolPolicyEvaluationMessage(event partybus.Event, reportOutput io.Writer) error {
	// show the report to stdout
	pt, err := xeolEventParsers.ParseEolPolicyEvaluationMessage(event)
	if err != nil {
		return fmt.Errorf("bad %s event: %w", event.Type, err)
	}

	var message string
	if pt.Action == policyTypes.PolicyActionDeny {
		message = color.Red.Sprintf("[%s][%s] Policy Violation: %s (v%s) needs to be upgraded to a newer version.\n", pt.Action, pt.Type, pt.ProductName, pt.Cycle)
	} else {
		if pt.FailDate != "" {
			message = color.Yellow.Sprintf("[%s][%s] Policy Violation: %s (v%s) needs to be upgraded to a newer version. This policy will fail builds starting on %s.\n", pt.Action, pt.Type, pt.ProductName, pt.Cycle, pt.FailDate)
		} else {
			message = color.Yellow.Sprintf("[%s][%s] Policy Violation: %s (v%s) needs to be upgraded to a newer version.\n", pt.Action, pt.Type, pt.ProductName, pt.Cycle)
		}
	}
	if _, err := reportOutput.Write([]byte(message)); err != nil {
		return fmt.Errorf("unable to show policy evaluation message: %w", err)
	}
	return nil
}

func handleEolScanningFinished(event partybus.Event, reportOutput io.Writer) error {
	// show the report to stdout
	pres, err := xeolEventParsers.ParseEolScanningFinished(event)
	if err != nil {
		return fmt.Errorf("bad CatalogerFinished event: %w", err)
	}

	if err := pres.Present(reportOutput); err != nil {
		return fmt.Errorf("unable to show eol report: %w", err)
	}
	return nil
}

func handleNonRootCommandFinished(event partybus.Event, reportOutput io.Writer) error {
	// show the report to stdout
	result, err := xeolEventParsers.ParseNonRootCommandFinished(event)
	if err != nil {
		return fmt.Errorf("bad NonRootCommandFinished event: %w", err)
	}

	if _, err := reportOutput.Write([]byte(*result)); err != nil {
		return fmt.Errorf("unable to show eol report: %w", err)
	}
	return nil
}
