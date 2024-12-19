Rest API
===============


Lambda Integration
---------------------------------------------

Proxy: (Easy)

-  the input to the integrated Lambda function can be expressed as any combination of request headers, path variables, query string parameters, and body.
-  Lambda function can use API configuration settings to influence its execution logic.
-  API Gateway configures the integration request and integration response for you
-  the integrated API method can evolve with the backend without modifying the existing settings.
-  the backend Lambda function developer parses the incoming request data and responds 
   -  with desired results to the client when nothing goes wrong or 
   -  with error messages when anything goes wrong.

* API Gateway accept any request header and parameters, transform it to fix data format and pass it to lambda function.
  * Fixed input data format to Lambda function
* Lambda function verify the API request and handle all responds include error code.

Non-Proxy:

- Dev ensure that input to the Lambda function is supplied as the integration request payload
- Dev map any input data the client supplied as request parameters into the proper integration request body to Lambda function
- Dev need to translate the client-supplied request body into a format recognized by the Lambda function.

* API Gateway verify the API request and handle responds
* Lambda just do jobs with correct request body