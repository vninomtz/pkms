package model

type Note struct {
  Name string
  Content string
  Meta map[interface{}]interface{}
  Parent string
  Size int64
}
