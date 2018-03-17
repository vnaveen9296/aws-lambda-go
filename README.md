# aws-lambda-go

AWS Lambda now supports Go. This provides details about how to set up users, roles, policies, create lambda function etc along with a very simple Go program that can be deployed as a Lambda function.

## Details
The following steps assume there is already an administrator user within your AWS account (called `admin`) and it's credentials are saved on your machine using `aws configure --profile admin` command

### Create a group
Create a group called `lambda-group` within your AWS admin account.

`aws iam create-group --group-name lambda-group --profile admin`

### Attach an AWS managed policy to the group
Attach an AWS managed policy to `lambda-group`. To enable users of this group with full access to Lambda, `AWSLambdaFullAccess` policy is used below

`aws iam attach-group-policy --group-name lambda-group --policy-arn arn:aws:iam::aws:policy/AWSLambdaFullAccess  --profile admin`

### Create a user
Create a user called `lambda-user` that works with Lambda functions

`aws iam create-user --user-name lambda-user --profile admin`

### Add the user to the group
Add `lambda-user` to `lambda-group`

`aws iam add-user-to-group --user-name lambda-user --group-name lambda-group --profile admin`

### Access key for the user
Create an access key for the user `lambda-user` (save the access key id and secret key safely)

`aws iam create-access-key --user-name lambda-user --profile admin`

### Grant certain IAM related permissions to the user
Add an inline policy to grant certain IAM related permissions to `lambda-user` (for example, PassRole permission is required when creating a Lambda function)

`aws iam put-user-policy --user-name lambda-user --policy-name iam-permissions-for-lambda --policy-document file://iam_permissions_for_lambda_user.json --profile admin`

### Create a role that will be assumed by the Lambda function
Create a role called `lambda-exec-role`. This is the role that will be associated with each Lambda function we create. In short, this role determines what AWS Lambda function can do when it executes (i.e can it access S3, can it access CloudWatch logs etc.)

`aws iam create-role --role-name lambda-exec-role --assume-role-policy file://D:\aws_work\trust_policy_for_lambda.json --profile admin`

### Attach a managed policy to the role
Attach an AWS managed permission policy to the role `lambda-exec-role`. In this case, `AWSLambdaExecute` policy is used. This grants permissions to CloudWatch Logs actions as well S3 get and put actions

`aws iam attach-role-policy --role-name lambda-exec-role --policy-arn arn:aws:iam::aws:policy/AWSLambdaExecute --profile admin`

### Save user credentials in local credentials file on your machine
Run `aws configure` to set credentials for `lambda-user` in the AWS credentials file on your local machine. Enter access key id and secret key when asked.

`aws configure --profile lambda-user`

A user called `lambda-user` is now set up. This user can now work with Lambda functions (create, update, add triggers etc.)

## Create and Invoke a Lambda function
The following steps can be used to create a new Lambda function and invoke it by taking an action such as "uploading an object to S3 bucket"

## Create a Lambda function
We create a function named as `lambda-hello` and specify the role that was created above (`lambda-exec-role`). The zip file contains our Go programs binary.

`aws lambda create-function --function-name lambda-hello --runtime go1.x --role arn:aws:iam::721266161153:role/lambda-exec-role --handler main --zip-file fileb://lambda-hello.zip --profile lambda-user`

### Grant permissions to S3
Grant permissions to an event source such as S3 to push events (or invoke this Lambda function)

`aws lambda add-permission --function-name arn:aws:lambda:us-west-2:721266161153:function:lambda-hello --statement-id 1 --action lambda:* --principal s3.amazonaws.com --source-arn arn:aws:s3:::vn-utts --profile lambda-user`

### Create a bucket in S3
Create a bucket named `vn-upload-for-txt`. We configure this as to push events and invoke the Lambda function whenever a new object with a suffix of `.txt` is created

`aws s3api create-bucket --bucket vn-upload-for-txt --region us-west-2 --create-bucket-configuration LocationConstraint=us-west-2 --profile lambda-user`

### Bucket notification configuration
S3 bucket you created above needs to be configured to push events whenever a new object of type `.txt` is created. We choose the Lambda function as the destination for these events (i.e. the Lambda function will be invoked in response to these events)

`aws s3api put-bucket-notification-configuration --bucket vn-utts --notification-configuration file://s3_configuration_to_trigger_lambda.json --profile lambda-user`

### Upload a file to S3
Upload a .txt file to S3 bucket `vn-upload-for-txt`

`aws s3api put-object --bucket vn-upload-for-txt --key utts1.txt --body ..\utts\utts1.txt --profile lambda-user`

Now, goto AWS console on your browser, choose Lambda from services and select the Lambda function `lambda-hello` to open its configuration panel. You can click on Monitoring to see some graphs. There is also a link to CloudWatch logs on this page. Click on this to see the log messages that confirm the invocation of your Lambda function
