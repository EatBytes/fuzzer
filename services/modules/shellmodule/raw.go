package shellmodule

import (
	"github.com/eatbytes/razboy"
	"github.com/eatbytes/razboy/adapter/phpadapter"
	"github.com/eatbytes/razboy/core"
	"github.com/eatbytes/razboynik/services/kernel"
)

func Raw(kc kernel.KernelCmd, request *core.REQUEST) (kernel.KernelCmd, error) {
	var (
		rzRes *razboy.RazResponse
		err   error
	)

	request.Type = "SHELL"

	request.SHLc.Cmd = kc.GetRaw()
	request.SHLc.Scope = kc.GetScope()
	phpadapter.CreateCMD(&request.SHLc)

	request.Build()

	rzRes, err = razboy.Send(request)
	kc.SetResult(rzRes)

	return kc, err
}
