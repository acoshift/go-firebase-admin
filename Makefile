SERVICE_ACCOUNT=$(shell cat private/service_account.json)
PROJECT_ID=$(shell cat private/project_id)
DATABASE_URL=$(shell cat private/database_url)
API_KEY=$(shell cat private/api_key)

test:
	env \
		SERVICE_ACCOUNT='$(SERVICE_ACCOUNT)' \
		PROJECT_ID='$(PROJECT_ID)' \
		DATABASE_URL='$(DATABASE_URL)' \
		API_KEY='$(API_KEY)' \
	go test .
