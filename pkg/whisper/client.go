package whisper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Client struct {
	whisperPath string
}

type TranscriptionOptions struct {
	Language     string
	OutputFormat string
	Model        string
	Temperature  float64
	Threads      int
}

func NewClient(whisperPath string) *Client {
	return &Client{
		whisperPath: whisperPath,
	}
}

func (c *Client) Transcribe(audioFile io.Reader, options TranscriptionOptions) (string, error) {
	tempFile, err := os.CreateTemp("", "whisper-audio-*.wav")
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // TODO: написать хендерл для ошибки
	defer tempFile.Close()           // TODO: написать хендерл для ошибки

	log.Printf("Created temp file: %s", tempFile.Name())

	_, err = io.Copy(tempFile, audioFile)
	if err != nil {
		log.Printf("Error writing audio to temp file: %v", err)
		return "", fmt.Errorf("failed to write audio to temp file: %w", err)
	}

	log.Printf("Audio written to temp file: %s", tempFile.Name())

	// TODO: сделать нормальную директорию для сохранения
	outputDir := "/root"

	cmdArgs := []string{tempFile.Name(), "--verbose", "False", "--output_dir", outputDir}

	optionsMap := map[string]string{
		"--language":      options.Language,
		"--output_format": options.OutputFormat,
		"--model":         options.Model,
	}

	optionsFloatMap := map[string]float64{
		"--temperature": options.Temperature,
	}

	optionsIntMap := map[string]int{
		"--threads": options.Threads,
	}

	for flag, value := range optionsMap {
		if value != "" {
			cmdArgs = append(cmdArgs, flag, value)
		}
	}

	for flag, value := range optionsFloatMap {
		if value != 0 {
			cmdArgs = append(cmdArgs, flag, fmt.Sprintf("%f", value))
		}
	}

	for flag, value := range optionsIntMap {
		if value != 0 {
			cmdArgs = append(cmdArgs, flag, fmt.Sprintf("%d", value))
		}
	}

	log.Printf("Executing whisper command with args: %v", cmdArgs)
	cmd := exec.Command(c.whisperPath, cmdArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running whisper command: %v", err)
		log.Printf("Whisper output: %s", out.String())
		return "", fmt.Errorf("whisper command failed: %w", err)
	}

	log.Printf("Whisper command output: %s", out.String())

	// TODO: сохранять не только те форматы что указаны в запросе, а все, выдавать только указанный формат
	outputFilePrefix := strings.TrimSuffix(filepath.Base(tempFile.Name()), ".wav")
	outputFile := filepath.Join(outputDir, outputFilePrefix+".txt")

	log.Printf("Looking for output file: %s", outputFile)
	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		log.Printf("Error reading output file: %v", err)
		return "", fmt.Errorf("failed to read output file: %w", err)
	}

	return string(content), nil
}
