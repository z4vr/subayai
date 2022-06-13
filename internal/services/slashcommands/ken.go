package slashcommands

import (
	"fmt"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/discord"
	"github.com/z4vr/subayai/internal/services/slashcommands/cmds"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/examples/usermsgcommands/commands"
)

func New(ctn di.Container) (k *ken.Ken, err error) {

	dc := ctn.Get("discord").(*discord.Discord)

	k, err = ken.New(dc.Session(), ken.Options{
		DependencyProvider: ctn,
		OnSystemError:      SystemErrorHandler,
		OnCommandError:     commandErrorHandler,
		EmbedColors: ken.EmbedColors{
			Default: 0xd32f2f,
			Error:   0x249ff2,
		},
	})
	if err != nil {
		return
	}

	err = k.RegisterCommands(
		new(cmds.Profile),
		new(commands.InfoUserCommand))
	if err != nil {
		return
	}

	return
}

func SystemErrorHandler(context string, err error, _ ...interface{}) {
	logrus.WithField("ctx", context).WithError(err).Error("slashcommand error")
}

func commandErrorHandler(err error, ctx *ken.Ctx) {
	// Is ignored if interaction has already been responded
	err = ctx.Defer()
	if err != nil {
		return
	}

	if err == ken.ErrNotDMCapable {
		err := ctx.RespondError("This command can not be used in DMs.", "An error occurred.")
		if err != nil {
			return
		}
		return
	}

	err = ctx.RespondError(
		fmt.Sprintf("The command execution failed unexpectedly:\n```\n%s\n```", err.Error()),
		"Command execution failed")
	if err != nil {
		return
	}
}
