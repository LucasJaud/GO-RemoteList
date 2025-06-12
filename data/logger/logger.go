package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type LoggerData struct{
	Service string
	Method string
	Params map[string]interface{}
}

type Logger struct {
	filename string
	mu sync.Mutex
}

func NewLogger(filename string) *Logger {
	return &Logger{
        filename: filename,
    }
}

func (l *Logger) Save(service string, method string, params map[string]interface{}) error{
	l.mu.Lock()
	defer l.mu.Unlock()

	data := LoggerData{
		Service: service,
		Method: method,
		Params: params,
	}

	file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil{
		return fmt.Errorf(" Error opening file: %v", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf(" Error parsing file: %v", err)
	}

	_, err = file.Write(append(jsonData,'\n'))
	if err != nil {
		return fmt.Errorf(" Error writing file: %v", err)
	}

	return file.Sync()
}

func (l *Logger) Load() ([]LoggerData,error){
	file, err := os.Open(l.filename)
	if err != nil {
		return nil, fmt.Errorf(" Error opening file: %v", err)
	}
	defer file.Close()

	var entries []LoggerData
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan(){
		var entry LoggerData 
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			fmt.Printf("Error converting line: %v\n", err)
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf(" Error reading file: %v", err)
    }
	return entries, nil
}

func (l *Logger) Clear() error {
    l.mu.Lock()
    defer l.mu.Unlock()
	
    return os.Truncate(l.filename, 0)
}

func (l *Logger) GetFilename() string {
    return l.filename
}
