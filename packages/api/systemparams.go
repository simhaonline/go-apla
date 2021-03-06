// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package api

import (
	"net/http"

	"github.com/AplaProject/go-apla/packages/consts"
	"github.com/AplaProject/go-apla/packages/model"

	log "github.com/sirupsen/logrus"
)

func getSystemParamsHandler(w http.ResponseWriter, r *http.Request) {
	form := &paramsForm{}
	if err := parseForm(r, form); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	logger := getLogger(r)

	list, err := model.GetAllSystemParameters(nil)
	if err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("Getting all system parameters")
	}

	result := &ecosystemParamsResult{
		List: make([]paramResult, 0),
	}

	acceptNames := form.AcceptNames()
	for _, item := range list {
		if len(acceptNames) > 0 && !acceptNames[item.Name] {
			continue
		}
		result.List = append(result.List, paramResult{
			Name:       item.Name,
			Value:      item.Value,
			Conditions: item.Conditions,
		})
	}

	if len(result.List) == 0 {
		errorResponse(w, errParamNotFound.Errorf(form.Names), http.StatusBadRequest)
		return
	}

	jsonResponse(w, result)
}
