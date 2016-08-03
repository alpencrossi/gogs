// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	api "github.com/gogits/go-gogs-client"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/context"
	"github.com/gogits/gogs/modules/log"
	"github.com/gogits/gogs/routers/api/v1/convert"
)

func ListLabels(ctx *context.APIContext) {
	labels, err := models.GetLabelsByRepoID(ctx.Repo.Repository.ID)
	if err != nil {
		ctx.Error(500, "Labels", err)
		return
	}

	apiLabels := make([]*api.Label, len(labels))
	for i := range labels {
		apiLabels[i] = convert.ToLabel(labels[i])
	}

	ctx.JSON(200, &apiLabels)
}

func GetLabel(ctx *context.APIContext) {
	label, err := models.GetLabelByID(ctx.ParamsInt64(":id"))
	if err != nil {
		if models.IsErrLabelNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetLabelByID", err)
		}
		return
	}

	ctx.JSON(200, convert.ToLabel(label))
}

func CreateLabel(ctx *context.APIContext, form api.LabelOption) {
	if !ctx.Repo.IsWriter() {
		ctx.Status(403)
		return
	}

	label := &models.Label{
		Name:   form.Name,
		Color:  form.Color,
		RepoID: ctx.Repo.Repository.ID,
	}
	err := models.NewLabel(label)
	if err != nil {
		ctx.Error(500, "NewLabel", err)
		return
	}

	label, err = models.GetLabelByID(label.ID)
	if err != nil {
		ctx.Error(500, "GetLabelByID", err)
		return
	}
	ctx.JSON(201, convert.ToLabel(label))
}

func EditLabel(ctx *context.APIContext, form api.LabelOption) {
	if !ctx.Repo.IsWriter() {
		ctx.Status(403)
		return
	}

	label, err := models.GetLabelByID(ctx.ParamsInt64(":id"))
	if err != nil {
		if models.IsErrLabelNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetLabelByID", err)
		}
		return
	}

	if len(form.Name) > 0 {
		label.Name = form.Name
	}
	if len(form.Color) > 0 {
		label.Color = form.Color
	}

	if err := models.UpdateLabel(label); err != nil {
		ctx.Handle(500, "UpdateLabel", err)
		return
	}
	ctx.JSON(200, convert.ToLabel(label))
}

func DeleteLabel(ctx *context.APIContext) {
	if !ctx.Repo.IsWriter() {
		ctx.Status(403)
		return
	}

	label, err := models.GetLabelByID(ctx.ParamsInt64(":id"))
	if err != nil {
		if models.IsErrLabelNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetLabelByID", err)
		}
		return
	}

	if err := models.DeleteLabel(ctx.Repo.Repository.ID, ctx.ParamsInt64(":id")); err != nil {
		ctx.Error(500, "DeleteLabel", err)
		return
	}

	log.Trace("Label deleted: %s %s", label.ID, label.Name)
	ctx.Status(204)
}