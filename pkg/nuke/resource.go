package nuke

import (
	"github.com/cgroschupp/stackit-nuke/pkg/stackit"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/sirupsen/logrus"
)

const Account registry.Scope = "project"

type ListerOpts struct {
	Client    *stackit.Client
	ProjectID string
	Region    string
	Logger    *logrus.Entry
}
