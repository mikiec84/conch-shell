// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package validation contains commands for validation related commands
package validation

import (
	"bufio"
	"os"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
)

type validationPlans []conch.ValidationPlan

func (vps validationPlans) renderTable() {
	table := util.GetMarkdownTable()
	table.SetHeader([]string{"Id", "Name", "Description"})

	for _, vp := range vps {
		table.Append([]string{vp.ID.String(), vp.Name, vp.Description})
	}

	table.Render()
}

func getValidationPlans(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	app.Action = func() {
		var validationPlans validationPlans
		validationPlans, err := util.API.GetValidationPlans()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(validationPlans)
			return
		}
		validationPlans.renderTable()
	}
}

func getValidationPlan(app *cli.Cmd) {
	app.Action = func() {
		validationPlan, err := util.API.GetValidationPlan(validationPlanUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(validationPlan)
			return
		}
		validationPlans := validationPlans{validationPlan}
		validationPlans.renderTable()
	}
}

func showValidationPlanValidations(app *cli.Cmd) {

	app.Action = func() {
		var validations validations
		validations, err := util.API.GetValidationPlanValidations(validationPlanUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(validations)
			return
		}
		validations.renderTable()
	}
}

func testValidationPlan(app *cli.Cmd) {
	var deviceSerial = app.StringArg("DEVICE_ID", "", "The Device ID (serial number) to test the validation plan against")

	app.Spec = "DEVICE_ID"

	app.Action = func() {
		body := bufio.NewReader(os.Stdin)
		var validationResults validationResults
		validationResults, err := util.API.RunDeviceValidationPlan(*deviceSerial, validationPlanUUID, body)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(validationResults)
			return
		}
		validationResults.renderTable()
	}
}
