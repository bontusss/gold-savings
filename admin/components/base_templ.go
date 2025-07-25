// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Base(title string) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\" dir=\"ltr\" data-nav-layout=\"vertical\" class=\"light\" data-header-styles=\"light\" data-menu-styles=\"dark\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `admin/components/base.templ`, Line: 10, Col: 18}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</title><meta content=\"A savings platform.\" name=\"description\"><meta content=\"Gold Co-op\" name=\"GoldCo-op\"><!-- Favicon --><link rel=\"shortcut icon\" href=\"/static/img/brand-logos/favicon.ico\"><!-- Main JS --><script src=\"/static/js/main.js\"></script><!-- Style Css --><link rel=\"stylesheet\" href=\"/static/css/style.css\"><!-- Simplebar Css --><link rel=\"stylesheet\" href=\"/static/libs/simplebar/simplebar.min.css\"><!-- Color Picker Css --><link rel=\"stylesheet\" href=\"/static/libs/@simonwep/pickr/themes/nano.min.css\"><!-- Vector Map Css--><link rel=\"stylesheet\" href=\"/static/libs/jsvectormap/css/jsvectormap.min.css\"><script src=\"https://unpkg.com/htmx.org@1.9.6\"></script></head><body class=\"\"><div class=\"page\"><div id=\"toast-container\" class=\"fixed top-4 right-4 z-50\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><!-- Apex Charts JS --><script src=\"/static/libs/apexcharts/apexcharts.min.js\"></script><!-- Chartjs Chart JS --><script src=\"/static/libs/chart.js/chart.min.js\"></script><!-- Index JS --><script src=\"/static/js/index.js\"></script><!-- Back To Top --><div class=\"scrollToTop\"><span class=\"arrow\"><i class=\"ri-arrow-up-s-fill text-xl\"></i></span></div><div id=\"responsive-overlay\"></div><!-- popperjs --><script src=\"/static/libs/@popperjs/core/umd/popper.min.js\"></script><!-- Color Picker JS --><script src=\"/static/libs/@simonwep/pickr/pickr.es5.min.js\"></script><!-- sidebar JS --><script src=\"/static/js/defaultmenu.js\"></script><!-- sticky JS --><script src=\"/static/js/sticky.js\"></script><!-- Switch JS --><script src=\"/static/js/switch.js\"></script><!-- Preline JS --><script src=\"/static/libs/preline/preline.js\"></script><!-- Simplebar JS --><script src=\"/static/libs/simplebar/simplebar.min.js\"></script><!-- Custom JS --><script src=\"/static/js/custom.js\"></script><script>\n        function showToast(message, type = \"success\") {\n            const toast = document.createElement(\"div\");\n            toast.className = `toast-message bg-${type === \"success\" ? \"green\" : \"red\"}-500 text-white px-4 py-2 rounded shadow mb-2 animate-fade-in`;\n            toast.innerText = message;\n            document.getElementById(\"toast-container\").appendChild(toast);\n            setTimeout(() => {\n                toast.classList.add(\"animate-fade-out\");\n                setTimeout(() => toast.remove(), 500);\n            }, 3000);\n        }\n\n        // Listen for custom htmx events\n        document.body.addEventListener(\"htmx:afterOnLoad\", function (evt) {\n            const msg = evt.detail.xhr.getResponseHeader(\"X-Toast\");\n            const type = evt.detail.xhr.getResponseHeader(\"X-Toast-Type\");\n            if (msg) showToast(msg, type || \"success\");\n        });\n    </script><style>\n        @keyframes fadeIn {\n            from {\n                opacity: 0;\n                transform: translateY(-10px);\n            }\n\n            to {\n                opacity: 1;\n                transform: translateY(0);\n            }\n        }\n\n        @keyframes fadeOut {\n            from {\n                opacity: 1;\n            }\n\n            to {\n                opacity: 0;\n                transform: translateY(-10px);\n            }\n        }\n\n        .animate-fade-in {\n            animation: fadeIn 0.3s;\n        }\n\n        .animate-fade-out {\n            animation: fadeOut 0.5s forwards;\n        }\n    </style></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
