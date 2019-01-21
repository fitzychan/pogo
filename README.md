PoGo API
========
Proof of concept Shell as a Service (ShaaS). Exposes execution of PowerShell commands and scripts via REST methods. Can easily be extended to any shell (planned in future releases).

### Request Contexts
* __/command/__ - Executes named PowerShell command.
* __/scripts/__ - Executes named script.

### Installing/Running as Windows Service
```batchfile
    sc create PoGo binPath= C:\path\to\bin\pogo.exe
    net start pogo
```

### Running Commands
Commands are handled by the /command/ context. Most common commands are readily supported and will return any structured data in JSON format. Parameters are passed via url query parameters. Named values will be broken out into key/value pairs and added to the commandstring.

Run __[Get-Date](https://technet.microsoft.com/en-us/library/hh849887.aspx)__ command.
    http://127.0.0.1:8080/command/Get-Date
```json
{
    "value":  "\/Date(1436237104231)\/",
    "DisplayHint":  2,
    "DateTime":  "Monday, July 6, 2015 9:45:04 PM"
}
```

Run __[Write-Host](https://technet.microsoft.com/en-us/library/ee177031.aspx)__ command.
    http://127.0.0.1:8080/command/Write-Host?"Hello,%20World!"
```json
Hello, World!
```

Run __[Get-Item](https://technet.microsoft.com/en-us/library/hh849788.aspx)__ command.
    http://127.0.0.1:8080/command/Get-Item?-Path="pogo.exe"
```json
{
    "Name":  "pogo.exe",
    "Length":  9076224,
    "DirectoryName":  "C:\\Users\\onyxhat\\Documents\\GitHub\\pogo",
    "Directory":  {
                      "Name":  "pogo",
                      "Parent":  {
                                     "Name":  "GitHub",
                                     "Parent":  "Documents",
                                     "Exists":  true,
                                     "Root":  "C:\\",
                                     "FullName":  "C:\\Users\\onyxhat\\Documents\\GitHub",
                                     "Extension":  "",
                                     "CreationTime":  "\/Date(1436236701491)\/",
                                     "CreationTimeUtc":  "\/Date(1436236701491)\/",
                                     "LastAccessTime":  "\/Date(1436236829685)\/",
                                     "LastAccessTimeUtc":  "\/Date(1436236829685)\/",
                                     "LastWriteTime":  "\/Date(1436236829685)\/",
                                     "LastWriteTimeUtc":  "\/Date(1436236829685)\/",
                                     "Attributes":  16
                                 },
                      "Exists":  true,
...
}
```

### Running Scripts
Scripts are handled by the __/scripts/__ context. By default custom scripts are stored in the __./scripts/__ directory relative to the service binary, but the location can be customized in the __config.json__. Much like the commands - parameters are passed via url query params. There is a simple test script included with the code that will allow the input of __-Name__. Adding/removing custom scripts requires no restarts nor recompiles.

Run script with no parameters (using defaults).
    http://127.0.0.1:8080/scripts/Hello-World
```json
Hello, World!
```

Run script to greet 'Isaac'.
    http://127.0.0.1:8080/scripts/Hello-World?-Name=Isaac
```json
Hello, Isaac!
```

### TODO
* Add ___Authentication___
* Implement POST method of form values as parameters (takes precedence over query params)
* Command restrictions
* Additional configuration values
* Remote configuration
* Fix relative path breakage while running as a service
