// demo of Tk9 (go tk) toplevel windows hiding/showing

package main

import (
	"fortio.org/cli"
	"fortio.org/log"
	tk "modernc.org/tk9.0"
)

func main() {
	cli.Main()
	log.Infof("Hello, world!")
	other := tk.Toplevel()
	tk.Pack(tk.Button(tk.Txt("Show"), tk.Command(func() { tk.WmDeiconify(other.Window) })))
	tk.Pack(other.Button(tk.Txt("Hide"), tk.Command(func() { tk.WmWithdraw(other.Window) })))
	tk.App.Wait()
}
