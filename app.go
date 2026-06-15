package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	window fyne.Window
	config Config

	nameEntry        *widget.Entry
	dpiSlider        *widget.Slider
	dpiLabel         *widget.Label
	smartShiftCheck  *widget.Check
	smartShiftSlider *widget.Slider
	smartShiftLabel  *widget.Label
	hiresCheck       *widget.Check
	hiresInvertCheck *widget.Check
	hiresTargetCheck *widget.Check

	thumbDivertCheck *widget.Check
	thumbInvertCheck *widget.Check
	thumbLeftLabel   *widget.Label
	thumbRightLabel  *widget.Label
	thumbTapLabel    *widget.Label
	thumbLeftBtn     *widget.Button
	thumbRightBtn    *widget.Button
	thumbTapBtn      *widget.Button

	buttonWidgets []*ButtonRowWidget

	previewEntry *widget.Entry
	statusLabel  *widget.Label
	langSelect   *widget.Select
}

type ButtonRowWidget struct {
	container *fyne.Container
	nameLabel *widget.Label
	actionLabel *widget.Label
	editBtn   *widget.Button
}

func NewApp(w fyne.Window) *App {
	a := &App{
		window: w,
		config: DefaultConfig(),
	}
	w.SetTitle("Logitech MX Master Configuration")
	w.Resize(fyne.NewSize(800, 600))
	return a
}

func (a *App) BuildUI() fyne.CanvasObject {
	a.langSelect = widget.NewSelect([]string{"English", "Português"}, func(s string) {
		if s == "Português" {
			currentLang = LangPT
		} else {
			currentLang = LangEN
		}
		a.refreshButtonLabels()
		a.refreshPreview()
	})
	a.langSelect.SetSelected("English")

	generalTab := container.NewVBox(
		widget.NewLabelWithStyle("General Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.buildGeneralSection(),
		widget.NewSeparator(),
		a.buildSmartShiftSection(),
		widget.NewSeparator(),
		a.buildHiResSection(),
	)

	buttonsTab := container.NewVBox(
		widget.NewLabelWithStyle("Button Configuration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.buildButtonsSection(),
	)

	thumbTab := container.NewVBox(
		widget.NewLabelWithStyle("Thumbwheel Configuration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.buildThumbwheelSection(),
	)

	previewTab := container.NewVBox(
		widget.NewLabelWithStyle("Configuration Preview", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.previewEntry,
	)

	serviceTab := a.buildServiceSection()

	topBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel("Language:"),
		a.langSelect,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("General", container.NewVBox(generalTab, layout.NewSpacer())),
		container.NewTabItem("Buttons", container.NewVBox(buttonsTab, layout.NewSpacer())),
		container.NewTabItem("Thumbwheel", container.NewVBox(thumbTab, layout.NewSpacer())),
		container.NewTabItem("Preview", container.NewVBox(previewTab, layout.NewSpacer())),
		container.NewTabItem("Service", serviceTab),
	)

	content := container.NewBorder(topBar, nil, nil, nil, tabs)
	return content
}

func (a *App) buildGeneralSection() fyne.CanvasObject {
	a.nameEntry = widget.NewEntry()
	a.nameEntry.SetText(a.config.Name)
	a.nameEntry.OnChanged = func(s string) {
		a.config.Name = s
		a.refreshPreview()
	}

	a.dpiLabel = widget.NewLabel(fmt.Sprintf("DPI: %d", a.config.DPI))
	a.dpiSlider = widget.NewSlider(200, 4000)
	a.dpiSlider.Value = float64(a.config.DPI)
	a.dpiSlider.Step = 100
	a.dpiSlider.OnChanged = func(v float64) {
		val := int(v)
		a.config.DPI = val
		a.dpiLabel.SetText(fmt.Sprintf("DPI: %d", val))
		a.refreshPreview()
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Device Name", Widget: a.nameEntry},
			{Text: "", Widget: a.dpiLabel},
			{Text: "", Widget: a.dpiSlider},
		},
	}
	return form
}

func (a *App) buildSmartShiftSection() fyne.CanvasObject {
	a.smartShiftCheck = widget.NewCheck("Enable SmartShift", func(b bool) {
		a.config.SmartShift.On = b
		a.refreshPreview()
	})
	a.smartShiftCheck.SetChecked(a.config.SmartShift.On)

	a.smartShiftLabel = widget.NewLabel(fmt.Sprintf("Threshold: %d", a.config.SmartShift.Threshold))
	a.smartShiftSlider = widget.NewSlider(1, 255)
	a.smartShiftSlider.Value = float64(a.config.SmartShift.Threshold)
	a.smartShiftSlider.Step = 1
	a.smartShiftSlider.OnChanged = func(v float64) {
		val := int(v)
		a.config.SmartShift.Threshold = val
		a.smartShiftLabel.SetText(fmt.Sprintf("Threshold: %d", val))
		a.refreshPreview()
	}

	box := container.NewVBox(
		widget.NewLabelWithStyle("SmartShift", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.smartShiftCheck,
		a.smartShiftLabel,
		a.smartShiftSlider,
	)
	return box
}

func (a *App) buildHiResSection() fyne.CanvasObject {
	a.hiresCheck = widget.NewCheck("Enable Hi-Res Scrolling", func(b bool) {
		a.config.HiResScroll.Hires = b
		a.refreshPreview()
	})
	a.hiresCheck.SetChecked(a.config.HiResScroll.Hires)

	a.hiresInvertCheck = widget.NewCheck("Invert Scroll", func(b bool) {
		a.config.HiResScroll.Invert = b
		a.refreshPreview()
	})
	a.hiresInvertCheck.SetChecked(a.config.HiResScroll.Invert)

	a.hiresTargetCheck = widget.NewCheck("HID++ Target (remap scroll wheel)", func(b bool) {
		a.config.HiResScroll.Target = b
		a.refreshPreview()
	})
	a.hiresTargetCheck.SetChecked(a.config.HiResScroll.Target)

	box := container.NewVBox(
		widget.NewLabelWithStyle("Hi-Res Scrolling", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.hiresCheck,
		a.hiresInvertCheck,
		a.hiresTargetCheck,
	)
	return box
}

func (a *App) buildThumbwheelSection() fyne.CanvasObject {
	a.thumbDivertCheck = widget.NewCheck("Divert Thumbwheel (handle in logid)", func(b bool) {
		a.config.Thumbwheel.Divert = b
		if a.thumbLeftBtn == nil {
			return
		}
		a.thumbLeftBtn.Disable()
		a.thumbRightBtn.Disable()
		a.thumbTapBtn.Disable()
		if b {
			a.thumbLeftBtn.Enable()
			a.thumbRightBtn.Enable()
			a.thumbTapBtn.Enable()
		}
		a.refreshPreview()
	})
	a.thumbDivertCheck.SetChecked(a.config.Thumbwheel.Divert)

	a.thumbInvertCheck = widget.NewCheck("Invert Direction", func(b bool) {
		a.config.Thumbwheel.Invert = b
		a.refreshPreview()
	})
	a.thumbInvertCheck.SetChecked(a.config.Thumbwheel.Invert)

	a.thumbLeftLabel = widget.NewLabel(a.formatThumbAction("thumbwheel_left", a.config.Thumbwheel.Left))
	a.thumbLeftBtn = widget.NewButton("Edit", func() {
		a.showActionDialog("thumbwheel_left", &a.config.Thumbwheel.Left)
	})

	a.thumbRightLabel = widget.NewLabel(a.formatThumbAction("thumbwheel_right", a.config.Thumbwheel.Right))
	a.thumbRightBtn = widget.NewButton("Edit", func() {
		a.showActionDialog("thumbwheel_right", &a.config.Thumbwheel.Right)
	})

	a.thumbTapLabel = widget.NewLabel(a.formatThumbAction("thumbwheel_tap", a.config.Thumbwheel.Tap))
	a.thumbTapBtn = widget.NewButton("Edit", func() {
		a.showActionDialog("thumbwheel_tap", &a.config.Thumbwheel.Tap)
	})

	leftRow := container.NewGridWithColumns(3,
		widget.NewLabel(ButtonName("thumbwheel_left", currentLang)),
		a.thumbLeftLabel, a.thumbLeftBtn,
	)
	rightRow := container.NewGridWithColumns(3,
		widget.NewLabel(ButtonName("thumbwheel_right", currentLang)),
		a.thumbRightLabel, a.thumbRightBtn,
	)
	tapRow := container.NewGridWithColumns(3,
		widget.NewLabel(ButtonName("thumbwheel_tap", currentLang)),
		a.thumbTapLabel, a.thumbTapBtn,
	)

	box := container.NewVBox(
		widget.NewLabelWithStyle("Thumbwheel", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.thumbDivertCheck,
		a.thumbInvertCheck,
		widget.NewSeparator(),
		leftRow,
		rightRow,
		tapRow,
	)
	return box
}

func (a *App) buildButtonsSection() fyne.CanvasObject {
	a.buttonWidgets = nil
	scrollContent := container.NewVBox()

	for _, btn := range a.config.Buttons {
		row := a.newButtonRow(btn)
		a.buttonWidgets = append(a.buttonWidgets, row)
		scrollContent.Add(row.container)
	}

	scroll := container.NewScroll(scrollContent)
	scroll.SetMinSize(fyne.NewSize(750, 300))
	return scroll
}

func (a *App) newButtonRow(btn ButtonConfig) *ButtonRowWidget {
	idx := -1
	for i, b := range a.config.Buttons {
		if b.CID == btn.CID {
			idx = i
			break
		}
	}

	var nameID string
	for _, bi := range MXMasterButtons {
		if bi.CID == btn.CID {
			nameID = bi.Name
			break
		}
	}

	nameLabel := widget.NewLabel(ButtonName(nameID, currentLang))
	actionLabel := widget.NewLabel(a.formatAction(btn.Action))

	editBtn := widget.NewButton("Edit", func() {
		if idx >= 0 {
			a.showButtonActionDialog(idx, nameID)
		}
	})

	row := container.NewGridWithColumns(3, nameLabel, actionLabel, editBtn)
	return &ButtonRowWidget{
		container:   row,
		nameLabel:   nameLabel,
		actionLabel: actionLabel,
		editBtn:     editBtn,
	}
}

func (a *App) showButtonActionDialog(idx int, nameID string) {
	btn := &a.config.Buttons[idx]
	title := fmt.Sprintf("Configure %s", ButtonName(nameID, currentLang))

	a.showActionEditorDialog(title, &btn.Action, func() {
		a.buttonWidgets[idx].actionLabel.SetText(a.formatAction(btn.Action))
		a.refreshPreview()
	})
}

func (a *App) showActionDialog(nameID string, action *Action) {
	title := fmt.Sprintf("Configure %s", ButtonName(nameID, currentLang))

	a.showActionEditorDialog(title, action, func() {
		switch nameID {
		case "thumbwheel_left":
			a.thumbLeftLabel.SetText(a.formatThumbAction(nameID, *action))
		case "thumbwheel_right":
			a.thumbRightLabel.SetText(a.formatThumbAction(nameID, *action))
		case "thumbwheel_tap":
			a.thumbTapLabel.SetText(a.formatThumbAction(nameID, *action))
		}
		a.refreshPreview()
	})
}

func (a *App) showActionEditorDialog(title string, action *Action, onSave func()) {
	actionTypes := []string{"None", "Keypress", "Gestures", "ToggleSmartShift", "ToggleHiresScroll", "CycleDPI", "ChangeDPI", "ChangeHost"}
	atMap := ActionTypes(currentLang)

	actionNames := make([]string, len(actionTypes))
	for i, t := range actionTypes {
		if n, ok := atMap[t]; ok {
			actionNames[i] = n
		} else {
			actionNames[i] = t
		}
	}

	typeSelect := widget.NewSelect(actionNames, nil)
	for i, t := range actionTypes {
		if t == action.Type {
			typeSelect.SetSelected(actionNames[i])
			break
		}
	}
	if typeSelect.Selected == "" {
		typeSelect.SetSelected(actionNames[0])
	}

	content := container.NewVBox(
		widget.NewLabel("Action Type:"),
		typeSelect,
	)

	keysEntry := widget.NewEntry()
	keysEntry.SetText(strings.Join(action.Keys, ", "))

	dpiEntry := widget.NewEntry()
	dpiStrs := make([]string, len(action.DPIs))
	for i, d := range action.DPIs {
		dpiStrs[i] = strconv.Itoa(d)
	}
	dpiEntry.SetText(strings.Join(dpiStrs, ", "))

	hostEntry := widget.NewEntry()
	hostEntry.SetText(action.Host)

	gestureContent := container.NewVBox()

	updateContent := func() {
		content.Objects = []fyne.CanvasObject{
			widget.NewLabel("Action Type:"),
			typeSelect,
		}

		var selectedType string
		for i, t := range actionTypes {
			if typeSelect.Selected == actionNames[i] {
				selectedType = t
				break
			}
		}

		switch selectedType {
		case "Keypress":
			content.Add(widget.NewLabel("Keys (comma-separated, e.g. KEY_LEFTCTRL, KEY_T):"))
			content.Add(keysEntry)
		case "Gestures":
			content.Add(widget.NewLabel("Configure gestures in the main window"))
			content.Add(gestureContent)
		case "CycleDPI":
			content.Add(widget.NewLabel("DPI values (comma-separated, e.g. 400, 800, 1000, 1200):"))
			content.Add(dpiEntry)
		case "ChangeDPI":
			content.Add(widget.NewLabel("DPI increment:"))
			content.Add(dpiEntry)
		case "ChangeHost":
			content.Add(widget.NewLabel("Host (number, 'next', or 'prev'):"))
			content.Add(hostEntry)
		}
		content.Refresh()
	}

	typeSelect.OnChanged = func(s string) {
		updateContent()
	}

	updateContent()

	dialog.ShowCustomConfirm(title, "Save", "Cancel", content, func(b bool) {
		if !b {
			return
		}

		var selectedType string
		for i, t := range actionTypes {
			if typeSelect.Selected == actionNames[i] {
				selectedType = t
				break
			}
		}

		action.Type = selectedType
		action.Keys = nil
		action.DPIs = nil
		action.Host = ""

		switch selectedType {
		case "Keypress":
			parts := strings.Split(keysEntry.Text, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					action.Keys = append(action.Keys, p)
				}
			}
		case "CycleDPI":
			parts := strings.Split(dpiEntry.Text, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if v, err := strconv.Atoi(p); err == nil {
					action.DPIs = append(action.DPIs, v)
				}
			}
		case "ChangeDPI":
			if v, err := strconv.Atoi(strings.TrimSpace(dpiEntry.Text)); err == nil {
				action.Inc = v
			}
		case "ChangeHost":
			action.Host = strings.TrimSpace(hostEntry.Text)
		}

		onSave()
	}, a.window)
}

func (a *App) buildServiceSection() fyne.CanvasObject {
	statusLabel := widget.NewLabel("")
	outputEntry := widget.NewMultiLineEntry()
	outputEntry.SetMinRowsVisible(8)
	outputEntry.Disable()

	configPathEntry := widget.NewEntry()
	configPathEntry.SetText("/etc/logid.cfg")

	updateStatus := func() {
		if IsServiceRunning() {
			statusLabel.SetText("Status: Running")
		} else {
			statusLabel.SetText("Status: Stopped")
		}
	}

	refreshBtn := widget.NewButton("Refresh Status", func() {
		updateStatus()
	})

	saveConfigBtn := widget.NewButton("Save Configuration", func() {
		content := a.config.Generate()
		path := configPathEntry.Text
		err := WriteConfigFile(content, path)
		if err != nil {
			outputEntry.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		outputEntry.SetText(fmt.Sprintf("Configuration saved to %s\n\n%s", path, content))
	})

	installServiceBtn := widget.NewButton("Install & Start Service", func() {
		outputEntry.SetText("Installing systemd service...\n")

		content := a.config.Generate()
		if err := WriteConfigFile(content, configPathEntry.Text); err != nil {
			outputEntry.SetText(fmt.Sprintf("Error writing config: %v", err))
			return
		}

		servicePath := "/etc/systemd/system/logid.service"
		if err := WriteServiceFile(servicePath); err != nil {
			outputEntry.SetText(fmt.Sprintf("Error writing service: %v", err))
			return
		}

		out, err := EnableAndStartService()
		if err != nil {
			outputEntry.SetText(fmt.Sprintf("%s\nError: %v", out, err))
			return
		}
		outputEntry.SetText(out)
		updateStatus()
	})

	stopServiceBtn := widget.NewButton("Stop Service", func() {
		out, err := StopService()
		if err != nil {
			outputEntry.SetText(fmt.Sprintf("Error: %v\n%s", err, out))
			return
		}
		outputEntry.SetText(out)
		updateStatus()
	})

	restartServiceBtn := widget.NewButton("Restart Service", func() {
		out, err := RestartService()
		if err != nil {
			outputEntry.SetText(fmt.Sprintf("Error: %v\n%s", err, out))
			return
		}
		outputEntry.SetText(out)
		updateStatus()
	})

	servicePreview := widget.NewMultiLineEntry()
	servicePreview.SetText(serviceContent)
	servicePreview.Disable()
	servicePreview.SetMinRowsVisible(10)

	updateStatus()

	box := container.NewVBox(
		widget.NewLabelWithStyle("Service Management", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel("Config path:"), configPathEntry),
		refreshBtn,
		statusLabel,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Actions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		saveConfigBtn,
		installServiceBtn,
		container.NewHBox(stopServiceBtn, restartServiceBtn),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Systemd Service Unit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		servicePreview,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Output", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		outputEntry,
	)

	scroll := container.NewScroll(box)
	scroll.SetMinSize(fyne.NewSize(750, 500))
	return scroll
}

func (a *App) refreshPreview() {
	if a.previewEntry == nil {
		a.previewEntry = widget.NewMultiLineEntry()
		a.previewEntry.Disable()
		a.previewEntry.SetMinRowsVisible(20)
	}
	a.previewEntry.SetText(a.config.Generate())
}

func (a *App) refreshButtonLabels() {
	if a.buttonWidgets == nil {
		return
	}
	for i, row := range a.buttonWidgets {
		if i < len(a.config.Buttons) {
			var nameID string
			for _, bi := range MXMasterButtons {
				if bi.CID == a.config.Buttons[i].CID {
					nameID = bi.Name
					break
				}
			}
			if nameID != "" {
				row.nameLabel.SetText(ButtonName(nameID, currentLang))
			}
		}
	}
}

func (a *App) formatAction(a2 Action) string {
	atMap := ActionTypes(currentLang)
	name, ok := atMap[a2.Type]
	if !ok {
		name = a2.Type
	}
	if name == "" {
		name = "None"
	}
	switch a2.Type {
	case "Keypress":
		return fmt.Sprintf("%s: %s", name, strings.Join(a2.Keys, " + "))
	case "Gestures":
		return fmt.Sprintf("%s (%d directions)", name, len(a2.Gestures))
	case "CycleDPI":
		dpiStrs := make([]string, len(a2.DPIs))
		for i, d := range a2.DPIs {
			dpiStrs[i] = strconv.Itoa(d)
		}
		return fmt.Sprintf("%s: [%s]", name, strings.Join(dpiStrs, ", "))
	case "ChangeDPI":
		return fmt.Sprintf("%s: %+d", name, a2.Inc)
	case "ChangeHost":
		return fmt.Sprintf("%s: %s", name, a2.Host)
	default:
		return name
	}
}

func (a *App) formatThumbAction(nameID string, a2 Action) string {
	if a2.Type == "" || a2.Type == "None" {
		return "None"
	}
	return a.formatAction(a2)
}

func (a *App) Run() {
	a.refreshPreview()
	a.window.SetContent(a.BuildUI())
	a.window.ShowAndRun()
}
