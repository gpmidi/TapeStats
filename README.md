# TapeStats
Public Tape Stats For LTO

## About
TapeStats.com is (will be) a site dedicated to providing information about common LTO media. 
This includes failure rates, lifetime expectancy, failure indicators, and anything else
we can pull from the data. 

# Where We Are
Once the submission code is complete and is live on www.tapestats.com work will begin on
adding stats code. That'll provide basic insight into the data that is collected. 

Note: No stats will be published until there are at least 500 tapes in the database.   

## Ideas/TODO
* Need to finish the basics first!
* Add data stats (coming later)
    * Use a RO Postgres follower for stats queries
    * Redis for caching
* Better auth system / options
* Rate limiting to RegisterAccountHandler - Maybe one per min per source IP
* Fix Created/Modified on tables
* Swagger file
* Add dict input for user defined data - maybe for later use/standardization?
* Allow accounts to delete themselves and their data
* Rewrite [mamtool](https://github.com/redrice/mamtool) in golang
* Allow accounts to get data about their tapes/submissions
* Allow account password changes
* Password resets for accounts?
* Allow input of purchase cost data for tape

## Dev Tips

### Using Cobra Commander
`cobra add --config .cobra.yaml -p rootCmd somethingnew`
