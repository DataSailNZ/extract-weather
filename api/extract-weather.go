package main

import (
    "fmt"
    "github.com/datasail/extract-weather/snowflake"
    "github.com/datasail/extract-weather/weather"
)

func main() {
    responses := weather.ExtractWeather()
    fmt.Println("Weather extraction complete")

    snowflake.LoadInSnowflake(responses)
    fmt.Println("Load complete")
}
