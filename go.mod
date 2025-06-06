module github.com/MichaelBabushkin/sammy_po

go 1.23

require (
	github.com/chromedp/cdproto v0.0.0-20240202021202-6d0b6a386732 // Added chromedp
	github.com/chromedp/chromedp v0.9.5 // Added chromedp
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.3.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

// Add require block for indirect dependencies if needed by go mod tidy
// require (
//	github.com/chromedp/sysutil v1.0.0 // indirect
//	github.com/gobwas/httphead v0.1.0 // indirect
//	github.com/gobwas/pool v0.2.1 // indirect
//	github.com/gobwas/ws v1.3.2 // indirect
//	github.com/josharian/intern v1.0.0 // indirect
//	github.com/mailru/easyjson v0.7.7 // indirect
//	golang.org/x/sys v0.18.0 // indirect
// )
