package main

import (
	"fmt"
	"html"
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/damian-sztankowski/lazygcp/internal/gcp"
	"github.com/damian-sztankowski/lazygcp/internal/tui"
)

type demoScene struct {
	Badge   string
	Caption string
	Frame   string
	DelayCs int // 100ths of a second for GIF
}

func main() {
	cols := 132
	rows := 36
	os.Setenv("COLUMNS", strconv.Itoa(cols+1))
	os.Setenv("LINES", strconv.Itoa(rows))

	proj := gcp.BuildMockState("acme-fintech-prod-01", "sre-lead@acme-cloud.io")
	app := tui.NewApp(proj)
	app.Width = cols
	app.Height = rows

	var scenes []demoScene

	// Scene 1: Module 1 — Cloud Run & Functions (Traffic Split + Leaked Secret Detection)
	app.ActiveCat = gcp.CatCloudRun
	app.SelectedIdx = 0
	app.InspectorTab = tui.TabOverview
	app.StatusBanner = "Module [1] Cloud Run: Press [y] for Knative YAML, [s] for Container Shell, [x] to Unmask Secrets, [r] to Rollback."
	scenes = append(scenes, demoScene{
		Badge:   "1/8 • [1] Cloud Run & Functions",
		Caption: "Inspect revision traffic splits, environment variables, leaked plaintext secrets [x], and Knative YAML [y]",
		Frame:   app.RenderFrame(),
		DelayCs: 240,
	})

	// Scene 2: Module 1 — Cloud Run Knative Service YAML Manifest ([y])
	app.InspectorTab = tui.TabServiceYAML
	app.StatusBanner = "Viewing Cloud Run / Function YAML Manifest. Use [ and ] to scroll."
	scenes = append(scenes, demoScene{
		Badge:   "2/8 • [y] Live Cloud Run Service YAML",
		Caption: "Press [y] on any Cloud Run service or Cloud Function to inspect its full Knative service.yaml manifest",
		Frame:   app.RenderFrame(),
		DelayCs: 230,
	})

	// Scene 3: Module 2 — Live Project IAM & Security Assessment ([z] Blast-Radius)
	app.ActiveCat = gcp.CatIAMZero
	app.SelectedIdx = 1
	app.InspectorTab = tui.TabIAMZero
	app.StatusBanner = "Module [2] IAM & Security Assessment: Live getIamPolicy & Service Account blast-radius analysis."
	scenes = append(scenes, demoScene{
		Badge:   "3/8 • [2] Live Project IAM & Blast-Radius",
		Caption: "Audits live project IAM bindings & Service Accounts — showing granted vs. minimal workload permissions [z]",
		Frame:   app.RenderFrame(),
		DelayCs: 250,
	})

	// Scene 4: Module 3 — GKE & Compute VMs: SSH Mode Selector Modal (Direct SSH vs IAP Tunnel)
	app.ActiveCat = gcp.CatGKECompute
	app.SelectedIdx = 1 // Running Compute VM
	app.InspectorTab = tui.TabOverview
	app.ShowSSHModal = true
	app.SSHModalIdx = 1
	app.SSHModalPromptIAP = false
	app.StatusBanner = "Choose SSH mode for 'legacy-jumpbox-vm-01': [1] Direct SSH, [2] IAP Tunnel SSH, or [3] Configure IAP Best Practices."
	scenes = append(scenes, demoScene{
		Badge:   "4/8 • [3] VM SSH Mode Selector (Direct vs IAP)",
		Caption: "Press [s] or [S] on any running Compute VM to choose Direct SSH or Zero-Trust IAP Tunnel SSH",
		Frame:   app.RenderFrame(),
		DelayCs: 250,
	})

	// Scene 5: Module 3 — IAP Best-Practices Auto-Configuration Prompt
	app.ShowSSHModal = true
	app.SSHModalIdx = 1
	app.SSHModalPromptIAP = true
	app.StatusBanner = "IAP TCP Forwarding is not detected for this VPC. Press [y / Enter] to configure IAP per Google Cloud Best Practices!"
	scenes = append(scenes, demoScene{
		Badge:   "5/8 • [3] 1-Key IAP Best-Practices Auto-Setup",
		Caption: "Automatically enables iap.googleapis.com, creates 35.235.240.0/20 -> tcp:22 firewall rule & binds IAM role",
		Frame:   app.RenderFrame(),
		DelayCs: 260,
	})

	// Scene 6: Module 3 — 1:1 Reverse-Engineered Terraform HCL for Compute VM
	app.ShowSSHModal = false
	app.SSHModalPromptIAP = false
	app.InspectorTab = tui.TabTerraform
	app.StatusBanner = "1:1 Reverse-Engineered Terraform HCL (mapped from live GCP state)! Press [w] to write .tf file."
	scenes = append(scenes, demoScene{
		Badge:   "6/8 • [t] 1:1 Live Reverse-Terraform HCL",
		Caption: "Generates state-accurate Terraform HCL + matching `terraform import` command from live GCP resources [t]",
		Frame:   app.RenderFrame(),
		DelayCs: 240,
	})

	// Scene 7: Module 8 — FinOps & Cloud Billing API v1 Breakdown
	app.ActiveCat = gcp.CatFinOps
	app.SelectedIdx = 0
	app.InspectorTab = tui.TabOverview
	app.StatusBanner = "Module [8] FinOps & Cloud Billing: Live run-rate forecast, service spend bars & 1-key waste reclamation [f]."
	scenes = append(scenes, demoScene{
		Badge:   "7/8 • [8] FinOps & Cloud Billing Perspective",
		Caption: "Tracks daily/monthly/annual run-rate forecast, service family spend breakdown & reclaimable idle waste",
		Frame:   app.RenderFrame(),
		DelayCs: 250,
	})

	// Scene 8: Module 9 — Cloud Audit Time-Machine (Git Diff of Live GCP Changes)
	app.ActiveCat = gcp.CatTimeMachine
	app.SelectedIdx = 0
	app.ShowFullAuditJSON = false
	app.StatusBanner = "Module [9] Cloud Audit Time-Machine: Inspect who changed what in git-diff format or press [r] to revert."
	scenes = append(scenes, demoScene{
		Badge:   "8/8 • [9] Cloud Audit Time-Machine",
		Caption: "Streams live cloudaudit.googleapis.com events as unified git diffs with 1-keystroke rollback [r]",
		Frame:   app.RenderFrame(),
		DelayCs: 250,
	})

	assetsDir := "assets"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		panic(err)
	}

	svgPath := filepath.Join(assetsDir, "lazygcp-demo.svg")
	if err := writeAnimatedSVG(svgPath, scenes, cols, rows); err != nil {
		panic(err)
	}
	fmt.Println("✅ Generated Retina Animated SVG:", svgPath)

	gifPath := filepath.Join(assetsDir, "lazygcp-demo.gif")
	if err := writeAnimatedGIF(gifPath, scenes, cols, rows); err != nil {
		panic(err)
	}
	fmt.Println("✅ Generated Animated GIF:", gifPath)

	tapePath := filepath.Join(assetsDir, "demo.tape")
	if err := writeVHSTape(tapePath); err != nil {
		panic(err)
	}
	fmt.Println("✅ Generated VHS Tape:", tapePath)
}

type styledSpan struct {
	Text string
	FG   string
	BG   string
	Bold bool
}

type cellColor struct {
	ch   rune
	fg   uint8
	bg   uint8
	bold bool
}

// parseANSIToGrid parses an ANSI TUI frame into a 2D grid of (rune, fgPaletteIdx, bgPaletteIdx, bold)
func parseANSIToGrid(frame string, cols, rows int) [][]cellColor {
	rawLines := strings.Split(frame, "\r\n")
	if len(rawLines) == 1 {
		rawLines = strings.Split(frame, "\n")
	}
	grid := make([][]cellColor, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]cellColor, cols)
		defaultRowBG := uint8(0)
		if r == 0 || r == 1 || r == rows-1 {
			defaultRowBG = 10 // BgHeader
		} else if r == rows-2 {
			defaultRowBG = 14 // BgActive
		}
		for c := 0; c < cols; c++ {
			grid[r][c] = cellColor{ch: ' ', fg: 1, bg: defaultRowBG}
		}
		if r >= len(rawLines) {
			continue
		}
		line := rawLines[r]
		c := 0
		curFG := uint8(1) // default white/slate
		curBG := defaultRowBG
		curBold := false

		for i := 0; i < len(line) && c < cols; {
			if line[i] == '\x1b' && i+1 < len(line) && line[i+1] == '[' {
				j := i + 2
				for j < len(line) && !((line[j] >= 'A' && line[j] <= 'Z') || (line[j] >= 'a' && line[j] <= 'z')) {
					j++
				}
				if j < len(line) {
					cmd := line[j]
					params := line[i+2 : j]
					if cmd == 'm' {
						applySGR(params, &curFG, &curBG, &curBold, defaultRowBG)
					}
					i = j + 1
					continue
				}
			}
			rn, sz := utf8.DecodeRuneInString(line[i:])
			i += sz
			if rn < 32 || rn == 0xFE0F || rn == 0xFE0E || rn == 0x200D {
				continue
			}
			grid[r][c] = cellColor{ch: rn, fg: curFG, bg: curBG, bold: curBold}
			w := runeDisplayWidth(rn)
			c += w
			if w == 2 && c-1 < cols {
				grid[r][c-1] = cellColor{ch: ' ', fg: curFG, bg: curBG, bold: curBold}
			}
		}
	}
	return grid
}

func runeDisplayWidth(r rune) int {
	if r == 0xFE0F || r == 0xFE0E || r == 0x200D {
		return 0
	}
	if (r >= 0x1F300 && r <= 0x1FAFF) || (r >= 0x2600 && r <= 0x27BF) || r == 0x23EA || r == 0x23E9 {
		return 2
	}
	return 1
}

// Palette indices:
// 0: Default BG (#080c14)
// 1: Default FG (#e2e8f0)
// 2: Gray (#64748b)
// 3: Cyan (#22d3ee)
// 4: Green (#34d399)
// 5: Yellow (#fbbf24)
// 6: Red (#fb7185)
// 7: Purple (#a78bfa)
// 8: Blue (#60a5fa)
// 9: Orange (#fb923c)
// 10: Header BG (#11192c)
// 11: Select BG (#1e3a5f)
// 12: Alert BG (#46141e)
// 13: Success BG (#0f3728)
// 14: Active BG (#162038)
// 15: Badge Dark BG (#0f172a)

var demoPalette = color.Palette{
	color.RGBA{8, 12, 20, 255},     // 0: Default BG
	color.RGBA{226, 232, 240, 255}, // 1: Default FG
	color.RGBA{100, 116, 139, 255}, // 2: Gray
	color.RGBA{34, 211, 238, 255},  // 3: Cyan
	color.RGBA{52, 211, 153, 255},  // 4: Green
	color.RGBA{251, 191, 36, 255},  // 5: Yellow
	color.RGBA{251, 113, 133, 255}, // 6: Red
	color.RGBA{167, 139, 250, 255}, // 7: Purple
	color.RGBA{96, 165, 250, 255},  // 8: Blue
	color.RGBA{251, 146, 60, 255},  // 9: Orange
	color.RGBA{17, 25, 44, 255},    // 10: Header BG
	color.RGBA{30, 58, 95, 255},    // 11: Select BG
	color.RGBA{70, 20, 30, 255},    // 12: Alert BG
	color.RGBA{15, 55, 40, 255},    // 13: Success BG
	color.RGBA{22, 32, 56, 255},    // 14: Active BG
	color.RGBA{15, 23, 42, 255},    // 15: Chrome BG
}

var svgHexColors = []string{
	"#080c14", // 0
	"#e2e8f0", // 1
	"#64748b", // 2
	"#22d3ee", // 3
	"#34d399", // 4
	"#fbbf24", // 5
	"#fb7185", // 6
	"#a78bfa", // 7
	"#60a5fa", // 8
	"#fb923c", // 9
	"#11192c", // 10
	"#1e3a5f", // 11
	"#46141e", // 12
	"#0f3728", // 13
	"#162038", // 14
	"#0f172a", // 15
}

func applySGR(params string, fg, bg *uint8, bold *bool, defaultBG uint8) {
	if params == "" || params == "0" {
		*fg = 1
		*bg = defaultBG
		*bold = false
		return
	}
	parts := strings.Split(params, ";")
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "0":
			*fg = 1
			*bg = defaultBG
			*bold = false
		case "1":
			*bold = true
		case "2":
			*fg = 2
		case "7":
			*fg, *bg = *bg, *fg
		case "38":
			if i+2 < len(parts) && parts[i+1] == "5" {
				code := parts[i+2]
				switch code {
				case "243", "245":
					*fg = 2
				case "255":
					*fg = 1
				case "51":
					*fg = 3
				case "46", "48":
					*fg = 4
				case "220", "226":
					*fg = 5
				case "196":
					*fg = 6
				case "141", "201", "213":
					*fg = 7
				case "39":
					*fg = 8
				case "208":
					*fg = 9
				}
				i += 2
			}
		case "48":
			if i+2 < len(parts) && parts[i+1] == "5" {
				code := parts[i+2]
				switch code {
				case "235":
					*bg = 10
				case "24":
					*bg = 11
				case "52":
					*bg = 12
				case "22":
					*bg = 13
				case "237":
					*bg = 14
				}
				i += 2
			} else if i+4 < len(parts) && parts[i+1] == "2" {
				rStr := parts[i+2]
				switch {
				case rStr == "17":
					*bg = 10
				case rStr == "30":
					*bg = 11
				case rStr == "70":
					*bg = 12
				case rStr == "15":
					*bg = 13
				case rStr == "22":
					*bg = 14
				}
				i += 4
			}
		}
	}
}

func isBoxBorderRune(r rune) bool {
	switch r {
	case '─', '━', '│', '┃', '╭', '╮', '╰', '╯', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼', '·':
		return true
	}
	return false
}

func writeAnimatedSVG(outPath string, scenes []demoScene, cols, rows int) error {
	charW := 8.75
	charH := 19.5
	padX := 14.0
	padTop := 48.0
	padBot := 40.0
	svgW := int(float64(cols)*charW + padX*2)
	svgH := int(float64(rows)*charH + padTop + padBot)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="100%%">`, svgW, svgH))
	sb.WriteString("\n<style>\n")
	sb.WriteString(`  .term-text { font-family: 'JetBrains Mono', 'SFMono-Regular', Menlo, Consolas, 'Liberation Mono', monospace; font-size: 14.5px; white-space: pre; }` + "\n")
	sb.WriteString(`  .title-text { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 13.5px; font-weight: 600; fill: #94a3b8; }` + "\n")
	sb.WriteString(`  .caption-bar { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 13.5px; font-weight: 600; }` + "\n")

	n := len(scenes)
	totalSec := n * 3 // 3 seconds per scene = 24s cycle
	for idx := range scenes {
		startPct := float64(idx*100) / float64(n)
		endPct := float64((idx+1)*100) / float64(n)
		sb.WriteString(fmt.Sprintf("  .scene-%d { opacity: 0; animation: anim-%d %ds infinite steps(1, end); }\n", idx, idx, totalSec))
		if idx == 0 {
			sb.WriteString(fmt.Sprintf("  @keyframes anim-%d { 0%%, %.3f%% { opacity: 1; } %.3f%%, 100%% { opacity: 0; } }\n",
				idx, endPct-0.01, endPct))
		} else {
			sb.WriteString(fmt.Sprintf("  @keyframes anim-%d { 0%%, %.3f%% { opacity: 0; } %.3f%%, %.3f%% { opacity: 1; } %.3f%%, 100%% { opacity: 0; } }\n",
				idx, startPct-0.01, startPct, endPct-0.01, endPct))
		}
	}
	sb.WriteString("</style>\n")

	// Outer macOS terminal window frame
	sb.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%d" height="%d" rx="12" fill="#080c14" stroke="#1e293b" stroke-width="2"/>`+"\n", svgW, svgH))
	sb.WriteString(fmt.Sprintf(`<rect x="1" y="1" width="%d" height="38" rx="11" fill="#0f172a"/>`+"\n", svgW-2))
	sb.WriteString(`<circle cx="22" cy="20" r="6.5" fill="#fb7185"/>` + "\n")
	sb.WriteString(`<circle cx="42" cy="20" r="6.5" fill="#fbbf24"/>` + "\n")
	sb.WriteString(`<circle cx="62" cy="20" r="6.5" fill="#34d399"/>` + "\n")
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="24" text-anchor="middle" class="title-text">⚡ lazygcp — Interactive Google Cloud Terminal Cockpit (zsh — %dx%d)</text>`+"\n", svgW/2-60, cols, rows))

	for idx, sc := range scenes {
		grid := parseANSIToGrid(sc.Frame, cols, rows)
		sb.WriteString(fmt.Sprintf(`<g class="scene-%d">`+"\n", idx))

		// Top-right scene indicator pill in titlebar
		pillW := 336
		pillX := svgW - pillW - 14
		sb.WriteString(fmt.Sprintf(`  <rect x="%d" y="7" width="%d" height="24" rx="6" fill="#1e293b" stroke="#22d3ee" stroke-width="1"/>`+"\n", pillX, pillW))
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="23" text-anchor="middle" class="caption-bar" fill="#22d3ee">%s</text>`+"\n", pillX+pillW/2, html.EscapeString(sc.Badge)))

		for r := 0; r < rows; r++ {
			yTop := padTop + float64(r)*charH
			yBase := yTop + 15.0

			// 1. Emit background rectangles for runs of non-default background colors
			c := 0
			for c < cols {
				bg := grid[r][c].bg
				if bg == 0 {
					c++
					continue
				}
				startC := c
				for c < cols && grid[r][c].bg == bg {
					c++
				}
				rx := padX + float64(startC)*charW
				rw := float64(c-startC) * charW
				sb.WriteString(fmt.Sprintf(`  <rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`+"\n",
					rx, yTop, rw, charH+0.6, svgHexColors[bg]))
			}

			// 2. Emit text spans with exact column X coordinates so borders & columns align across 100% of the SVG width
			var lineSpans strings.Builder
			c = 0
			for c < cols {
				if grid[r][c].ch == ' ' || grid[r][c].ch == 0 {
					c++
					continue
				}
				startC := c
				fg := grid[r][c].fg
				bold := grid[r][c].bold
				isBorder := isBoxBorderRune(grid[r][c].ch)

				var run strings.Builder
				for c < cols {
					cell := grid[r][c]
					if cell.ch == 0 {
						break
					}
					// Break before 2+ spaces so the next non-space token gets an exact column x=... coordinate
					if cell.ch == ' ' {
						if isBorder || (c+1 < cols && grid[r][c+1].ch == ' ') {
							break
						}
					}
					if cell.fg != fg || cell.bold != bold || isBoxBorderRune(cell.ch) != isBorder {
						break
					}
					ch := cell.ch
					if runeDisplayWidth(ch) == 2 {
						ch = '◆'
					}
					run.WriteRune(ch)
					c++
				}

				txt := run.String()
				if strings.TrimSpace(txt) == "" {
					continue
				}
				escaped := html.EscapeString(txt)
				xPos := padX + float64(startC)*charW
				weight := ""
				if bold {
					weight = ` font-weight="700"`
				}
				extraAttr := ""
				runLen := utf8.RuneCountInString(txt)
				if isBorder && runLen > 1 {
					extraAttr = fmt.Sprintf(` textLength="%.1f" lengthAdjust="spacingAndGlyphs"`, float64(runLen)*charW)
				}
				lineSpans.WriteString(fmt.Sprintf(`<tspan x="%.1f" fill="%s"%s%s>%s</tspan>`,
					xPos, svgHexColors[fg], weight, extraAttr, escaped))
			}

			if lineSpans.Len() > 0 {
				sb.WriteString(fmt.Sprintf(`  <text y="%.1f" class="term-text">%s</text>`+"\n",
					yBase, lineSpans.String()))
			}
		}

		// Bottom caption footer bar
		footerY := svgH - 34
		sb.WriteString(fmt.Sprintf(`  <rect x="1" y="%d" width="%d" height="33" rx="8" fill="#0f172a"/>`+"\n", footerY, svgW-2))
		sb.WriteString(fmt.Sprintf(`  <text x="18" y="%d" class="caption-bar" fill="#34d399">▶ %s</text>`+"\n",
			footerY+21, html.EscapeString(sc.Caption)))

		sb.WriteString("</g>\n")
	}

	sb.WriteString("</svg>\n")
	return os.WriteFile(outPath, []byte(sb.String()), 0644)
}

// writeAnimatedGIF renders each demoScene into a paletted image.Paletted using a crisp 8x16 bitmap font & box-drawing renderer.
func writeAnimatedGIF(outPath string, scenes []demoScene, cols, rows int) error {
	cw, ch := 8, 16
	padX, padTop, padBot := 12, 34, 26
	imgW := cols*cw + padX*2
	imgH := rows*ch + padTop + padBot

	anim := &gif.GIF{}

	for _, sc := range scenes {
		grid := parseANSIToGrid(sc.Frame, cols, rows)
		palImg := image.NewPaletted(image.Rect(0, 0, imgW, imgH), demoPalette)

		// Fill default dark background
		fillRect(palImg, 0, 0, imgW, imgH, 0)
		// Top window bar
		fillRect(palImg, 0, 0, imgW, 28, 15)
		// Traffic lights
		fillRect(palImg, 12, 10, 8, 8, 6)
		fillRect(palImg, 26, 10, 8, 8, 5)
		fillRect(palImg, 40, 10, 8, 8, 4)
		drawBitmapString(palImg, 60, 7, "lazygcp — Interactive Google Cloud Terminal Cockpit  ["+sc.Badge+"]", 3, cw, ch)

		for r := 0; r < rows; r++ {
			py := padTop + r*ch
			for c := 0; c < cols; c++ {
				px := padX + c*cw
				cell := grid[r][c]
				if cell.bg != 0 {
					fillRect(palImg, px, py, cw, ch, cell.bg)
				}
				if cell.ch != ' ' && cell.ch != 0 {
					drawGlyph(palImg, px, py, cw, ch, cell.ch, cell.fg, cell.bold)
				}
			}
		}

		// Bottom caption strip
		fillRect(palImg, 0, imgH-22, imgW, 22, 15)
		drawBitmapString(palImg, 12, imgH-18, ">> "+sc.Caption, 4, cw, ch)

		anim.Image = append(anim.Image, palImg)
		anim.Delay = append(anim.Delay, sc.DelayCs)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

func fillRect(img *image.Paletted, x, y, w, h int, colIdx uint8) {
	maxX := x + w
	maxY := y + h
	if maxX > img.Bounds().Max.X {
		maxX = img.Bounds().Max.X
	}
	if maxY > img.Bounds().Max.Y {
		maxY = img.Bounds().Max.Y
	}
	for yy := y; yy < maxY; yy++ {
		if yy < 0 {
			continue
		}
		for xx := x; xx < maxX; xx++ {
			if xx >= 0 {
				img.SetColorIndex(xx, yy, colIdx)
			}
		}
	}
}

func drawBitmapString(img *image.Paletted, x, y int, s string, fg uint8, cw, ch int) {
	col := 0
	for _, r := range s {
		drawGlyph(img, x+col*cw, y, cw, ch, r, fg, true)
		col++
	}
}

// drawGlyph renders box-drawing characters, block elements, and ASCII glyphs into a cw x ch cell.
func drawGlyph(img *image.Paletted, x, y, cw, ch int, r rune, fg uint8, bold bool) {
	midX := x + cw/2
	midY := y + ch/2

	switch r {
	case '─', '━', '·', '-':
		for xx := x; xx < x+cw; xx++ {
			if r == '·' && xx%2 == 1 {
				continue
			}
			img.SetColorIndex(xx, midY, fg)
		}
		return
	case '│', '┃', '|':
		for yy := y; yy < y+ch; yy++ {
			img.SetColorIndex(midX, yy, fg)
		}
		return
	case '╭', '┌':
		for xx := midX; xx < x+cw; xx++ {
			img.SetColorIndex(xx, midY, fg)
		}
		for yy := midY; yy < y+ch; yy++ {
			img.SetColorIndex(midX, yy, fg)
		}
		return
	case '╮', '┐':
		for xx := x; xx <= midX; xx++ {
			img.SetColorIndex(xx, midY, fg)
		}
		for yy := midY; yy < y+ch; yy++ {
			img.SetColorIndex(midX, yy, fg)
		}
		return
	case '╰', '└':
		for xx := midX; xx < x+cw; xx++ {
			img.SetColorIndex(xx, midY, fg)
		}
		for yy := y; yy <= midY; yy++ {
			img.SetColorIndex(midX, yy, fg)
		}
		return
	case '╯', '┘':
		for xx := x; xx <= midX; xx++ {
			img.SetColorIndex(xx, midY, fg)
		}
		for yy := y; yy <= midY; yy++ {
			img.SetColorIndex(midX, yy, fg)
		}
		return
	case '█':
		fillRect(img, x, y+2, cw, ch-4, fg)
		return
	case '░':
		for yy := y + 2; yy < y+ch-2; yy += 2 {
			for xx := x; xx < x+cw; xx += 2 {
				img.SetColorIndex(xx, yy, fg)
			}
		}
		return
	}

	glyph, ok := font5x7[r]
	if !ok && r >= 'a' && r <= 'z' {
		glyph, ok = font5x7[r-32]
	}
	if !ok {
		// Diamond/dot indicator for symbols/emojis
		fillRect(img, x+2, y+5, 3, 4, fg)
		return
	}

	offX := x + 1
	offY := y + 3
	for row := 0; row < 7; row++ {
		bits := glyph[row]
		for col := 0; col < 5; col++ {
			if (bits>>(4-col))&1 == 1 {
				img.SetColorIndex(offX+col, offY+row, fg)
				if bold && col < 4 {
					img.SetColorIndex(offX+col+1, offY+row, fg)
				}
			}
		}
	}
}

var font5x7 = map[rune][7]uint8{
	'A': {0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'B': {0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E},
	'C': {0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E},
	'D': {0x1C, 0x12, 0x11, 0x11, 0x11, 0x12, 0x1C},
	'E': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F},
	'F': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x10},
	'G': {0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0F},
	'H': {0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'I': {0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E},
	'J': {0x07, 0x02, 0x02, 0x02, 0x02, 0x12, 0x0C},
	'K': {0x11, 0x12, 0x14, 0x18, 0x14, 0x12, 0x11},
	'L': {0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x1F},
	'M': {0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11},
	'N': {0x11, 0x19, 0x15, 0x13, 0x11, 0x11, 0x11},
	'O': {0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'P': {0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10},
	'Q': {0x0E, 0x11, 0x11, 0x11, 0x15, 0x12, 0x0D},
	'R': {0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11},
	'S': {0x0F, 0x10, 0x10, 0x0E, 0x01, 0x01, 0x1E},
	'T': {0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04},
	'U': {0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'V': {0x11, 0x11, 0x11, 0x11, 0x11, 0x0A, 0x04},
	'W': {0x11, 0x11, 0x11, 0x15, 0x15, 0x15, 0x0A},
	'X': {0x11, 0x11, 0x0A, 0x04, 0x0A, 0x11, 0x11},
	'Y': {0x11, 0x11, 0x0A, 0x04, 0x04, 0x04, 0x04},
	'Z': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x10, 0x1F},
	'0': {0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E},
	'1': {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E},
	'2': {0x0E, 0x11, 0x01, 0x02, 0x04, 0x08, 0x1F},
	'3': {0x1F, 0x02, 0x04, 0x02, 0x01, 0x11, 0x0E},
	'4': {0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02},
	'5': {0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E},
	'6': {0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E},
	'7': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08},
	'8': {0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E},
	'9': {0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C},
	'[': {0x0E, 0x08, 0x08, 0x08, 0x08, 0x08, 0x0E},
	']': {0x0E, 0x02, 0x02, 0x02, 0x02, 0x02, 0x0E},
	'(': {0x02, 0x04, 0x08, 0x08, 0x08, 0x04, 0x02},
	')': {0x08, 0x04, 0x02, 0x02, 0x02, 0x04, 0x08},
	'{': {0x02, 0x04, 0x04, 0x08, 0x04, 0x04, 0x02},
	'}': {0x08, 0x04, 0x04, 0x02, 0x04, 0x04, 0x08},
	':': {0x00, 0x04, 0x04, 0x00, 0x04, 0x04, 0x00},
	'.': {0x00, 0x00, 0x00, 0x00, 0x00, 0x0C, 0x0C},
	',': {0x00, 0x00, 0x00, 0x00, 0x04, 0x04, 0x08},
	'/': {0x01, 0x01, 0x02, 0x04, 0x08, 0x10, 0x10},
	'$': {0x04, 0x0F, 0x14, 0x0E, 0x05, 0x1E, 0x04},
	'%': {0x18, 0x19, 0x02, 0x04, 0x08, 0x13, 0x03},
	'+': {0x00, 0x04, 0x04, 0x1F, 0x04, 0x04, 0x00},
	'=': {0x00, 0x00, 0x1F, 0x00, 0x1F, 0x00, 0x00},
	'>': {0x10, 0x08, 0x04, 0x02, 0x04, 0x08, 0x10},
	'<': {0x01, 0x02, 0x04, 0x08, 0x04, 0x02, 0x01},
	'?': {0x0E, 0x11, 0x01, 0x02, 0x04, 0x00, 0x04},
	'!': {0x04, 0x04, 0x04, 0x04, 0x04, 0x00, 0x04},
	'@': {0x0E, 0x11, 0x17, 0x15, 0x17, 0x10, 0x0E},
	'#': {0x0A, 0x0A, 0x1F, 0x0A, 0x1F, 0x0A, 0x0A},
	'_': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1F},
	'"': {0x0A, 0x0A, 0x0A, 0x00, 0x00, 0x00, 0x00},
	'\'': {0x04, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00},
	'&': {0x0C, 0x12, 0x14, 0x08, 0x15, 0x12, 0x0D},
	'*': {0x00, 0x04, 0x15, 0x0E, 0x15, 0x04, 0x00},
}

func writeVHSTape(outPath string) error {
	tape := `# Charmbracelet VHS Tape for recording lazygcp
Output assets/lazygcp-demo.gif
Set FontSize 14
Set Width 1440
Set Height 860
Set Theme "Catppuccin Mocha"

Type "lazygcp --demo"
Enter
Sleep 2.5s

# Module 2: Live Project IAM & Blast Radius
Type "2"
Sleep 1.5s
Type "z"
Sleep 2.5s

# Module 3: GKE & Compute VMs + SSH / IAP Mode Selector
Type "3"
Sleep 1.5s
Type "j"
Sleep 1s
Type "s"
Sleep 2.5s
Type "2"
Sleep 2.5s
Escape
Sleep 1s
Type "t"
Sleep 2.5s

# Module 8: FinOps & Cloud Billing
Type "8"
Sleep 2.5s

# Module 9: Cloud Audit Time-Machine
Type "9"
Sleep 2.5s

Type "q"
`
	return os.WriteFile(outPath, []byte(tape), 0644)
}
