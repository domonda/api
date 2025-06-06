module github.com/domonda/api/golang/domonda

go 1.24

tool github.com/ungerik/go-enum

require (
	github.com/domonda/go-types v0.0.0-20250527163512-252e849a39ce
	github.com/ungerik/go-fs v0.0.0-20250527162931-1691110c1708
)

// Pinned versions to avoid breaking updates
replace github.com/olekukonko/tablewriter => github.com/olekukonko/tablewriter v0.0.5 // Don't upgrade this to v1, it breaks the build!

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a // indirect
	github.com/domonda/go-errs v0.0.0-20250527162518-c9fdfcc032a1 // indirect
	github.com/domonda/go-pretty v0.0.0-20240110134850-17385799142f // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gogs/chardet v0.0.0-20211120154057-b7413eaefb8f // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/jhillyerd/enmime v1.3.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/teamwork/tnef v0.0.0-20200108124832-7deabccfdb32 // indirect
	github.com/ungerik/go-astvisit v0.0.0-20231019122241-2d1ef5bbb4cf // indirect
	github.com/ungerik/go-enum v0.0.0-20241119153159-5b13f22868ae // indirect
	github.com/ungerik/go-reflection v0.0.0-20240905081803-708928fe0862 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/xurls/v2 v2.6.0 // indirect
)
