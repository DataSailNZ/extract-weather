package snowflake

import (
	"database/sql"
	"fmt"
	"github.com/datasail/extract-weather/weather"
	"github.com/google/uuid"
	_ "github.com/snowflakedb/gosnowflake"
	"log"
	"os"
	"time"
)

func LoadInSnowflake(responses *[]weather.Response) {

	user := os.Getenv("SNOWFLAKE_USER")
	password := os.Getenv("SNOWFLAKE_PWD")
	account := os.Getenv("SNOWFLAKE_ACCOUNT")
	db := os.Getenv("SNOWFLAKE_DB")
	schema := os.Getenv("SNOWFLAKE_SCHEMA")

	dataSourceName := user + ":" + password + "@" + account + "/" + db + "/" + schema

	dbConn, err := sql.Open("snowflake", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open db: %s", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	stmt := `
create table if not exists LOAD_CITY (
    ID VARCHAR(64),
    NAME VARCHAR(255) NOT NULL,
    LAT NUMBER(6,3),
    LONG NUMBER(6,3),
    PRIMARY KEY (ID)
);
`
	_, err = dbConn.Exec(stmt)
	if err != nil {
		panic(err)
	}

	stmt = `
create table if not exists LOAD_WEATHER (
    ID VARCHAR(64),
    TIMEPOINT INT,
    CLOUDCOVER INT,
    LIFTED_INDEX INT,
    PREC_TYPE VARCHAR(16),
    PREC_AMOUNT INT,
    TEMP_2M INT,
    RH_2M VARCHAR(16),
    WEATHER VARCHAR(16),
    CITY_ID VARCHAR(64) REFERENCES LOAD_CITY(ID),
    EXTRACT_DATE TIMESTAMP,
    PRIMARY KEY (ID)
);
`
	_, err = dbConn.Exec(stmt)
	if err != nil {
		panic(err)
	}

	stmt = `
create table if not exists LOAD_WIND_10M (
    LOAD_WEATHER_ID VARCHAR(64),
    DIRECTION VARCHAR(16),
    SPEED INT,
    EXTRACT_DATE TIMESTAMP,
    FOREIGN KEY (LOAD_WEATHER_ID) REFERENCES LOAD_WEATHER(ID)
);
`
	_, err = dbConn.Exec(stmt)
	if err != nil {
		panic(err)
	}

	stmt = `
delete from LOAD_CITY;
`
	_, err = dbConn.Exec(stmt)
	if err != nil {
		panic(err)
	}

	cityId := "e6ca8b9b-b9a8-4eb1-83b2-d61431c98c7d"
	stmt = `
insert into LOAD_CITY (ID, NAME, LAT, LONG) VALUES ('%v', 'Auckland', -36.843, 174.766);
`
	_, err = dbConn.Exec(fmt.Sprintf(stmt, cityId))
	if err != nil {
		panic(err)
	}

	now := time.Now()
	timestamp := now.Format("2006-01-02 15:04:05")

	for _, response := range *responses {

		for _, dataSeries := range response.DataSeries {

			weatherId := uuid.New().String()

			stmt = `
insert into LOAD_WEATHER 
    (
     	ID, TIMEPOINT, CLOUDCOVER, LIFTED_INDEX, PREC_TYPE, PREC_AMOUNT, TEMP_2M, RH_2M, WEATHER, CITY_ID, EXTRACT_DATE
    ) 
VALUES
    (
            '%v', %v, %v, %v, '%v', %v, %v, '%v', '%v', '%v', '%v'
    );
`
			_, err = dbConn.Exec(fmt.Sprintf(stmt,
				weatherId,
				dataSeries.Timepoint,
				dataSeries.Cloudcover,
				dataSeries.LiftedIndex,
				dataSeries.PrecType,
				dataSeries.PrecAmount,
				dataSeries.Temp2m,
				dataSeries.Rh2m,
				dataSeries.Weather,
				cityId,
				timestamp))
			if err != nil {
				panic(err)
			}

			stmt = `
insert into LOAD_WIND_10M 
    (
     	LOAD_WEATHER_ID, DIRECTION, SPEED, EXTRACT_DATE
    ) 
VALUES
    (
            '%v', '%v', %v, '%v'
    );
`
			_, err = dbConn.Exec(fmt.Sprintf(
				stmt,
				weatherId,
				dataSeries.Wind10m.Direction,
				dataSeries.Wind10m.Speed,
				timestamp))
			if err != nil {
				panic(err)
			}

		}
	}
}
