.PHONY: build

build:
	sam build

deploy-infra:
	sam build && aws-vault exec alvis --no-session -- sam deploy

deploy-site:
	aws-vault exec alvis --no-session -- aws s3 sync ./resume-site/_site s3://my-starter-website

invoke-put:
	sam build && aws-vault exec alvis --no-session -- sam local invoke PutFunction

invoke-get:
	sam build && aws-vault exec alvis --no-session -- sam local invoke GetFunction
