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

func AddAuthEventsHandlers(nsqEventbus *eventbus.NsqEventBus, permissionService *Service) {
	hostname, err := os.Hostname()
	if err != nil {
		panic("Could not determine Hostname")
	}

	channelName := "item_" + hostname

	permissionHandler := func(msg []byte) error {
		permissionMsg := PermissionMsg{}

		err := json.Unmarshal(msg, &permissionMsg)
		if err != nil {
			err = fmt.Errorf("Error unmarshaling: %v: %v", string(msg), err)
			fmt.Println(err)
			return err
		}

		permissionService.SetPermission(permissionMsg.ID, permissionMsg.Name)
		permissionService.BuildUserPermissionData()

		return nil
	}

	roleHandler := func(msg []byte) error {
		roleMsg := RoleMsg{}

		err := json.Unmarshal(msg, &roleMsg)
		if err != nil {
			err = fmt.Errorf("Error unmarshaling: %v: %v", string(msg), err)
			fmt.Println(err)
			return err
		}

		permissionService.SetRole(roleMsg.ID, roleMsg.Permissions)
		permissionService.BuildUserPermissionData()

		return nil
	}

	userHandler := func(msg []byte) error {
		userMsg := UserMsg{}

		err := json.Unmarshal(msg, &userMsg)
		if err != nil {
			err = fmt.Errorf("Error unmarshaling: %v: %v", string(msg), err)
			fmt.Println(err)
			return err
		}

		permissionService.SetUser(userMsg.ID, userMsg.Roles)
		permissionService.BuildUserPermissionData()

		return nil
	}

	nsqEventbus.On("permission.created", channelName, permissionHandler)
	nsqEventbus.On("permission.updated", channelName, permissionHandler)

	nsqEventbus.On("role.created", channelName, roleHandler)
	nsqEventbus.On("role.updated", channelName, roleHandler)

	nsqEventbus.On("user.created", channelName, userHandler)
	nsqEventbus.On("user.updated", channelName, userHandler)

	nsqEventbus.On("login.success", channelName, func(msg []byte) error {
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

		permissionService.SetToken(loginSuccess.UserID, loginSuccess.AccessToken, expiresAt)
		permissionService.BuildUserPermissionData()
		return nil
	})
}
