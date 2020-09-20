## Primitive Server


A web server that acts as an image transformation service. It allows a user to upload an image, and 
then guides them through a selection process using various options on the [primitive](https://github.com/fogleman/primitive) 
CLI to transform the image.

## Requirements
1. Must have the ```primitive``` package installed. To install, run:
   ```
   go get -u github.com/fogleman/primitive
   ```

## Usage
1. Clone the repository
   
2. Run the server
   ```bash
   go run server.go
   ```

3. Go to `http://localhost:5000` in your browser, and upload an image file
   (currently only supports `.jpg`, `.jpeg` and `.png` file extensions)
  
 ## Note:
 Currently a work in progress, will make a front-end client in React.js to view various image transformations, and eventually deploy the server onto the cloud.