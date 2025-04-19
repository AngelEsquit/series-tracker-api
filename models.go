package main

// Series struct representa una serie de TV
type Series struct {
	// ID es el identificador único de la serie
	ID int `json:"id"`
	// Title es el título de la serie
	Title string `json:"title"`
	// Status es el estado de la serie (e.g., "Plan to Watch", "Watching", "Dropped", "Completed")
	Status string `json:"status"`
	// LastEpisodeWatched es el último episodio visto por el usuario
	LastEpisodeWatched int `json:"lastEpisodeWatched"`
	// TotalEpisodes es el total de episodios de la serie
	TotalEpisodes int `json:"totalEpisodes"`
	// Ranking es la calificación de la serie
	Ranking int `json:"ranking"`
}

// SeriesList es una lista de Series
type ErrorResponse struct {
	// Error es el mensaje de error
	Message string `json:"message"`
}

type SuccessResponse struct {
	// Message es el mensaje de éxito
	Message string `json:"message"`
}
