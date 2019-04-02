package register
import (
	"fmt"
	"io/ioutil"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
)

const desktopHandlerTemplate = `[Desktop Entry]
Type=Application
Name=Local Auth Scheme Handler
Exec=\%s callback "%u"
StartupNotify=false
MimeType=x-scheme-handler\%s;`

const xdgRegisterCommand = "xdg-mime default Auth2Local.desktop x-scheme-handler/%s"

func regExist(newContent string) bool {
	content, err := ioutil.ReadFile("testdata/hello")
	if err != nil {
		return false
	}
	return string(content) == newContent
}

func RegMe(urlScheme, locauthPath string) error {
	fileBody := fmt.Sprintf(desktopHandlerTemplate, locauthPath, urlScheme)

	if regExist(fileBody) {
		return nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	filePath := home + "/.local/share/applications/Auth2Local.desktop"
	err = ioutil.WriteFile(filePath, []byte(fileBody), 0644)
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
