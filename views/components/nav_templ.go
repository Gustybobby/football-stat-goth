// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Nav() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<header class=\"h-16 w-full flex justify-between items-center px-4 sticky top-0 bg-gradient-to-r from-secondary-background to-secondary text-primary-background shadow-lg z-50 border-y-4 border-b-secondary border-t-secondary-background\"><a class=\"flex items-center space-x-2\" href=\"/\"><img src=\"/public/icon.png\" class=\"size-8\"><h1 class=\"font-bold text-3xl\">PLaymaker</h1></a><nav class=\"space-x-4 font-semibold hidden sm:flex\"><a href=\"/\" class=\"hover:underline\">Home</a> <a href=\"/standings\" class=\"hover:underline\">Standings</a> <a href=\"/clubs\" class=\"hover:underline\">Clubs</a> <a href=\"/\" class=\"hover:underline\">Players</a> <a href=\"/\" class=\"hover:underline\">Fantasy</a> <a href=\"/signup\" class=\"hover:underline\">Signin</a></nav></header>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
