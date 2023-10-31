package http

type App struct {
	AppName string `yaml:"appName"`
	Address string `yaml:"innerAddress"`
	Port    int    `yaml:"port"`
}

func (h *App) GetAppName() string {
	return h.AppName
}

func (h *App) DAddress() string {
	return h.Address
}

func (h *App) DPort() int {
	return h.Port
}
