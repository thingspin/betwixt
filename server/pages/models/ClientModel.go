package models

type ClientModel struct {
    Endpoint         string
    RegistrationID   string
    RegistrationDate string
    LastUpdate       string
    Objects          map[string]ObjectModel
}