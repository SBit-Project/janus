package transformer

import (
	"github.com/SBit-Project/janus/pkg/notifier"
	"github.com/labstack/echo"
)

func getNotifier(c echo.Context) *notifier.Notifier {
	storedValue := c.Get("notifier")
	notifier, ok := storedValue.(*notifier.Notifier)
	if !ok {
		return nil
	}
	return notifier
}
