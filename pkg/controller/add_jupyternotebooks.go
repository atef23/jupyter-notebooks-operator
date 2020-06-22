package controller

import (
	"github.com/atef23/jupyter-notebooks-operator/pkg/controller/jupyternotebooks"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, jupyternotebooks.Add)
}
