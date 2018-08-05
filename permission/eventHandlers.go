package permission

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dukfaar/goUtils/eventbus"
)

type LoginSuccessMsg struct {
	UserID               string `json:"userId,omitempty"`
	AccessToken          string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt string `json:"accessTokenExpiresAt,omitempty"`
}

type RoleMsg struct {
	ID          string   `json:"_id,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type PermissionMsg struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserMsg struct {
	ID    string   `json:"_id,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

func permissionHandler(permissionService *Service, msg []byte) error {
	permissionMsg := PermissionMsg{}

	err := json.Unmarshal(msg, &permissionMsg)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling: %v: %v", string(msg), err)
		fmt.Println(err)
		return err
	}

	permissionService.SetPermission(permissionMsg.ID, permissionMsg.Name)

	return nil
}

func roleHandler(permissionService *Service, msg []byte) error {
	roleMsg := RoleMsg{}

	err := json.Unmarshal(msg, &roleMsg)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling: %v: %v", string(msg), err)
		fmt.Println(err)
		return err
	}

	permissionService.SetRole(roleMsg.ID, roleMsg.Permissions)
	permissionService.BuildAllUserPermissionData()

	return nil
}

func userHandler(permissionService *Service, msg []byte) error {
	userMsg := UserMsg{}

	err := json.Unmarshal(msg, &userMsg)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling: %v: %v", string(msg), err)
		fmt.Println(err)
		return err
	}

	permissionService.SetUser(userMsg.ID, userMsg.Roles)
	permissionService.BuildUserPermissionData(userMsg.ID)

	return nil
}

func tokenHandler(permissionService *Service, msg []byte) error {
	loginSuccess := LoginSuccessMsg{}

	err := json.Unmarshal(msg, &loginSuccess)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling: %v: %v", string(msg), err)
		fmt.Println(err)
		return err
	}

	expiresAt, err := time.Parse(time.RFC3339Nano, loginSuccess.AccessTokenExpiresAt)
	if err != nil {
		err = fmt.Errorf("Error parsing time: %v: %v", string(msg), err)
		fmt.Println(err)
		return err
	}

	permissionService.SetToken(loginSuccess.AccessToken, loginSuccess.UserID, expiresAt)
	return nil
}

func AddAuthEventsHandlers(nsqEventbus *eventbus.NsqEventBus, permissionService *Service) {
	hostname, err := os.Hostname()
	if err != nil {
		panic("Could not determine Hostname")
	}

	channelName := "item_" + hostname

	pHandler := func(msg []byte) error { return permissionHandler(permissionService, msg) }
	rHandler := func(msg []byte) error { return roleHandler(permissionService, msg) }
	uHandler := func(msg []byte) error { return userHandler(permissionService, msg) }
	tHandler := func(msg []byte) error { return tokenHandler(permissionService, msg) }

	nsqEventbus.On("permission.created", channelName, pHandler)
	nsqEventbus.On("permission.updated", channelName, pHandler)

	nsqEventbus.On("role.created", channelName, rHandler)
	nsqEventbus.On("role.updated", channelName, rHandler)

	nsqEventbus.On("user.created", channelName, uHandler)
	nsqEventbus.On("user.updated", channelName, uHandler)

	nsqEventbus.On("login.success", channelName, tHandler)
}
