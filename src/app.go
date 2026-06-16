package main

import (
	"fmt"
	"log"
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

	nameSelect       *widget.Select
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
	deviceInfoLabel *widget.Label

	previewEntry *widget.Entry
	statusLabel  *widget.Label
	langSelect   *widget.Select
}

type ButtonRowWidget struct {
	container   *fyne.Container
	nameLabel   *widget.Label
	actionLabel *widget.Label
	editBtn     *widget.Button
}

func NewApp(w fyne.Window) *App {
	cfg := DefaultConfig()
	if loaded, err := LoadConfigFromFile("/etc/logid.cfg"); err == nil {
		cfg = loaded
	}
	a := &App{
		window: w,
		config: cfg,
	}
	w.SetTitle(Translate("Configuração do Logitech MX Master", "Logitech MX Master Configuration", currentLang))
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
		widget.NewLabelWithStyle(Translate("Configurações Gerais", "General Settings", currentLang), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.buildGeneralSection(),
		widget.NewSeparator(),
		a.buildSmartShiftSection(),
		widget.NewSeparator(),
		a.buildHiResSection(),
	)

	buttonsTab := container.NewVBox(
		widget.NewLabelWithStyle(Translate("Configuração de Botões", "Button Configuration", currentLang), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.buildButtonsSection(),
	)

	thumbTab := container.NewVBox(
		widget.NewLabelWithStyle(Translate("Configuração da Roda Lateral", "Thumbwheel Configuration", currentLang), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.buildThumbwheelSection(),
	)

	previewTab := container.NewVBox(
		widget.NewLabelWithStyle(Translate("Visualização da Configuração", "Configuration Preview", currentLang), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		a.previewEntry,
	)

	serviceTab := a.buildServiceSection()

	topBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel(Translate("Idioma:", "Language:", currentLang)),
		a.langSelect,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem(Translate("Geral", "General", currentLang), container.NewVBox(generalTab, layout.NewSpacer())),
		container.NewTabItem(Translate("Botões", "Buttons", currentLang), container.NewVBox(buttonsTab, layout.NewSpacer())),
		container.NewTabItem(Translate("Roda Lateral", "Thumbwheel", currentLang), container.NewVBox(thumbTab, layout.NewSpacer())),
		container.NewTabItem(Translate("Visualizar", "Preview", currentLang), container.NewVBox(previewTab, layout.NewSpacer())),
		container.NewTabItem(Translate("Salvar", "Save", currentLang), serviceTab),
	)

	content := container.NewBorder(topBar, nil, nil, nil, tabs)
	return content
}

func (a *App) buildGeneralSection() fyne.CanvasObject {
	a.nameSelect = widget.NewSelect(DeviceNames, func(s string) {
		a.config.Name = s
		a.refreshPreview()
	})
	a.nameSelect.SetSelected(a.config.Name)
	if a.nameSelect.Selected == "" {
		a.nameSelect.SetSelected(DeviceNames[0])
		a.config.Name = DeviceNames[0]
	}

	a.deviceInfoLabel = widget.NewLabelWithStyle(DeviceNameInfo(currentLang), fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

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
			{Text: Translate("Nome do Dispositivo", "Device Name", currentLang), Widget: a.nameSelect},
			{Text: "", Widget: a.deviceInfoLabel},
			{Text: "", Widget: a.dpiLabel},
			{Text: "", Widget: a.dpiSlider},
		},
	}
	return form
}

func (a *App) buildSmartShiftSection() fyne.CanvasObject {
	a.smartShiftCheck = widget.NewCheck(Translate("Ativar SmartShift", "Enable SmartShift", currentLang), func(b bool) {
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
		widget.NewLabelWithStyle(Translate("SmartShift", "SmartShift", currentLang), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.smartShiftCheck,
		a.smartShiftLabel,
		a.smartShiftSlider,
	)
	return box
}

func (a *App) buildHiResSection() fyne.CanvasObject {
	a.hiresCheck = widget.NewCheck(Translate("Ativar Scroll de Alta Resolução", "Enable Hi-Res Scrolling", currentLang), func(b bool) {
		a.config.HiResScroll.Hires = b
		a.refreshPreview()
	})
	a.hiresCheck.SetChecked(a.config.HiResScroll.Hires)

	a.hiresInvertCheck = widget.NewCheck(Translate("Inverter Scroll", "Invert Scroll", currentLang), func(b bool) {
		a.config.HiResScroll.Invert = b
		a.refreshPreview()
	})
	a.hiresInvertCheck.SetChecked(a.config.HiResScroll.Invert)

	a.hiresTargetCheck = widget.NewCheck(Translate("Destino HID++ (remapear roda de scroll)", "HID++ Target (remap scroll wheel)", currentLang), func(b bool) {
		a.config.HiResScroll.Target = b
		a.refreshPreview()
	})
	a.hiresTargetCheck.SetChecked(a.config.HiResScroll.Target)

	box := container.NewVBox(
		widget.NewLabelWithStyle(Translate("Scroll de Alta Resolução", "Hi-Res Scrolling", currentLang), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.hiresCheck,
		a.hiresInvertCheck,
		a.hiresTargetCheck,
	)
	return box
}

func (a *App) buildThumbwheelSection() fyne.CanvasObject {
	a.thumbDivertCheck = widget.NewCheck(Translate("Desviar Roda Lateral (gerenciar no logid)", "Divert Thumbwheel (handle in logid)", currentLang), func(b bool) {
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

	a.thumbInvertCheck = widget.NewCheck(Translate("Inverter Direção", "Invert Direction", currentLang), func(b bool) {
		a.config.Thumbwheel.Invert = b
		a.refreshPreview()
	})
	a.thumbInvertCheck.SetChecked(a.config.Thumbwheel.Invert)

	a.thumbLeftLabel = widget.NewLabel(a.formatThumbAction(a.config.Thumbwheel.Left))
	a.thumbLeftBtn = widget.NewButton(Translate("Editar", "Edit", currentLang), func() {
		a.showActionDialog("thumbwheel_left", &a.config.Thumbwheel.Left)
	})

	a.thumbRightLabel = widget.NewLabel(a.formatThumbAction(a.config.Thumbwheel.Right))
	a.thumbRightBtn = widget.NewButton(Translate("Editar", "Edit", currentLang), func() {
		a.showActionDialog("thumbwheel_right", &a.config.Thumbwheel.Right)
	})

	a.thumbTapLabel = widget.NewLabel(a.formatThumbAction(a.config.Thumbwheel.Tap))
	a.thumbTapBtn = widget.NewButton(Translate("Editar", "Edit", currentLang), func() {
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
		widget.NewLabelWithStyle(Translate("Roda Lateral", "Thumbwheel", currentLang), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
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

	editBtn := widget.NewButton(Translate("Editar", "Edit", currentLang), func() {
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
			a.thumbLeftLabel.SetText(a.formatThumbAction(*action))
		case "thumbwheel_right":
			a.thumbRightLabel.SetText(a.formatThumbAction(*action))
		case "thumbwheel_tap":
			a.thumbTapLabel.SetText(a.formatThumbAction(*action))
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
		widget.NewLabel(Translate("Tipo de Ação:", "Action Type:", currentLang)),
		typeSelect,
	)

	keySelects := make([]*widget.Select, 0)
	keyContainer := container.NewVBox()

	for _, k := range action.Keys {
		s := widget.NewSelect(SortedKeyDisplayNames(currentLang), nil)
		s.SetSelected(KeyDisplayName(k, currentLang))
		keySelects = append(keySelects, s)
	}
	if len(keySelects) == 0 {
		s := widget.NewSelect(SortedKeyDisplayNames(currentLang), nil)
		keySelects = append(keySelects, s)
	}

	dpiEntry := widget.NewEntry()
	dpiStrs := make([]string, len(action.DPIs))
	for i, d := range action.DPIs {
		dpiStrs[i] = strconv.Itoa(d)
	}
	dpiEntry.SetText(strings.Join(dpiStrs, ", "))

	hostEntry := widget.NewEntry()
	hostEntry.SetText(action.Host)

	gestureContent := container.NewVBox()

	var rebuildKeySelects func()

	updateContent := func() {
		content.Objects = []fyne.CanvasObject{
			widget.NewLabel(Translate("Tipo de Ação:", "Action Type:", currentLang)),
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
			content.Add(widget.NewLabel(Translate("Teclas:", "Keys:", currentLang)))
			rebuildKeySelects()
			content.Add(keyContainer)
		case "Gestures":
			content.Add(widget.NewLabel(Translate("Configure gestos na janela principal", "Configure gestures in the main window", currentLang)))
			content.Add(gestureContent)
		case "CycleDPI":
			content.Add(widget.NewLabel(Translate("Valores de DPI (separados por vírgula, ex: 400, 800, 1000, 1200):", "DPI values (comma-separated, e.g. 400, 800, 1000, 1200):", currentLang)))
			content.Add(dpiEntry)
		case "ChangeDPI":
			content.Add(widget.NewLabel(Translate("Incremento de DPI:", "DPI increment:", currentLang)))
			content.Add(dpiEntry)
		case "ChangeHost":
			content.Add(widget.NewLabel(Translate("Host (número, 'next' ou 'prev'):", "Host (number, 'next', or 'prev'):", currentLang)))
			content.Add(hostEntry)
		}
		content.Refresh()
	}

	rebuildKeySelects = func() {
		keyContainer.Objects = nil
		sortedNames := SortedKeyDisplayNames(currentLang)

		for idx, s := range keySelects {
			row := container.NewHBox(
				widget.NewLabel(fmt.Sprintf(Translate("Tecla %d:", "Key %d:", currentLang), idx+1)),
				s,
			)
			if len(keySelects) > 1 {
				removeIdx := idx
				removeBtn := widget.NewButton("-", func() {
					keySelects = append(keySelects[:removeIdx], keySelects[removeIdx+1:]...)
					rebuildKeySelects()
				})
				row.Add(removeBtn)
			}
			keyContainer.Add(row)
		}

		addBtn := widget.NewButton(Translate("+ Adicionar Tecla", "+ Add Key", currentLang), func() {
			s := widget.NewSelect(sortedNames, nil)
			keySelects = append(keySelects, s)
			rebuildKeySelects()
		})
		keyContainer.Add(addBtn)
		content.Refresh()
	}

	typeSelect.OnChanged = func(s string) {
		updateContent()
	}

	updateContent()

	dialog.ShowCustomConfirm(title, Translate("Salvar", "Save", currentLang), Translate("Cancelar", "Cancel", currentLang), content, func(b bool) {
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
			for _, s := range keySelects {
				if s.Selected == "" {
					continue
				}
				k := KeyCodeFromDisplay(s.Selected, currentLang)
				if k != "" {
					action.Keys = append(action.Keys, k)
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
			statusLabel.SetText(Translate("Status: Executando", "Status: Running", currentLang))
		} else {
			statusLabel.SetText(Translate("Status: Parado", "Status: Stopped", currentLang))
		}
	}

	refreshBtn := widget.NewButton(Translate("Atualizar Status", "Refresh Status", currentLang), func() {
		log.Printf("Refresh Status clicked")
		updateStatus()
	})

	saveConfigBtn := widget.NewButton(Translate("Salvar Configuração", "Save Configuration", currentLang), func() {
		log.Printf("Save Configuration clicked, path=%s", configPathEntry.Text)
		content := a.config.Generate()
		path := configPathEntry.Text
		err := WriteConfigFile(content, path)
		if err != nil {
			outputEntry.SetText(fmt.Sprintf(Translate("Erro: %v", "Error: %v", currentLang), err))
			dialog.ShowError(fmt.Errorf(Translate("Falha ao salvar configuração:\n%v", "Failed to save configuration:\n%v", currentLang), err), a.window)
			return
		}
		outputEntry.SetText(fmt.Sprintf(Translate("Configuração salva em %s\n\n%s", "Configuration saved to %s\n\n%s", currentLang), path, content))
		dialog.ShowInformation(Translate("Sucesso", "Success", currentLang), fmt.Sprintf(Translate("Configuração salva em %s", "Configuration saved to %s", currentLang), path), a.window)
	})

	resetConfigBtn := widget.NewButton(Translate("Redefinir para Padrão", "Reset to Default", currentLang), func() {
		log.Printf("Reset to Default clicked")
		a.config = DefaultConfig()
		a.applyConfigToUI()
		outputEntry.SetText(Translate("Configuração redefinida para valores padrão", "Configuration reset to default values", currentLang))
		dialog.ShowInformation(Translate("Sucesso", "Success", currentLang), Translate("Configuração redefinida para valores padrão", "Configuration reset to default values", currentLang), a.window)
	})

	installServiceBtn := widget.NewButton(Translate("Instalar e Iniciar Serviço", "Install & Start Service", currentLang), func() {
		log.Printf("Install & Start Service clicked")
		outputEntry.SetText(Translate("Instalando serviço systemd...\n", "Installing systemd service...\n", currentLang))

		content := a.config.Generate()
		if err := WriteConfigFile(content, configPathEntry.Text); err != nil {
			outputEntry.SetText(fmt.Sprintf(Translate("Erro ao escrever config: %v", "Error writing config: %v", currentLang), err))
			dialog.ShowError(fmt.Errorf(Translate("Falha ao escrever config:\n%v", "Failed to write config:\n%v", currentLang), err), a.window)
			return
		}

		servicePath := "/etc/systemd/system/logid.service"
		if err := WriteServiceFile(servicePath); err != nil {
			outputEntry.SetText(fmt.Sprintf(Translate("Erro ao escrever serviço: %v", "Error writing service: %v", currentLang), err))
			dialog.ShowError(fmt.Errorf(Translate("Falha ao escrever arquivo de serviço:\n%v", "Failed to write service file:\n%v", currentLang), err), a.window)
			return
		}

		out, err := EnableAndStartService()
		if err != nil {
			outputEntry.SetText(fmt.Sprintf(Translate("%s\nErro: %v", "%s\nError: %v", currentLang), out, err))
			dialog.ShowError(fmt.Errorf(Translate("Falha ao iniciar serviço:\n%v", "Failed to start service:\n%v", currentLang), err), a.window)
			return
		}
		outputEntry.SetText(out)
		dialog.ShowInformation(Translate("Sucesso", "Success", currentLang), Translate("Serviço instalado e iniciado com sucesso", "Service installed and started successfully", currentLang), a.window)
		updateStatus()
	})

	stopServiceBtn := widget.NewButton(Translate("Parar Serviço", "Stop Service", currentLang), func() {
		log.Printf("Stop Service clicked")
		out, err := StopService()
		if err != nil {
			outputEntry.SetText(fmt.Sprintf(Translate("Erro: %v\n%s", "Error: %v\n%s", currentLang), err, out))
			dialog.ShowError(fmt.Errorf(Translate("Falha ao parar serviço:\n%v", "Failed to stop service:\n%v", currentLang), err), a.window)
			return
		}
		outputEntry.SetText(out)
		dialog.ShowInformation(Translate("Sucesso", "Success", currentLang), Translate("Serviço parado", "Service stopped", currentLang), a.window)
		updateStatus()
	})

	restartServiceBtn := widget.NewButton(Translate("Reiniciar Serviço", "Restart Service", currentLang), func() {
		log.Printf("Restart Service clicked")
		out, err := RestartService()
		if err != nil {
			outputEntry.SetText(fmt.Sprintf(Translate("Erro: %v\n%s", "Error: %v\n%s", currentLang), err, out))
			dialog.ShowError(fmt.Errorf(Translate("Falha ao reiniciar serviço:\n%v", "Failed to restart service:\n%v", currentLang), err), a.window)
			return
		}
		outputEntry.SetText(out)
		dialog.ShowInformation(Translate("Sucesso", "Success", currentLang), Translate("Serviço reiniciado", "Service restarted", currentLang), a.window)
		updateStatus()
	})

	removeServiceBtn := widget.NewButton(Translate("Remover Serviço", "Remove Service", currentLang), func() {
		log.Printf("Remove Service clicked")
		outputEntry.SetText(Translate("Removendo serviço systemd...\n", "Removing systemd service...\n", currentLang))
		out, err := RemoveService()
		if err != nil {
			outputEntry.SetText(fmt.Sprintf(Translate("%s\nErro: %v", "%s\nError: %v", currentLang), out, err))
			dialog.ShowError(fmt.Errorf(Translate("Falha ao remover serviço:\n%v", "Failed to remove service:\n%v", currentLang), err), a.window)
			return
		}
		outputEntry.SetText(out)
		dialog.ShowInformation(Translate("Sucesso", "Success", currentLang), Translate("Serviço removido com sucesso", "Service removed successfully", currentLang), a.window)
		updateStatus()
	})

	servicePreview := widget.NewMultiLineEntry()
	servicePreview.SetText(serviceContent)
	servicePreview.Disable()
	servicePreview.SetMinRowsVisible(10)

	updateStatus()

	box := container.NewVBox(
		widget.NewLabelWithStyle(Translate("Gerenciamento de Serviço", "Service Management", currentLang), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewBorder(nil, nil, widget.NewLabel(Translate("Caminho do Config:", "Config path:", currentLang)), nil, configPathEntry),
		refreshBtn,
		statusLabel,
		widget.NewSeparator(),
		widget.NewLabelWithStyle(Translate("Ações", "Actions", currentLang), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		saveConfigBtn,
		resetConfigBtn,
		installServiceBtn,
		container.NewHBox(stopServiceBtn, restartServiceBtn),
		removeServiceBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle(Translate("Unidade de Serviço Systemd", "Systemd Service Unit", currentLang), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		servicePreview,
		widget.NewSeparator(),
		widget.NewLabelWithStyle(Translate("Saída", "Output", currentLang), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
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
	if a.deviceInfoLabel != nil {
		a.deviceInfoLabel.SetText(DeviceNameInfo(currentLang))
	}
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

func (a *App) formatThumbAction(a2 Action) string {
	if a2.Type == "" || a2.Type == "None" {
		return "None"
	}
	return a.formatAction(a2)
}

func (a *App) applyConfigToUI() {
	a.nameSelect.SetSelected(a.config.Name)
	if a.nameSelect.Selected == "" {
		a.nameSelect.SetSelected(DeviceNames[0])
		a.config.Name = DeviceNames[0]
	}
	a.dpiSlider.Value = float64(a.config.DPI)
	a.dpiSlider.Refresh()
	a.dpiLabel.SetText(fmt.Sprintf("DPI: %d", a.config.DPI))

	a.smartShiftCheck.SetChecked(a.config.SmartShift.On)
	a.smartShiftSlider.Value = float64(a.config.SmartShift.Threshold)
	a.smartShiftSlider.Refresh()
	a.smartShiftLabel.SetText(fmt.Sprintf("Threshold: %d", a.config.SmartShift.Threshold))

	a.hiresCheck.SetChecked(a.config.HiResScroll.Hires)
	a.hiresInvertCheck.SetChecked(a.config.HiResScroll.Invert)
	a.hiresTargetCheck.SetChecked(a.config.HiResScroll.Target)

	a.thumbDivertCheck.SetChecked(a.config.Thumbwheel.Divert)
	a.thumbInvertCheck.SetChecked(a.config.Thumbwheel.Invert)
	a.thumbLeftLabel.SetText(a.formatThumbAction(a.config.Thumbwheel.Left))
	a.thumbRightLabel.SetText(a.formatThumbAction(a.config.Thumbwheel.Right))
	a.thumbTapLabel.SetText(a.formatThumbAction(a.config.Thumbwheel.Tap))

	for i, btn := range a.config.Buttons {
		if i < len(a.buttonWidgets) {
			a.buttonWidgets[i].actionLabel.SetText(a.formatAction(btn.Action))
		}
	}

	a.refreshPreview()
}

func (a *App) Run() {
	a.refreshPreview()
	a.window.SetContent(a.BuildUI())
	a.window.ShowAndRun()
}
