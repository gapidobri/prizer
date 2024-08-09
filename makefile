update-db:
	atlas schema apply --env local --auto-approve

generate-spec:
	SWAGGER_GENERATE_EXTENSION=false swagger generate spec --exclude-tag=admin -m --exclude-deps -o ./api/public.yaml
	SWAGGER_GENERATE_EXTENSION=false swagger generate spec --include-tag=admin -m --exclude-deps -o ./api/admin.yaml