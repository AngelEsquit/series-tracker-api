package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @summary Crear una nueva serie
// @description Crea una nueva serie en la base de datos
// @tags series
// @produce json
// @router /api/series [post]
// @param series body Series true "Datos de la serie"
// @success 200 {string} string "Serie creada"
func createSeries(w http.ResponseWriter, r *http.Request) {
	var series Series
	err := json.NewDecoder(r.Body).Decode(&series)

	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	err = InsertSeries(series)
	if err != nil {
		http.Error(w, "Error al insertar serie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Serie '%s' creada", series.Title),
	})
}

// @summary Obtener todas las series
// @description Obtiene una lista de todas las series en la base de datos
// @tags series
// @produce json
// @router /api/series [get]
// @param search query string false "Buscar por título"
// @param status query string false "Filtrar por estado (e.g., 'Watching', 'Completed')" enums(Watching, Completed, Dropped, PlanToWatch)
// @param sort query string false "Ordenar por (e.g., 'asc', 'desc')" enums(asc, desc)
// @success 200 {array} Series "Lista de series"
func getAllSeries(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")

	seriesList, err := GetSeriesWithFilters(search, status, sort)
	if err != nil {
		http.Error(w, "Error al obtener la lista de series", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seriesList)
}

// @summary Obtener una serie por ID
// @description Obtiene una serie específica por su ID
// @tags series
// @produce json
// @router /api/series/{id} [get]
// @param id path int true "ID de la serie"
// @success 200 {object} Series "Serie encontrada"
// @failure 404 {object} ErrorResponse "Serie no encontrada"
func getSeriesByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")

	seriesList, err := GetSeriesWithFilters(search, status, sort)
	if err != nil {
		http.Error(w, "Error al obtener la serie", http.StatusInternalServerError)
		return
	}

	for _, series := range seriesList {
		if fmt.Sprintf("%d", series.ID) == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(series)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: fmt.Sprintf("Serie con ID '%s' no encontrada", id),
	})
}

// @summary Actualizar una serie por ID
// @description Actualiza una serie específica por su ID
// @tags series
// @produce json
// @router /api/series/{id} [put]
// @param id path int true "ID de la serie"
// @param series body Series true "Datos de la serie"
// @success 200 {object} SuccessResponse "Serie actualizada"
// failure 404 {object} ErrorResponse "Error al actualizar la serie"
func updateSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "ID inválido",
		})
		return
	}

	var serie Series
	err = json.NewDecoder(r.Body).Decode(&serie)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Error al decodificar JSON",
		})
		return
	}

	serie.ID = id

	err = UpdateSeries(serie)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al actualizar la serie con ID '%d': %v", id, err),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: fmt.Sprintf("Serie con ID '%d' actualizada", id),
	})
}

// @summary Eliminar una serie por ID
// @description Elimina una serie específica por su ID
// @tags series
// @produce json
// @router /api/series/{id} [delete]
// @param id path int true "ID de la serie"
// @success 200 {object} SuccessResponse "Serie eliminada"
// @failure 404 {object} ErrorResponse "Error al eliminar la serie"
func deleteSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "ID inválido",
		})
		return
	}

	err = DeleteSeries(id)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al eliminar la serie con ID '%d': %v", id, err),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: fmt.Sprintf("Serie con ID '%d' eliminada", id),
	})
}

// @summary Actualizar el estado de una serie por ID
// @description Actualiza el estado de una serie específica por su ID
// @tags series
// @produce json
// @router /api/series/{id}/status [patch]
// @param id path int true "ID de la serie"
// @param status query string true "Nuevo estado de la serie (e.g., 'Watching', 'Completed')" enums(Watching, Completed, Dropped, PlanToWatch)
// @success 200 {object} SuccessResponse "Estado actualizado"
// @failure 404 {object} ErrorResponse "Serie no encontrada"
func updateStatus(w http.ResponseWriter, r *http.Request) {
	// Manejo del preflight
	params := mux.Vars(r)
	idStr := params["id"]

	// Convertir el id de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "ID inválido",
		})
		return
	}

	var serie Series

	err = json.NewDecoder(r.Body).Decode(&serie)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Error al decodificar JSON",
		})
		return
	}

	// Llamar a la función UpdateStatus
	err = UpdateStatus(id, serie.Status)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al actualizar el estado de la serie con ID '%d': %v", id, err),
		})
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: fmt.Sprintf("Estado de la serie con ID '%d' actualizado a '%s'", id, serie.Status),
	})
}

// @summary Actualizar el episodio de una serie por ID
// @description Actualiza el episodio de una serie específica por su ID
// @tags series
// @produce json
// @router /api/series/{id}/episode [patch]
// @param id path int true "ID de la serie"
// @success 200 {object} SuccessResponse "Episodio actualizado"
// @failure 404 {object} ErrorResponse "Serie no encontrada"
func updateEpisode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "ID inválido",
		})
		return
	}

	var serie Series

	serie, err = GetSeriesByID(id)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al obtener la serie con ID '%d': %v", id, err),
		})
		return
	}

	// Llamar a la función UpdateEpisode
	err = UpdateEpisode(id, serie.LastEpisodeWatched+1)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al actualizar el episodio con ID '%d': %v", id, err),
		})
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: fmt.Sprintf("Episodio de la serie con ID '%d' actualizado a '%d'", id, (serie.LastEpisodeWatched)+1),
	})
}

// @summary Aumentar el ranking de una serie por ID
// @description Aumenta el ranking de una serie específica por su ID
// @tags series
// @produce json
// @router /api/series/{id}/upvote [patch]
// @param id path int true "ID de la serie"
// @success 200 {object} SuccessResponse "Ranking actualizado"
// @failure 404 {object} ErrorResponse "Serie no encontrada"
func upVote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "ID inválido",
		})
		return
	}

	var serie Series

	serie, err = GetSeriesByID(id)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al obtener la serie con ID '%d': %v", id, err),
		})
		return
	}

	// Llamar a la función UpdateEpisode
	err = UpVote(id)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al actualizar el episodio con ID '%d': %v", id, err),
		})
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: fmt.Sprintf("Episodio de la serie con ID '%d' actualizado a '%d'", id, (serie.Ranking)-1),
	})
}

// @summary Disminuir el ranking de una serie por ID
// @description Disminuye el ranking de una serie específica por su ID
// @tags series
// @produce json
// @router /api/series/{id}/downvote [patch]
// @param id path int true "ID de la serie"
// @success 200 {object} SuccessResponse "Ranking actualizado"
// @failure 404 {object} ErrorResponse "Serie no encontrada"
func downVote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "ID inválido",
		})
		return
	}

	var serie Series

	serie, err = GetSeriesByID(id)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al obtener la serie con ID '%d': %v", id, err),
		})
		return
	}

	// Llamar a la función UpdateEpisode
	err = DownVote(id)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: fmt.Sprintf("Error al actualizar el episodio con ID '%d': %v", id, err),
		})
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: fmt.Sprintf("Episodio de la serie con ID '%d' actualizado a '%d'", id, (serie.Ranking)+1),
	})
}
