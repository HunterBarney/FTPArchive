# FTPArchive
FTPArchive is a CLI tool written in the go language that downloads files from an FTP or SFTP site,
Archives them, and then uploads them to either AWS S3 or Google Cloud Storage. <br>
You can also set SMTP credentials up so that when the program fails or succeeds it will send an email. <br>
To use the program you will need to create a profile file and a configuration file.

# Configuration File
If you run the program without a config.json file it will create one for you. The default configuration file is as follows:
```json
{
  "retryCount":3,
  "retryDelay":3,
  "logDirectory":"logs",
  "downloadDirectory":"downloads",
  "archiveDirectory":"archives",
  "sendEmail":true,
  "smtp":
    {
      "host":"",
      "port":0,
      "username":"",
      "password":"",
      "from":"",
      "to":[],
      "cc":[],
      "bcc":[]
    },
  "sendLogOverEmail":true
}
```
## Configuration field descriptions
**retryCount**: How many times the program will retry a download before throwing an error and exiting. <br>
**retryDelay**: The amount of seconds the program will wait before retrying a download. <br>
**logDirectory**: The directory that logs will be stored in. Accepts both absolute and relative paths. <br>
**downloadDirectory**: The directory that downloads will be stored in. Accepts both absolute and relative paths <br>
**archiveDirectory**: The directory that archives will be stored in. Accepts both absolute and relative paths <br>
**sendEmail**: If set to true, the program will send an email when there the program exits. It will include an error in the body of the email if the program failed. <br>
**smtp**: A set of values that decide the mail server and recipients that mail will be sent to. If the profile file also contains SMTP info, the ones in the configuration file will be overwritten. <br>
**sendLogOverEmail**: If set to true, the program will attach the log 

# Profile File
The program by default searches for a profile in the same directory as it called "profile.json". If your profile is 
named anything other than "profile.json" or is in any other folder other than the one the program is in you must pass 
the program the file by using the argument `-profile name.json` where "name.json" is replaced with the profile path and name.
<br>
Example profile:
```json
{
  "hostName": "ftp.site.com",
  "port": 21,
  "username": "username",
  "password": "password",
  "protocol": "FTP",
  "downloads": [
    "file1.txt",
    "file2.pdf",
    "/home/user/testdir"
  ],
  "outputName": "backup",
  "uploadPlatform": "aws",
  "bucketName": "archive",
  "smtp": {
    "host": "",
    "port": 0,
    "username": "",
    "password": "",
    "from": "",
    "to": [
      ""
    ],
    "cc": [
      ""
    ],
    "bcc": [
      ""
    ],
    "cleanupDownloads": true,
    "cleanupArchives": true,
    "cleanupOnFail": true
  }
}
```

## Profile field descriptions
**hostName**: The web address of the remote site you want to connect to. <br>
**port**: The port of the remote site you want to connect to. <br>
**username**: The username you will use to authenticate to the remote site with. <br>
**password**: The password you will use to authenticate to the remote site with. <br>
**protocol**: The protocol you will use to connect to the remote site with. Can be either FTP or SFTP. <br>
**downloads**: A list of files you would like to download. Can handle whole directories and singular files. <br>
**outputName**: The name of the Archive that the downloaded files will be placed into. <br>
**uploadPlatform**: The platform you will be uploading the archive to. Can be either AWS or GCP. <br>
**bucketName**: The bucket that the archive will be placed into on the upload platform. <br>
**smtp**: A set of values to connect to a mail server and recipients that mail will be sent to. These values will overwrite anything put in the smtp section of the config file. <br>
**cleanupDownloads**: Deletes the downloaded files after program execution. <br>
**cleanupArchives**: Deletes the archived files after program exectuion. <br>
**cleanupOnFail**: Enables the cleanup function to run when the program fails to finish.


# Manual mode
Manual mode will allow you to run only specific parts of the program by passing the following args <br>
**-download**: will enable downloading in manual mode<br>
**-archive**: will enable archiving in manual mode<br>
**-upload**: will enable uploading in manual mode

example usage: <br>
`./ftparchive -profile testprofile.json -download -archive` will only download and zip the files from a remote site. <br>
`./ftparchive -profile testprofile.json -archive` will only run the archive part of the program. This will fail if the files have not been downloaded already.
