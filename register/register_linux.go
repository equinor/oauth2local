package register

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
)

const desktopHandlerTemplate = `[Desktop Entry]
Type=Application
Terminal=true
Name=Local Auth Scheme Handler
Exec=%s callback %s
StartupNotify=false
MimeType=x-scheme-handler/%s;
NoDisplay=true`

const xdgRegisterCommand = "xdg-mime default auth2local.desktop x-scheme-handler/%s"

func regExist(newContent string) bool {
	content, err := ioutil.ReadFile("testdata/hello")
	if err != nil {
		return false
	}
	return string(content) == newContent
}

func RegMe(urlScheme, locauthPath string) error {
	fileBody := fmt.Sprintf(desktopHandlerTemplate, locauthPath, "%u", urlScheme)

	if regExist(fileBody) {
		return nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	filePath := home + "/.local/share/applications/auth2local.desktop"
	err = ioutil.WriteFile(filePath, []byte(fileBody), 0744)
	if err != nil {
		return err
	}
	cmd := exec.Command("xdg-mime", "default", "Auth2Local.desktop", "x-scheme-handler/"+urlScheme)

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
