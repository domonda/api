module github.com/domonda/api/golang/domonda

go 1.24

tool github.com/ungerik/go-enum

require (
	github.com/domonda/go-types v0.0.0-20250808081339-6ceb00be8516
	github.com/ungerik/go-fs v0.0.0-20250625190701-a26b03a1a7ca
)

// Pinned versions to avoid breaking updates
replace github.com/olekukonko/tablewriter => github.com/olekukonko/tablewriter v0.0.5 // Don't upgrade this to v1, it breaks the build!

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a // indirect
	github.com/domonda/go-errs v0.0.0-20250603150208-71d6de0c48ea // indirect
	github.com/domonda/go-pretty v0.0.0-20250602142956-1b467adc6387 // indirect
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
	github.com/ungerik/go-astvisit v0.0.0-20250122155250-e994358a002f // indirect
	github.com/ungerik/go-enum v0.0.0-20250819122115-78bcc1aa940b // indirect
	github.com/ungerik/go-reflection v0.0.0-20250602142243-03da83aecd0d // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/xurls/v2 v2.6.0 // indirect
)
