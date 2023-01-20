package model

type Note struct {
  Name string `json:"name"`
  Content string `json:"content"`
  Meta map[string]interface{} `json:"metadata"`
  Parent string `json:"folder"`
  Size int64 `json:"size"`
}
