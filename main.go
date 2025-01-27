package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/patient/:id", getPatientData)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Servidor corriendo en :%s", port)
	r.Run(":" + port)
}

func getPatientData(c *gin.Context) {
	patientID := c.Param("id")
	query := "SELECT * FROM medical_records WHERE id_patient = $1" 

	queryRequest := map[string]interface{}{
		"sql":  query,
		"args": []interface{}{patientID},
	}

	queryBody, err := json.Marshal(queryRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar la consulta"})
		return
	}

	queryServiceURL := "http://localhost:8001/query"
	resp, err := http.Post(queryServiceURL, "application/json", bytes.NewBuffer(queryBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al conectar con el servicio de consulta"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Error al obtener los datos del paciente"})
		return
	}

	var queryResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&queryResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar los datos del paciente"})
		return
	}

	c.JSON(http.StatusOK, queryResponse)
}
