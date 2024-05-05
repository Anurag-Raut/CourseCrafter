package utils

import (
	"sync"
)

var (
	CourseChannels = make(map[string]chan []byte)

	CourseMutex sync.Mutex
)

type CourseContent struct{
	Content string `json:"content"`
	ContentMutext sync.Mutex
}

var (
	CourseContentMap = make(map[string]CourseContent)
	CourseContentMutex sync.Mutex
)

var (
	CourseStreamChannels = make(map[string]chan []byte)
	CourseStreamMutex sync.Mutex

)

type Course struct {
	Title          string          `json:"title"`
	Mode           int             `json:"mode"`
	Docs           []string        `json:"docs"`
	Pyqs           []string        `json:"pyqs"`
	UserId         string          `json:"userId"`
	ProcessingData map[string]ProcessingData `json:"processingData"`
}

type ProcessingData struct {
	Status bool   `json:"status"`
	Type   string `json:"type"`
}


