# orbcomm-api-integration-go

Go application to ingest the Orbcomm API.

DO NOT USE: Test using https://bolt.new application.


I would like to create a new project with a new database.

technology: go, postgres, s3, csv, echo
project: We want to integrate with the Orbcomm API. The application will do a number things 1 - Receive a list of ships to manage from a csv file stored in s3. 3 - Manage a list of ships in a list on the orbcomm platform using the Vessel List API 4 - Make a call every 5 minutes to the orbcomm API to get the latest data for each ship. 5 - Store the data in a postgres database of each ship 6 - Export a new CSV file in s3 with the latest data for each ship 7 - Create an API endpoint to get the latest data for a given ship, list of ships, or all ships
reference: https://globalais3.orbcomm.net/api/v3_0/
testing: i wish to have a director of unit and integrations tests for the code.