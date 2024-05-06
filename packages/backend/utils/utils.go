package utils

import (
	"fmt"
	"sync"
)

var (
	CourseChannels = make(map[string]chan []byte)

	CourseMutex sync.Mutex
)

type CourseContent struct {
	Content       string `json:"content"`
	ContentMutext sync.Mutex
}

var (
	CourseContentMap   = make(map[string]CourseContent)
	CourseContentMutex sync.Mutex
)

var (
	CourseStreamChannels = make(map[string]chan StreamResponse)
	CourseStreamMutex    sync.Mutex
)

type Course struct {
	Title          string                    `json:"title"`
	Mode           int                       `json:"mode"`
	Docs           []string                  `json:"docs"`
	Pyqs           []string                  `json:"pyqs"`
	UserId         string                    `json:"userId"`
	ProcessingData map[string]ProcessingData `json:"processingData"`
}

type ProcessingData struct {
	Status bool   `json:"status"`
	Type   string `json:"type"`
}
type StreamResponse struct {
	Message string  `json:"status"`
	Error   *string `json:"error"`
	Done    bool    `json:"done"`
}

func InputPrompt(courseJson string) string {
	return fmt.Sprintf(`
	"The following is the json format in which the input will be provided to you:"
	
	
	Input JSON Format:
	[
	  {
		"content": "Content on a page of a PPT/PDF.",
		"pageNumber": Page number
	  },
	  {
		"content": "Next content on another page.",
		"pageNumber": Page number
	  },
	  ...
	]
	
	%s
	
	Instructions for Note Generation:
	
	Detail-Oriented Notes: Break down the extracted text into detailed study notes. Include explanations, examples, and definitions to ensure comprehensive coverage of the topic. Provide real-world examples to illustrate key concepts and enhance understanding.
	
	Clarity and Simplicity: Ensure that the generated notes are clear and easy to understand. Use concise language and keep the explanations straightforward to facilitate quick comprehension.
	
	Scoring Optimization: Aim to produce study materials that can help students score well in exams. Prioritize accuracy, relevance, and completeness in the generated notes.
	
	Minimum 5 Key Points: Each set of notes should contain a minimum of 5 key points that are essential for understanding the topic thoroughly and scoring well in exams.
	
	JSON Format: Provide the notes in the following JSON format:
	[
	  {
		"topicName": "Topic Name",
		"content": "Detailed notes on the topic.",
		"fileAndPage": [
		  {
			"filename": "Name of the file",
			"pagenumber": Page number,
			"pyqs": ["Possible exam questions related to the topic."]
		  }
		]
	  }
	]

	You should only provide the json and not any other information also dont add code indecators at start and end of file" .
	
	Overall, the generated notes should be detailed, informative, and engaging enough for a student that studies from these notes shouldn't have to search for the same topic ever again.
	
	Additional Context:
	The extracted text contains key concepts, definitions, and explanations presented in a lecture. The goal is to create detailed study notes that include examples and explanations in simple language to assist students in understanding the material thoroughly and quickly, thereby improving their academic performance.
	`, courseJson)
}
