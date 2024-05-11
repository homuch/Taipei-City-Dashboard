// Package controllers stores all the controllers for the Gin router.
package controllers

import (
	"net/http"
	"strconv"

	"TaipeiCityDashboardBE/app/models"
	"TaipeiCityDashboardBE/app/util"

	"github.com/gin-gonic/gin"

	"fmt"
	"strings"
)

/*
GetComponentChartData retrieves the chart data for a component.
/api/v1/components/:id/chart

header: time_from, time_to (optional)
*/
func GetComponentChartData(c *gin.Context) {
	// 1. Get the component id from the URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid component ID"})
		return
	}
	var queryString string

	// 2. Get the chart data query and chart data type from the database
	queryType, queryStringOrigined, queryStringFiltered, err := models.GetComponentChartDataQuery(id)
	// fmt.Println(queryStringOrigined)
	// fmt.Println("-----------------")
	// fmt.Println(queryStringFiltered)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if (queryStringOrigined == "") || (queryType == "") {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No chart data available"})
		return
	}
	

	filter_by_distance:= c.Query("filter_by_distance")
	// fmt.Println(filter_by_distance)
	if filter_by_distance == "true" &&
	(strings.Count(queryStringFiltered, "%s") == 3 ||
	strings.Count(queryStringOrigined, "%s") == 5){
		// if strings.Count(queryStringFiltered, "%s") != 3 ||
		// 	strings.Count(queryStringOrigined, "%s") != 5 {
		// 	queryString = queryStringOrigined
		// }
		filter_distance := c.Query("filter_distance")
		filter_lat := c.Query("filter_lat")
		filter_long := c.Query("filter_long")
		if strings.Count(queryStringOrigined, "%s") == 5 {
			queryString = fmt.Sprintf(queryStringFiltered, "%s", "%s", filter_distance, filter_lat, filter_long)
		} else {
			queryString = fmt.Sprintf(queryStringFiltered, filter_distance, filter_lat, filter_long)
		}
	} else {
		queryString = queryStringOrigined
	}
	// fmt.Println("++++++++++++++++++++++++++++")
	// fmt.Println(queryString)
	timeFrom, timeTo := util.GetTime(c)

	// 3. Get and parse the chart data based on chart data type
	if queryType == "two_d" {
		chartData, err := models.GetTwoDimensionalData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	} else if queryType == "three_d" || queryType == "percent" {
		chartData, categories, err := models.GetThreeDimensionalData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData, "categories": categories})
	} else if queryType == "time" {
		chartData, err := models.GetTimeSeriesData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	} else if queryType == "map_legend" {
		chartData, err := models.GetMapLegendData(&queryString, timeFrom, timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
	}
}

/*
GetComponentHistoryData retrieves the history data for a component.
/api/v1/components/:id/history

header: time_from, time_to (mandatory)
timesteps are automatically determined based on the time range:
  - Within 24hrs: hour
  - Within 1 month: day
  - Within 3 months: week
  - Within 2 years: month
  - More than 2 years: year
*/
func GetComponentHistoryData(c *gin.Context) {
	// 1. Get the component id from the URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid component ID"})
		return
	}

	timeFrom, timeTo := util.GetTime(c)

	// 2. Get the history data query from the database
	queryHistory, err := models.GetComponentHistoryDataQuery(id, timeFrom, timeTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if queryHistory == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No history data available"})
		return
	}

	// 3. Get and parse the history data
	chartData, err := models.GetTimeSeriesData(&queryHistory, timeFrom, timeTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": chartData})
}
