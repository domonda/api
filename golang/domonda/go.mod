module github.com/domonda/api/golang/domonda

go 1.24.0

tool github.com/ungerik/go-enum

require (
	github.com/domonda/go-types v0.0.0-20260115133137-07f43dd1f81f
	github.com/ungerik/go-fs v0.0.0-20260118110456-0ae82a14cadb
)

// Pinned versions to avoid breaking updates
replace github.com/olekukonko/tablewriter => github.com/olekukonko/tablewriter v0.0.5 // Don't upgrade this to v1, it breaks the build!

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a // indirect
	github.com/clipperhouse/uax29/v2 v2.6.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/domonda/go-errs v0.0.0-20260113110342-222f906bd7a6 // indirect
	github.com/domonda/go-pretty v0.0.0-20260112082908-96fe37692898 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gogs/chardet v0.0.0-20211120154057-b7413eaefb8f // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/jhillyerd/enmime v1.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.9.1 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/xattr v0.4.12 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/teamwork/tnef v0.0.0-20200108124832-7deabccfdb32 // indirect
	github.com/ungerik/go-astvisit v0.0.0-20251017171216-b7bb0384dd33 // indirect
	github.com/ungerik/go-enum v0.0.0-20251216115906-f928944aa546 // indirect
	github.com/ungerik/go-reflection v0.0.0-20251017081454-aea4ca25282d // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/xurls/v2 v2.6.0 // indirect
)
