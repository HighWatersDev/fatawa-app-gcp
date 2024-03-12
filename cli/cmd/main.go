package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const baseURL = "http://localhost:8080"

// Models for BubbleTea
type model struct {
	choices    []string          // API actions
	cursor     int               // Menu cursor
	inputs     []textinput.Model // Input fields
	activeView string            // Current view state: "menu" or "inputs"
	response   string            // API response message
	err        error             // Error state
}

type errMsg error
type responseMsg string

// Initialize the model
func initialModel() model {
	m := model{
		choices:    []string{"Call Endpoint 1 (GET)", "Call Endpoint 2 (POST)"},
		activeView: "menu",
	}

	// Initialize text inputs for endpoint2
	fileInput := textinput.New()
	fileInput.Placeholder = "File Path"
	fileInput.Focus()

	authorInput := textinput.New()
	authorInput.Placeholder = "Author"

	m.inputs = []textinput.Model{fileInput, authorInput}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

// Model Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlQ:
			return m, tea.Quit
		}

		switch m.activeView {
		case "menu":
			return m.updateMenu(msg)
		case "inputs":
			return m.updateInputs(msg)
		case "response":
			m.activeView = "menu"
			return m, nil
		}

	case responseMsg:
		m.response = string(msg)
		m.activeView = "response" // Switch to response view after API call
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		switch m.cursor {
		case 0:
			// Example GET request
			endpoint := "/v1/processor"
			return m, callAPI("GET", endpoint, nil, nil)
		case 1:
			// Transition to inputs for POST request
			m.activeView = "inputs"
			m.inputs[0].Focus()
			return m, nil
		}
	case tea.KeyUp, tea.KeyCtrlK:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown, tea.KeyCtrlJ:
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	default:
		return m, nil
	}

	return m, nil
}

func (m model) updateInputs(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEnter:
		filePath := m.inputs[0].Value()
		author := m.inputs[1].Value()

		// Assuming /v1/processor/storage/upload endpoint expects query parameters and no body.
		endpoint := "/v1/processor"
		queryParams := map[string]string{"path": filePath, "author": author}

		// Adjusted for new callAPI signature, no body data for this case
		return m, callAPI("POST", endpoint, queryParams, nil)

	case tea.KeyTab, tea.KeyShiftTab:
		// Cycling focus through inputs
		for i := range m.inputs {
			if m.inputs[i].Focused() {
				m.inputs[i].Blur()
				nextIndex := (i + 1) % len(m.inputs)
				m.inputs[nextIndex].Focus()
				break
			}
		}
	default:
		return m, nil
	}

	// Update inputs based on key press
	for i := range m.inputs {
		m.inputs[i], cmd = m.inputs[i].Update(msg)
	}
	return m, cmd
}

// Model View function
func (m model) View() string {
	switch m.activeView {
	case "menu":
		return m.menuView()
	case "inputs":
		return m.inputsView()
	case "response":
		return m.responseView()
	}
	return ""
}

func (m model) menuView() string {
	s := "Choose an API call:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}

func (m model) inputsView() string {
	s := "Enter details:\n\n"
	for _, input := range m.inputs {
		s += input.View()
		s += "\n"
	}
	s += "\nPress Enter to submit, Ctrl+C or Esc to quit."
	return s
}

func (m model) responseView() string {
	doc := lipgloss.NewStyle().Width(50).Render(m.response) // Set the desired width
	return fmt.Sprintf("Response:\n%s\n\nPress any key to return.", doc)
}

// Command to call the API
func callAPI(method, endpoint string, queryParams map[string]string, data interface{}) tea.Cmd {
	return func() tea.Msg {
		// Prepare the URL
		u, err := url.Parse(endpoint)
		if err != nil {
			return errMsg(fmt.Errorf("error parsing URL: %w", err))
		}

		// Add query parameters if present
		if queryParams != nil && len(queryParams) > 0 {
			q := u.Query()
			for key, value := range queryParams {
				q.Add(key, value)
			}
			u.RawQuery = q.Encode()
		}

		// Encode the data to JSON if it's a POST or PUT request with a body
		var buf bytes.Buffer
		if data != nil && (method == "POST" || method == "PUT") {
			err := json.NewEncoder(&buf).Encode(data)
			if err != nil {
				return errMsg(fmt.Errorf("error encoding request body: %w", err))
			}
		}

		// Create the request
		req, err := http.NewRequest(method, u.String(), &buf)
		if err != nil {
			return errMsg(fmt.Errorf("error creating request: %w", err))
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		authToken := os.Getenv("AUTH_TOKEN")
		if authToken != "" {
			req.Header.Set("Authorization", "Bearer "+authToken)
		}

		// Execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return errMsg(fmt.Errorf("error making request: %w", err))
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		// Read and return the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errMsg(fmt.Errorf("error reading response body: %w", err))
		}

		return responseMsg(string(body))
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if err, _ := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err) // This logs an error and exits with code 1
	}
}
