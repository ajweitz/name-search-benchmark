# name-search-benchmark

## Datasets
There are 3 files in the `/data` directory
- names.txt - 18k rows of names
- words.txt - 460k rows of words
- strings.txt - 1m rows names, words & random strings

## Running Backend
### Run MySQL & Redis
1. Navigate to: `backend/infra/`
2. Run `./generate-data.sh strings` - for the 1 million items (you can speicify `words` or `names` if you want smaller dataset)
3. Run Docker - `docker-compose up`

### Running Go Application
1. Navigate to `backend/app/`
2. Run `make`

## Running Frontend
1. Navigate to `frontend/`
2. Open `index.html` in your browser

## Running Tests
1. Navigate to `test/`
2. For testing a particular case: run `locust -f <filename>`.
3. Access http://localhost:8089/ 
4. Specify `http://localhost:8080` as the Host.