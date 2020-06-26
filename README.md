# datagen
Data Generation Tool

Import V1 branch 

Navigate into datagen folder

RUN export PATH=$PATH:/usr/local/go/bin [To set go variables. Skip it if already set]

go get gopkg.in/yaml.v2

There will be 'output.txt' int /datagen. you can either flush the contents/delete the file to have just the output of this run. Or change the variable: output_file_name in file: sampleconfigV2.yaml

RUN go build main.go

RUN go run main.go
