package route

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/cloud"
	"github.com/evergreen-ci/evergreen/model/distro"
	"github.com/evergreen-ci/evergreen/model/event"
	"github.com/evergreen-ci/evergreen/rest/data"
	"github.com/evergreen-ci/evergreen/rest/model"
	"github.com/evergreen-ci/evergreen/units"
	"github.com/evergreen-ci/evergreen/util"
	"github.com/evergreen-ci/evergreen/validator"
	"github.com/evergreen-ci/gimlet"
	"github.com/mongodb/grip"
	"github.com/pkg/errors"
)

///////////////////////////////////////////////////////////////////////
//
// GET /rest/v2/distros/{distro_id}/setup

type distroIDGetSetupHandler struct {
	distroID string
	sc       data.Connector
}

func makeGetDistroSetup(sc data.Connector) gimlet.RouteHandler {
	return &distroIDGetSetupHandler{
		sc: sc,
	}
}

func (h *distroIDGetSetupHandler) Factory() gimlet.RouteHandler {
	return &distroIDGetSetupHandler{
		sc: h.sc,
	}
}

// Parse fetches the distroId from the http request.
func (h *distroIDGetSetupHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]

	return nil
}

// Run returns the given distro's setup script.
func (h *distroIDGetSetupHandler) Run(ctx context.Context) gimlet.Responder {
	d, err := h.sc.FindDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	apiDistro := &model.APIDistro{}
	if err = apiDistro.BuildFromService(d); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error converting from distro.Distro to model.APIDistro"))
	}

	return gimlet.NewJSONResponse(apiDistro.Setup)
}

///////////////////////////////////////////////////////////////////////
//
// PATCH /rest/v2/distros/{distro_id}/setup

type distroIDChangeSetupHandler struct {
	Setup    string
	distroID string
	sc       data.Connector
}

func makeChangeDistroSetup(sc data.Connector) gimlet.RouteHandler {
	return &distroIDChangeSetupHandler{
		sc: sc,
	}
}

func (h *distroIDChangeSetupHandler) Factory() gimlet.RouteHandler {
	return &distroIDChangeSetupHandler{
		sc: h.sc,
	}
}

// Parse fetches the distroId and JSON payload from the http request.
func (h *distroIDChangeSetupHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]
	body := util.NewRequestReader(r)
	defer body.Close()

	if err := util.ReadJSONInto(body, h); err != nil {
		return errors.Wrap(err, "Argument read error")
	}

	return nil
}

// Run updates the setup script for the given distroId.
func (h *distroIDChangeSetupHandler) Run(ctx context.Context) gimlet.Responder {
	d, err := h.sc.FindDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	d.Setup = h.Setup
	if err = h.sc.UpdateDistro(d, d); err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for update() by distro id '%s'", h.distroID))
	}

	apiDistro := &model.APIDistro{}
	if err = apiDistro.BuildFromService(d); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error converting from distro.Distro to model.APIDistro"))
	}

	return gimlet.NewJSONResponse(apiDistro)
}

///////////////////////////////////////////////////////////////////////
//
// GET /rest/v2/distros/{distro_id}/teardown

type distroIDGetTeardownHandler struct {
	distroID string
	sc       data.Connector
}

func makeGetDistroTeardown(sc data.Connector) gimlet.RouteHandler {
	return &distroIDGetTeardownHandler{
		sc: sc,
	}
}

func (h *distroIDGetTeardownHandler) Factory() gimlet.RouteHandler {
	return &distroIDGetTeardownHandler{
		sc: h.sc,
	}
}

// Parse fetches the distroId from the http request.
func (h *distroIDGetTeardownHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]

	return nil
}

// Run returns the given distro's teardown script.
func (h *distroIDGetTeardownHandler) Run(ctx context.Context) gimlet.Responder {
	d, err := h.sc.FindDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	apiDistro := &model.APIDistro{}
	if err = apiDistro.BuildFromService(d); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error converting from distro.Distro to model.APIDistro"))
	}

	return gimlet.NewJSONResponse(apiDistro.Teardown)
}

///////////////////////////////////////////////////////////////////////
//
// PATCH /rest/v2/distros/{distro_id}/teardown

type distroIDChangeTeardownHandler struct {
	Teardown string
	distroID string
	sc       data.Connector
}

func makeChangeDistroTeardown(sc data.Connector) gimlet.RouteHandler {
	return &distroIDChangeTeardownHandler{
		sc: sc,
	}
}

func (h *distroIDChangeTeardownHandler) Factory() gimlet.RouteHandler {
	return &distroIDChangeTeardownHandler{
		sc: h.sc,
	}
}

// Parse fetches the distroId and JSON payload from the http request.
func (h *distroIDChangeTeardownHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]
	body := util.NewRequestReader(r)
	defer body.Close()

	if err := util.ReadJSONInto(body, h); err != nil {
		return errors.Wrap(err, "Argument read error")
	}

	return nil
}

// Run updates the teardown script for the given distroId.
func (h *distroIDChangeTeardownHandler) Run(ctx context.Context) gimlet.Responder {
	d, err := h.sc.FindDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	d.Teardown = h.Teardown
	if err = h.sc.UpdateDistro(d, d); err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for update() by distro id '%s'", h.distroID))
	}

	apiDistro := &model.APIDistro{}
	if err = apiDistro.BuildFromService(d); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error converting from distro.Distro to model.APIDistro"))
	}

	return gimlet.NewJSONResponse(apiDistro)
}

///////////////////////////////////////////////////////////////////////
//
// PUT /rest/v2/distros/{distro_id}

type distroIDPutHandler struct {
	distroID string
	body     []byte
	sc       data.Connector
	settings *evergreen.Settings
}

func makePutDistro(sc data.Connector, settings *evergreen.Settings) gimlet.RouteHandler {
	return &distroIDPutHandler{
		sc:       sc,
		settings: settings,
	}
}

func (h *distroIDPutHandler) Factory() gimlet.RouteHandler {
	return &distroIDPutHandler{
		sc:       h.sc,
		settings: h.settings,
	}
}

// Parse fetches the distroId and JSON payload from the http request.
func (h *distroIDPutHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]

	body := util.NewRequestReader(r)
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "Argument read error")
	}
	h.body = b

	return nil
}

// Run either:
// (a) replaces an existing resource with the entity defined in the JSON payload, or
// (b) creates a new resource based on the Request-URI and JSON payload
func (h *distroIDPutHandler) Run(ctx context.Context) gimlet.Responder {
	user := MustHaveUser(ctx)

	original, err := h.sc.FindDistroById(h.distroID)
	if err != nil && err.(gimlet.ErrorResponse).StatusCode != http.StatusNotFound {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	apiDistro := &model.APIDistro{
		Name: model.ToStringPtr(h.distroID),
		FinderSettings: model.APIFinderSettings{
			Version: model.ToStringPtr(evergreen.FinderVersionLegacy),
		},
		PlannerSettings: model.APIPlannerSettings{
			Version: model.ToStringPtr(evergreen.PlannerVersionLegacy),
		},
		DispatcherSettings: model.APIDispatcherSettings{
			Version: model.ToStringPtr(evergreen.DispatcherVersionRevised),
		},
		HostAllocatorSettings: model.APIHostAllocatorSettings{
			Version: model.ToStringPtr(evergreen.HostAllocatorUtilization),
		},
		BootstrapSettings: model.APIBootstrapSettings{
			Method:        model.ToStringPtr(distro.BootstrapMethodLegacySSH),
			Communication: model.ToStringPtr(distro.CommunicationMethodLegacySSH),
		},
		CloneMethod: model.ToStringPtr(distro.CloneMethodLegacySSH),
	}
	if err = json.Unmarshal(h.body, apiDistro); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error while unmarshalling JSON"))
	}

	// Existing resource
	if original != nil {
		newDistro, respErr := validateDistro(ctx, apiDistro, h.distroID, h.settings, false)
		if respErr != nil {
			return respErr
		}

		if err = h.sc.UpdateDistro(original, newDistro); err != nil {
			return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for update() distro with distro id '%s'", h.distroID))
		}
		event.LogDistroModified(h.distroID, user.Username(), newDistro)
		return gimlet.NewJSONResponse(struct{}{})
	}
	// New resource
	newDistro, respErr := validateDistro(ctx, apiDistro, h.distroID, h.settings, true)
	if respErr != nil {
		return respErr
	}

	responder := gimlet.NewJSONResponse(struct{}{})
	if err = responder.SetStatus(http.StatusCreated); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrapf(err, "Cannot set HTTP status code to %d", http.StatusCreated))
	}
	if err = h.sc.CreateDistro(newDistro); err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for insert() distro with distro id '%s'", h.distroID))
	}

	return responder
}

///////////////////////////////////////////////////////////////////////
//
// DELETE /rest/v2/distros/{distro_id}

type distroIDDeleteHandler struct {
	distroID string
	sc       data.Connector
}

func makeDeleteDistroByID(sc data.Connector) gimlet.RouteHandler {
	return &distroIDDeleteHandler{
		sc: sc,
	}
}

func (h *distroIDDeleteHandler) Factory() gimlet.RouteHandler {
	return &distroIDDeleteHandler{
		sc: h.sc,
	}
}

// Parse fetches the distroId from the http request.
func (h *distroIDDeleteHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]

	return nil
}

// Run deletes a distro by id.
func (h *distroIDDeleteHandler) Run(ctx context.Context) gimlet.Responder {
	_, err := h.sc.FindDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	err = h.sc.DeleteDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for remove() by distro id '%s'", h.distroID))
	}

	return gimlet.NewJSONResponse(struct{}{})
}

////////////////////////////////////////////////////////////////////////
//
// PATCH /rest/v2/distros/{distro_id}

type distroIDPatchHandler struct {
	distroID string
	body     []byte
	sc       data.Connector
	settings *evergreen.Settings
}

func makePatchDistroByID(sc data.Connector, settings *evergreen.Settings) gimlet.RouteHandler {
	return &distroIDPatchHandler{
		sc:       sc,
		settings: settings,
	}
}

func (h *distroIDPatchHandler) Factory() gimlet.RouteHandler {
	return &distroIDPatchHandler{
		sc:       h.sc,
		settings: h.settings,
	}
}

// Parse fetches the distroId from the http request.
func (h *distroIDPatchHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]

	body := util.NewRequestReader(r)
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "Argument read error")
	}
	h.body = b

	return nil
}

// Run updates a distro by id.
func (h *distroIDPatchHandler) Run(ctx context.Context) gimlet.Responder {
	old, err := h.sc.FindDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	apiDistro := &model.APIDistro{}
	if err = apiDistro.BuildFromService(old); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error converting from distro.Distro to model.APIDistro"))
	}
	oldSettingsList := apiDistro.ProviderSettingsList
	apiDistro.ProviderSettingsList = nil
	if err = json.Unmarshal(h.body, apiDistro); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error while unmarshalling JSON"))
	}
	if len(apiDistro.ProviderSettingsList) == 0 {
		apiDistro.ProviderSettingsList = oldSettingsList
	}

	d, respErr := validateDistro(ctx, apiDistro, h.distroID, h.settings, false)
	if respErr != nil {
		return respErr
	}

	if err = h.sc.UpdateDistro(old, d); err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for update() by distro id '%s'", h.distroID))
	}

	return gimlet.NewJSONResponse(apiDistro)
}

////////////////////////////////////////////////////////////////////////
//
// GET /rest/v2/distros/{distro_id}

type distroIDGetHandler struct {
	distroID string
	sc       data.Connector
}

func makeGetDistroByID(sc data.Connector) gimlet.RouteHandler {
	return &distroIDGetHandler{
		sc: sc,
	}
}

func (h *distroIDGetHandler) Factory() gimlet.RouteHandler {
	return &distroIDGetHandler{
		sc: h.sc,
	}
}

// Parse fetches the distroId from the http request.
func (h *distroIDGetHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distroID = gimlet.GetVars(r)["distro_id"]

	return nil
}

// Run calls the data FindDistroById function and returns the distro from the provider.
func (h *distroIDGetHandler) Run(ctx context.Context) gimlet.Responder {
	d, err := h.sc.FindDistroById(h.distroID)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "Database error for find() by distro id '%s'", h.distroID))
	}

	apiDistro := &model.APIDistro{}
	if err = apiDistro.BuildFromService(d); err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error converting from distro.Distro to model.APIDistro"))
	}

	return gimlet.NewJSONResponse(apiDistro)
}

////////////////////////////////////////////////////////////////////////
//
// GET /rest/v2/distros

type distroGetHandler struct {
	sc data.Connector
}

func makeDistroRoute(sc data.Connector) gimlet.RouteHandler {
	return &distroGetHandler{
		sc: sc,
	}
}

func (h *distroGetHandler) Factory() gimlet.RouteHandler {
	return &distroGetHandler{
		sc: h.sc,
	}
}

func (h *distroGetHandler) Parse(ctx context.Context, r *http.Request) error {
	return nil
}

func (h *distroGetHandler) Run(ctx context.Context) gimlet.Responder {
	distros, err := h.sc.FindAllDistros()
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrap(err, "Database error for find() all distros"))
	}

	resp := gimlet.NewResponseBuilder()
	if err = resp.SetFormat(gimlet.JSON); err != nil {
		return gimlet.MakeJSONErrorResponder(err)
	}

	for _, d := range distros {
		distroModel := &model.APIDistro{}
		if err = distroModel.BuildFromService(d); err != nil {
			return gimlet.MakeJSONErrorResponder(err)
		}

		err = resp.AddData(distroModel)
		if err != nil {
			return gimlet.MakeJSONErrorResponder(err)
		}
	}

	return resp
}

////////////////////////////////////////////////////////////////////////

func validateDistro(ctx context.Context, apiDistro *model.APIDistro, resourceID string, settings *evergreen.Settings, isNewDistro bool) (*distro.Distro, gimlet.Responder) {
	if apiDistro.ProviderSettings != nil && len(apiDistro.ProviderSettings) > 0 {
		return nil, gimlet.MakeJSONErrorResponder(errors.New("must use provider_settings list to update settings"))
	}

	i, err := apiDistro.ToService()
	if err != nil {
		return nil, gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "API error converting from model.APIDistro to distro.Distro"))
	}
	d, ok := i.(*distro.Distro)
	if !ok {
		return nil, gimlet.MakeJSONErrorResponder(gimlet.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Unexpected type %T for distro.Distro", i),
		})
	}

	if err = cloud.UpdateProviderSettings(d); err != nil {
		return nil, gimlet.MakeJSONErrorResponder(gimlet.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	id := model.FromStringPtr(apiDistro.Name)
	if resourceID != id {
		return nil, gimlet.MakeJSONErrorResponder(gimlet.ErrorResponse{
			StatusCode: http.StatusForbidden,
			Message:    fmt.Sprintf("A distro's name is immutable; cannot rename distro '%s'", resourceID),
		})
	}

	vErrors, err := validator.CheckDistro(ctx, d, settings, isNewDistro)
	if err != nil {
		return nil, gimlet.MakeJSONErrorResponder(gimlet.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}
	if len(vErrors) != 0 {
		return nil, gimlet.MakeJSONErrorResponder(gimlet.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    vErrors.String(),
		})
	}

	return d, nil
}

///////////////////////////////////////////////////////////////////////
//
// PATCH /rest/v2/distros/{distro}/execute

type distroExecuteHandler struct {
	opts   model.APIDistroScriptOptions
	distro string
	sc     data.Connector
	env    evergreen.Environment
}

func makeDistroExecute(sc data.Connector, env evergreen.Environment) gimlet.RouteHandler {
	return &distroExecuteHandler{
		sc:  sc,
		env: env,
	}
}

func (h *distroExecuteHandler) Factory() gimlet.RouteHandler {
	return &distroExecuteHandler{
		sc:  h.sc,
		env: h.env,
	}
}

// Parse fetches the distro and JSON payload from the http request.
func (h *distroExecuteHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distro = gimlet.GetVars(r)["distro_id"]
	body := util.NewRequestReader(r)
	defer body.Close()

	if err := util.ReadJSONInto(body, &h.opts); err != nil {
		return errors.Wrap(err, "could not read request")
	}

	if h.opts.Script == "" {
		return errors.New("cannot execute an empty script")
	}
	if !h.opts.IncludeTaskHosts && !h.opts.IncludeSpawnHosts {
		return errors.New("cannot exclude both spawn hosts and task hosts from script execution")
	}

	return nil
}

// Run enqueues a job to run a script on all selected hosts that are not down
// for the given given distro ID.
func (h *distroExecuteHandler) Run(ctx context.Context) gimlet.Responder {
	hosts, err := h.sc.FindHostsByDistro(h.distro)
	if err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrapf(err, "could not find hosts for the distro %s", h.distro))
	}

	var allHostIDs []string
	for _, host := range hosts {
		allHostIDs = append(allHostIDs, host.Id)
	}
	catcher := grip.NewBasicCatcher()
	var hostIDs []string
	for _, host := range hosts {
		ts := util.RoundPartOfMinute(0).Format(units.TSFormat)
		if (host.StartedBy == evergreen.User && h.opts.IncludeTaskHosts) || (host.StartedBy != evergreen.User && h.opts.IncludeSpawnHosts) {
			if err = h.env.RemoteQueue().Put(ctx, units.NewHostExecuteJob(h.env, host, h.opts.Script, h.opts.Sudo, h.opts.SudoUser, ts)); err != nil {
				catcher.Wrapf(err, "problem enqueueing job to run script on host '%s'", host.Id)
				continue
			}
			hostIDs = append(hostIDs, host.Id)
		}
	}
	if catcher.HasErrors() {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrapf(err, "problem enqueueing jobs to run script on hosts"))
	}

	return gimlet.NewJSONResponse(struct {
		HostIDs []string `json:"host_ids"`
	}{HostIDs: hostIDs})
}

///////////////////////////////////////////////////////////////////////
//
// PATCH /rest/v2/distros/{distro}/icecream_config

type distroIcecreamConfigHandler struct {
	distro string
	opts   model.APIDistroScriptOptions
	sc     data.Connector
	env    evergreen.Environment
}

func makeDistroIcecreamConfig(sc data.Connector, env evergreen.Environment) gimlet.RouteHandler {
	return &distroIcecreamConfigHandler{
		sc:  sc,
		env: env,
	}
}

func (h *distroIcecreamConfigHandler) Factory() gimlet.RouteHandler {
	return &distroIcecreamConfigHandler{
		sc:  h.sc,
		env: h.env,
	}
}

// Parse extracts the distro and JSON payload from the http request.
func (h *distroIcecreamConfigHandler) Parse(ctx context.Context, r *http.Request) error {
	h.distro = gimlet.GetVars(r)["distro_id"]
	body := util.NewRequestReader(r)
	defer body.Close()

	if err := util.ReadJSONInto(body, &h.opts); err != nil {
		return errors.Wrap(err, "could not read request body")
	}

	return nil
}

// Run enqueues a job to run a script on all hosts that are not down for the
// given given distro ID.
func (h *distroIcecreamConfigHandler) Run(ctx context.Context) gimlet.Responder {
	hosts, err := h.sc.FindHostsByDistro(h.distro)
	if err != nil {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "could not find hosts for the distro '%s'", h.distro))
	}

	dat, err := distro.NewDistroAliasesLookupTable()
	if err != nil {
		return gimlet.MakeJSONInternalErrorResponder(errors.Wrap(err, "could not get distro lookup table"))
	}

	catcher := grip.NewBasicCatcher()
	var hostIDs []string
	for _, host := range hosts {
		if host.StartedBy == evergreen.User || !host.IsVirtualWorkstation {
			continue
		}

		// If the distro exists, we use the settings directly from that distro;
		// if the distro in the host document is deleted, we make a best-effort
		// attempt to resolve it to a real distro by attempting to pick any
		// existing distro with an alias that matches the deleted distro.
		distroIDs := dat.Expand([]string{host.Distro.Id})
		if len(distroIDs) == 0 {
			catcher.Errorf("could not look up distro '%s'", host.Distro.Id)
			continue
		}
		var distros []distro.Distro
		distros, err = distro.Find(distro.ByIds(distroIDs))
		if err != nil {
			catcher.Errorf("could not find distros matching '%s' for host '%s'", host.Distro.Id, host.Id)
			continue
		}
		var d distro.Distro
		var distroFound bool
		for _, d = range distros {
			if d.IcecreamSettings.Populated() {
				distroFound = true
				break
			}
		}
		if !distroFound {
			catcher.Wrapf(err, "could not resolve distro '%s' for host '%s'", host.Distro.Id, host.Id)
			continue
		}

		script := d.IcecreamSettings.GetUpdateConfigScript()
		ts := util.RoundPartOfMinute(0).Format(units.TSFormat)
		if err = h.env.RemoteQueue().Put(ctx, units.NewHostExecuteJob(h.env, host, script, true, "root", ts)); err != nil {
			catcher.Wrapf(err, "problem enqueueing job to update icecream config file on host '%s'", host.Id)
			continue
		}
		hostIDs = append(hostIDs, host.Id)
	}

	if catcher.HasErrors() {
		return gimlet.MakeJSONErrorResponder(errors.Wrapf(err, "problem enqueueing jobs to update icecream config on hosts"))
	}

	return gimlet.NewJSONResponse(struct {
		HostIDs []string `json:"host_ids"`
	}{HostIDs: hostIDs})
}
