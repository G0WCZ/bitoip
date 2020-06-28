/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"context"
	"fmt"
	"github.com/G0WCZ/cwc/bitoip"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func duration(t time.Time) string {
	if t.Equal(time.Time{}) {
		return "-"
	}

	duration := time.Now().Sub(t)

	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"d", days},
		{"h", hours},
		{"m", minutes},
		{"s", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d%s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d%s", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, "")
}

func renderer() multitemplate.Renderer {
	funcMap := template.FuncMap{
		"timefmt":  func(t time.Time, f string) string { return t.Format(f) },
		"duration": duration,
	}
	r := multitemplate.NewRenderer()
	r.AddFromFilesFuncs("index", funcMap, "web/tmpl/base.html", "web/tmpl/index.html")
	r.AddFromFilesFuncs("channel", funcMap, "web/tmpl/base.html", "web/tmpl/channel.html")
	return r
}

func APIServer(ctx context.Context, channels *ChannelMap, config *ReflectorConfig) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Static("/static", "./web/root")
	router.HTMLRender = renderer()

	router.GET("/api/channels", func(c *gin.Context) {
		c.JSON(http.StatusOK, *channels)
	})
	router.GET("/api/activity", func(c *gin.Context) {
		c.JSON(http.StatusOK, ChannelActivity)
	})
	router.GET("/channels/:cid", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("cid"))
		channel := *GetChannel(bitoip.ChannelIdType(id))

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.HTML(200, "channel", gin.H{
				"ServerName":              config.ReflectorName,
				"Channel":                 channel,
				"ListenSortedSubscribers": channel.GetListenSortedSubscribers(),
			})
		}
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"ServerName": config.ReflectorName,
			"Channels":   channels,
			"Activity":   ChannelActivity,
		})
	})

	router.Run(config.WebAddress)
}
